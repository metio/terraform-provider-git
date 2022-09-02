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

func TestConfig(t *testing.T, repository *git.Repository) *config.Config {
	cfg := ReadConfig(t, repository)
	cfg.User.Name = "user name"
	cfg.User.Email = "user@example.com"
	cfg.Author.Name = "author name"
	cfg.Author.Email = "author@example.com"
	cfg.Committer.Name = "committer name"
	cfg.Committer.Email = "committer@example.com"
	WriteConfig(t, repository, cfg)
	return cfg
}

func ReadConfig(t *testing.T, repository *git.Repository) *config.Config {
	cfg, err := repository.Config()
	if err != nil {
		t.Fatal(err)
	}
	return cfg
}

func WriteConfig(t *testing.T, repository *git.Repository, cfg *config.Config) {
	err := repository.SetConfig(cfg)
	if err != nil {
		t.Fatal(err)
	}
}
