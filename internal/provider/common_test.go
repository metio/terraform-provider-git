/*
 * This file is part of terraform-provider-git. It is subject to the license terms in the LICENSE file found in the top-level
 * directory of this distribution and at https://creativecommons.org/publicdomain/zero/1.0/. No part of terraform-provider-git,
 * including this file, may be copied, modified, propagated, or distributed except according to the terms contained
 * in the LICENSE file.
 */

package provider_test

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/metio/terraform-provider-git/internal/provider"
	"io/ioutil"
	"path/filepath"
	"testing"
	"time"
)

func protoV6ProviderFactories() map[string]func() (tfprotov6.ProviderServer, error) {
	return map[string]func() (tfprotov6.ProviderServer, error){
		"git": providerserver.NewProtocol6WithError(provider.New()),
	}
}

func initializeGitRepository(t *testing.T) (string, *git.Repository) {
	directory := temporaryDirectory(t)
	repository := gitInit(t, directory, false)
	return directory, repository
}

func gitInit(t *testing.T, directory string, bare bool) *git.Repository {
	repository, err := git.PlainInit(directory, bare)
	if err != nil {
		t.Fatal(err)
	}
	return repository
}

func temporaryDirectory(t *testing.T) string {
	directory, err := ioutil.TempDir("", "terraform-provider-git")
	if err != nil {
		t.Fatal(err)
	}
	return filepath.ToSlash(directory)
}

func createBranch(t *testing.T, repository *git.Repository, branch *config.Branch) {
	err := repository.CreateBranch(branch)
	if err != nil {
		t.Fatal(err)
	}
}

func createWorktree(t *testing.T, repository *git.Repository) *git.Worktree {
	worktree, err := repository.Worktree()
	if err != nil {
		t.Fatal(err)
	}
	return worktree
}

func addFile(t *testing.T, worktree *git.Worktree, name string) {
	filename := filepath.Join(worktree.Filesystem.Root(), name)
	err := ioutil.WriteFile(filename, []byte("hello world!"), 0644)
	if err != nil {
		t.Fatal(err)
	}
	_, err = worktree.Add(name)
	if err != nil {
		t.Fatal(err)
	}
}

func commitStaged(t *testing.T, worktree *git.Worktree) {
	_, err := worktree.Commit("example go-git commit", &git.CommitOptions{
		Author: signature(),
	})
	if err != nil {
		t.Fatal(err)
	}
}

func addAndCommitNewFile(t *testing.T, worktree *git.Worktree, name string) {
	addFile(t, worktree, name)
	commitStaged(t, worktree)
}

func initTestConfig(t *testing.T, repository *git.Repository) *config.Config {
	cfg := readConfig(t, repository)
	cfg.User.Name = "user name"
	cfg.User.Email = "user@example.com"
	cfg.Author.Name = "author name"
	cfg.Author.Email = "author@example.com"
	cfg.Committer.Name = "committer name"
	cfg.Committer.Email = "committer@example.com"
	writeConfig(t, repository, cfg)
	return cfg
}

func readConfig(t *testing.T, repository *git.Repository) *config.Config {
	cfg, err := repository.Config()
	if err != nil {
		t.Fatal(err)
	}
	return cfg
}

func writeConfig(t *testing.T, repository *git.Repository, cfg *config.Config) {
	err := repository.SetConfig(cfg)
	if err != nil {
		t.Fatal(err)
	}
}

func createRemote(t *testing.T, repository *git.Repository, remote string) {
	_, err := repository.CreateRemote(&config.RemoteConfig{
		Name: remote,
		URLs: []string{"https://example.com/metio/terraform-provider-git.git"},
	})
	if err != nil {
		t.Fatal(err)
	}
}

func createTag(t *testing.T, repository *git.Repository, tag string) {
	head, err := repository.Head()
	if err != nil {
		t.Fatal(err)
	}
	_, err = repository.CreateTag(tag, head.Hash(), &git.CreateTagOptions{
		Message: tag,
		Tagger:  signature(),
	})
	if err != nil {
		t.Fatal(err)
	}
}

func signature() *object.Signature {
	return &object.Signature{
		Name:  "Some Person",
		Email: "person@example.com",
		When:  time.Now(),
	}
}

func testCheckLen(expectedLength int) func(input string) error {
	return func(input string) error {
		if len(input) != expectedLength {
			return fmt.Errorf("expected length %d, actual length %d", expectedLength, len(input))
		}

		return nil
	}
}
