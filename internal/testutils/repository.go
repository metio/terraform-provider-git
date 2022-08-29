/*
 * SPDX-FileCopyrightText: The terraform-provider-git Authors
 * SPDX-License-Identifier: 0BSD
 */

package testutils

import (
	"github.com/go-git/go-git/v5"
	"testing"
)

func TestRepository(t *testing.T) (string, *git.Repository) {
	directory := testTemporaryDirectory(t)
	repository := testGitInit(t, directory, false)
	return directory, repository
}

func TestRepositoryBare(t *testing.T) string {
	directory := testTemporaryDirectory(t)
	testGitInit(t, directory, true)
	return directory
}
