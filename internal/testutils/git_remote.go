//go:build testing

/*
 * SPDX-FileCopyrightText: The terraform-provider-git Authors
 * SPDX-License-Identifier: 0BSD
 */

package testutils

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"testing"
)

func CreateRemote(t *testing.T, repository *git.Repository, remote string) {
	CreateRemoteWithName(t, repository, remote, "terraform-provider-git")
}

func CreateRemoteWithName(t *testing.T, repository *git.Repository, remote string, name string) {
	CreateRemoteWithUrls(t, repository, remote, []string{fmt.Sprintf("https://example.com/metio/%s.git", name)})
}

func CreateRemoteWithUrls(t *testing.T, repository *git.Repository, remote string, urls []string) {
	_, err := repository.CreateRemote(&config.RemoteConfig{
		Name: remote,
		URLs: urls,
	})
	if err != nil {
		t.Fatal(err)
	}
}
