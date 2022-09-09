//go:build testing

/*
 * SPDX-FileCopyrightText: The terraform-provider-git Authors
 * SPDX-License-Identifier: 0BSD
 */

package testutils

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"path/filepath"
	"testing"
)

func CreateRepository(t *testing.T) (string, *git.Repository) {
	directory := TemporaryDirectory(t)
	repository := GitInit(t, directory, false)
	return directory, repository
}

func CreateBareRepository(t *testing.T) string {
	directory := TemporaryDirectory(t)
	GitInit(t, directory, true)
	return directory
}

func GetRepositoryWorktree(t *testing.T, repository *git.Repository) *git.Worktree {
	worktree, err := repository.Worktree()
	if err != nil {
		t.Fatal(err)
	}
	return worktree
}

func GetRepositoryHead(t *testing.T, repository *git.Repository) *plumbing.Reference {
	head, err := repository.Head()
	if err != nil {
		t.Fatal(err)
	}
	return head
}

func WriteFileInWorktree(t *testing.T, worktree *git.Worktree, name string) {
	filename := FileInWorktree(worktree, name)
	WriteFile(t, filename)
}

func FileInWorktree(worktree *git.Worktree, name string) string {
	return filepath.Join(worktree.Filesystem.Root(), name)
}

func AddAndCommitNewFile(t *testing.T, worktree *git.Worktree, name string) {
	WriteFileInWorktree(t, worktree, name)
	GitAdd(t, worktree, name)
	GitCommit(t, worktree)
}
