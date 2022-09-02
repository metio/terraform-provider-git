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

func CreateTag(t *testing.T, repository *git.Repository, tag string) {
	head := GetRepositoryHead(t, repository)
	_, err := repository.CreateTag(tag, head.Hash(), &git.CreateTagOptions{
		Message: tag,
		Tagger:  Signature(),
	})
	if err != nil {
		t.Fatal(err)
	}
}
