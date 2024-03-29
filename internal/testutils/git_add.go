//go:build testing

/*
 * SPDX-FileCopyrightText: The terraform-provider-git Authors
 * SPDX-License-Identifier: 0BSD
 */

package testutils

import (
	"github.com/go-git/go-git/v5"
	"testing"
)

func GitAdd(t *testing.T, worktree *git.Worktree, name string) {
	_, err := worktree.Add(name)
	if err != nil {
		t.Fatal(err)
	}
}
