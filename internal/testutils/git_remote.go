//go:build testing

/*
 * SPDX-FileCopyrightText: The terraform-provider-git Authors
 * SPDX-License-Identifier: 0BSD
 */

package testutils

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"testing"
)

func CreateRemote(t *testing.T, repository *git.Repository, remote string) {
	_, err := repository.CreateRemote(&config.RemoteConfig{
		Name: remote,
		URLs: []string{"https://example.com/metio/terraform-provider-git.git"},
	})
	if err != nil {
		t.Fatal(err)
	}
}
