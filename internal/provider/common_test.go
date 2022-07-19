/*
 * This file is part of terraform-provider-git. It is subject to the license terms in the LICENSE file found in the top-level
 * directory of this distribution and at https://creativecommons.org/publicdomain/zero/1.0/. No part of terraform-provider-git,
 * including this file, may be copied, modified, propagated, or distributed except according to the terms contained
 * in the LICENSE file.
 */

package provider

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"io/ioutil"
	"path/filepath"
	"testing"
	"time"
)

func protoV6ProviderFactories() map[string]func() (tfprotov6.ProviderServer, error) {
	return map[string]func() (tfprotov6.ProviderServer, error){
		"git": providerserver.NewProtocol6WithError(New()),
	}
}

func initializeGitRepository(t *testing.T) (string, *git.Repository) {
	directory, err := ioutil.TempDir("", "terraform-provider-git")
	if err != nil {
		t.Fatal(err)
	}
	repository, err := git.PlainInit(directory, false)
	if err != nil {
		t.Fatal(err)
	}
	return directory, repository
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

func addAndCommitNewFile(t *testing.T, worktree *git.Worktree) {
	filename := filepath.Join(worktree.Filesystem.Root(), "example-git-file")
	err := ioutil.WriteFile(filename, []byte("hello world!"), 0644)
	if err != nil {
		t.Fatal(err)
	}
	_, err = worktree.Add("example-git-file")
	if err != nil {
		t.Fatal(err)
	}
	_, err = worktree.Commit("example go-git commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Some Person",
			Email: "person@example.com",
			When:  time.Now(),
		},
	})
	if err != nil {
		t.Fatal(err)
	}
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

func testCheckLen(expectedLength int) func(input string) error {
	return func(input string) error {
		if len(input) != expectedLength {
			return fmt.Errorf("expected length %d, actual length %d", expectedLength, len(input))
		}

		return nil
	}
}
