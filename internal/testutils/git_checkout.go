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

func TestGitCheckout(t *testing.T, worktree *git.Worktree, hash plumbing.Hash) {
	err := worktree.Checkout(&git.CheckoutOptions{
		Hash: hash,
	})
	if err != nil {
		t.Fatal(err)
	}
}
