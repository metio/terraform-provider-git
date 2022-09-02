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

func GitInit(t *testing.T, directory string, bare bool) *git.Repository {
	repository, err := git.PlainInit(directory, bare)
	if err != nil {
		t.Fatal(err)
	}
	return repository
}
