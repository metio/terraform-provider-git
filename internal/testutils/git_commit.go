//go:build testing

/*
 * SPDX-FileCopyrightText: The terraform-provider-git Authors
 * SPDX-License-Identifier: 0BSD
 */

package testutils

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"testing"
)

func GitCommit(t *testing.T, worktree *git.Worktree) plumbing.Hash {
	return GitCommitWith(t, worktree, &git.CommitOptions{
		Author:    Signature(),
		Committer: Signature(),
	})
}

func GitCommitWith(t *testing.T, worktree *git.Worktree, options *git.CommitOptions) plumbing.Hash {
	commit, err := worktree.Commit("example go-git commit", options)
	if err != nil {
		t.Fatal(err)
	}
	return commit
}
