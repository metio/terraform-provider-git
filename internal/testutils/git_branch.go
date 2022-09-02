//go:build testing

/*
 * SPDX-FileCopyrightText: The terraform-provider-git Authors
 * SPDX-License-Identifier: 0BSD
 */

package testutils

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"testing"
)

func CreateBranch(t *testing.T, repository *git.Repository, branch *config.Branch) {
	err := repository.CreateBranch(branch)
	if err != nil {
		t.Fatal(err)
	}
	head, err := repository.Head()
	if err != nil {
		t.Fatal(err)
	}
	reference := plumbing.NewHashReference(plumbing.NewBranchReferenceName(branch.Name), head.Hash())
	err = repository.Storer.SetReference(reference)
	if err != nil {
		t.Fatal(err)
	}
}
