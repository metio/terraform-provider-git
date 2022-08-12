/*
 * SPDX-FileCopyrightText: Sebastian Hoß <seb@hoß.de>
 * SPDX-License-Identifier: BSD0
 */

package provider_test

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/metio/terraform-provider-git/internal/provider"
	"io/ioutil"
	"path/filepath"
	"testing"
	"time"
)

func testProviderFactories() map[string]func() (tfprotov6.ProviderServer, error) {
	return map[string]func() (tfprotov6.ProviderServer, error){
		"git": providerserver.NewProtocol6WithError(provider.New()),
	}
}

func testRepository(t *testing.T) (string, *git.Repository) {
	directory := testTemporaryDirectory(t)
	repository := testGitInit(t, directory, false)
	return directory, repository
}

func testGitInit(t *testing.T, directory string, bare bool) *git.Repository {
	repository, err := git.PlainInit(directory, bare)
	if err != nil {
		t.Fatal(err)
	}
	return repository
}

func testTemporaryDirectory(t *testing.T) string {
	directory, err := ioutil.TempDir("", "terraform-provider-git")
	if err != nil {
		t.Fatal(err)
	}
	return filepath.ToSlash(directory)
}

func testCreateBranch(t *testing.T, repository *git.Repository, branch *config.Branch) {
	err := repository.CreateBranch(branch)
	if err != nil {
		t.Fatal(err)
	}
}

func testWorktree(t *testing.T, repository *git.Repository) *git.Worktree {
	worktree, err := repository.Worktree()
	if err != nil {
		t.Fatal(err)
	}
	return worktree
}

func testWriteFile(t *testing.T, worktree *git.Worktree, name string) {
	filename := filepath.Join(worktree.Filesystem.Root(), name)
	err := ioutil.WriteFile(filename, []byte("hello world!"), 0644)
	if err != nil {
		t.Fatal(err)
	}
}

func testGitAdd(t *testing.T, worktree *git.Worktree, name string) {
	_, err := worktree.Add(name)
	if err != nil {
		t.Fatal(err)
	}
}

func testGitCommit(t *testing.T, worktree *git.Worktree) plumbing.Hash {
	commit, err := worktree.Commit("example go-git commit", &git.CommitOptions{
		Author:    testSignature(),
		Committer: testSignature(),
	})
	if err != nil {
		t.Fatal(err)
	}
	return commit
}

func testAddAndCommitNewFile(t *testing.T, worktree *git.Worktree, name string) {
	testWriteFile(t, worktree, name)
	testGitAdd(t, worktree, name)
	testGitCommit(t, worktree)
}

func testConfig(t *testing.T, repository *git.Repository) *config.Config {
	cfg := testReadConfig(t, repository)
	cfg.User.Name = "user name"
	cfg.User.Email = "user@example.com"
	cfg.Author.Name = "author name"
	cfg.Author.Email = "author@example.com"
	cfg.Committer.Name = "committer name"
	cfg.Committer.Email = "committer@example.com"
	testWriteConfig(t, repository, cfg)
	return cfg
}

func testReadConfig(t *testing.T, repository *git.Repository) *config.Config {
	cfg, err := repository.Config()
	if err != nil {
		t.Fatal(err)
	}
	return cfg
}

func testWriteConfig(t *testing.T, repository *git.Repository, cfg *config.Config) {
	err := repository.SetConfig(cfg)
	if err != nil {
		t.Fatal(err)
	}
}

func testCreateRemote(t *testing.T, repository *git.Repository, remote string) {
	_, err := repository.CreateRemote(&config.RemoteConfig{
		Name: remote,
		URLs: []string{"https://example.com/metio/terraform-provider-git.git"},
	})
	if err != nil {
		t.Fatal(err)
	}
}

func testCreateTag(t *testing.T, repository *git.Repository, tag string) {
	head := testReadHead(t, repository)
	_, err := repository.CreateTag(tag, head.Hash(), &git.CreateTagOptions{
		Message: tag,
		Tagger:  testSignature(),
	})
	if err != nil {
		t.Fatal(err)
	}
}

func testReadHead(t *testing.T, repository *git.Repository) *plumbing.Reference {
	head, err := repository.Head()
	if err != nil {
		t.Fatal(err)
	}
	return head
}

func testSignature() *object.Signature {
	return &object.Signature{
		Name:  "Some Person",
		Email: "person@example.com",
		When:  time.Now(),
	}
}

func testCheckExactLength(expectedLength int) func(input string) error {
	return func(input string) error {
		if len(input) != expectedLength {
			return fmt.Errorf("expected length %d, actual length %d", expectedLength, len(input))
		}

		return nil
	}
}

func testCheckMinLength(minimumLength int) func(input string) error {
	return func(input string) error {
		if len(input) < minimumLength {
			return fmt.Errorf("minimum length %d, actual length %d", minimumLength, len(input))
		}

		return nil
	}
}
