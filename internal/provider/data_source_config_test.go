/*
 * SPDX-FileCopyrightText: The terraform-provider-git Authors
 * SPDX-License-Identifier: 0BSD
 */

package provider_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/metio/terraform-provider-git/internal/testutils"
)

func TestDataSourceGitConfig(t *testing.T) {
	t.Parallel()
	directory, repository := testutils.CreateRepository(t)
	cfg := testutils.TestConfig(t, repository)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "git_config" "test" {
						directory = "%s"
					}
				`, directory),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.git_config.test", "directory", directory),
					resource.TestCheckResourceAttr("data.git_config.test", "id", directory),
					resource.TestCheckResourceAttr("data.git_config.test", "scope", "global"),
					resource.TestCheckResourceAttr("data.git_config.test", "user_name", cfg.User.Name),
					resource.TestCheckResourceAttr("data.git_config.test", "user_email", cfg.User.Email),
					resource.TestCheckResourceAttr("data.git_config.test", "author_name", cfg.Author.Name),
					resource.TestCheckResourceAttr("data.git_config.test", "author_email", cfg.Author.Email),
					resource.TestCheckResourceAttr("data.git_config.test", "committer_name", cfg.Committer.Name),
					resource.TestCheckResourceAttr("data.git_config.test", "committer_email", cfg.Committer.Email),
				),
			},
		},
	})
}

func TestDataSourceGitConfig_InvalidRepository(t *testing.T) {
	t.Parallel()
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: `
					data "git_config" "test" {
						directory = "/some/random/path"
					}
				`,
				ExpectError: regexp.MustCompile(`Cannot open repository`),
			},
		},
	})
}

func TestDataSourceGitConfig_InvalidScope(t *testing.T) {
	t.Parallel()
	directory, _ := testutils.CreateRepository(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "git_config" "test" {
						directory = "%s"
						scope     = "unknown-scope"
					}
				`, directory),
				ExpectError: regexp.MustCompile(`Invalid Attribute Value Match`),
			},
		},
	})
}

func TestDataSourceGitConfig_ScopeLocal(t *testing.T) {
	t.Parallel()
	directory, repository := testutils.CreateRepository(t)
	cfg := testutils.TestConfig(t, repository)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "git_config" "test" {
						directory = "%s"
						scope     = "local"
					}
				`, directory),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.git_config.test", "directory", directory),
					resource.TestCheckResourceAttr("data.git_config.test", "id", directory),
					resource.TestCheckResourceAttr("data.git_config.test", "scope", "local"),
					resource.TestCheckResourceAttr("data.git_config.test", "user_name", cfg.User.Name),
					resource.TestCheckResourceAttr("data.git_config.test", "user_email", cfg.User.Email),
					resource.TestCheckResourceAttr("data.git_config.test", "author_name", cfg.Author.Name),
					resource.TestCheckResourceAttr("data.git_config.test", "author_email", cfg.Author.Email),
					resource.TestCheckResourceAttr("data.git_config.test", "committer_name", cfg.Committer.Name),
					resource.TestCheckResourceAttr("data.git_config.test", "committer_email", cfg.Committer.Email),
				),
			},
		},
	})
}

func TestDataSourceGitConfig_ScopeGlobal(t *testing.T) {
	t.Parallel()
	directory, repository := testutils.CreateRepository(t)
	cfg := testutils.TestConfig(t, repository)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "git_config" "test" {
						directory = "%s"
						scope     = "global"
					}
				`, directory),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.git_config.test", "directory", directory),
					resource.TestCheckResourceAttr("data.git_config.test", "id", directory),
					resource.TestCheckResourceAttr("data.git_config.test", "scope", "global"),
					resource.TestCheckResourceAttr("data.git_config.test", "user_name", cfg.User.Name),
					resource.TestCheckResourceAttr("data.git_config.test", "user_email", cfg.User.Email),
					resource.TestCheckResourceAttr("data.git_config.test", "author_name", cfg.Author.Name),
					resource.TestCheckResourceAttr("data.git_config.test", "author_email", cfg.Author.Email),
					resource.TestCheckResourceAttr("data.git_config.test", "committer_name", cfg.Committer.Name),
					resource.TestCheckResourceAttr("data.git_config.test", "committer_email", cfg.Committer.Email),
				),
			},
		},
	})
}

func TestDataSourceGitConfig_ScopeSystem(t *testing.T) {
	t.Parallel()
	directory, repository := testutils.CreateRepository(t)
	cfg := testutils.TestConfig(t, repository)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "git_config" "test" {
						directory = "%s"
						scope     = "system"
					}
				`, directory),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.git_config.test", "directory", directory),
					resource.TestCheckResourceAttr("data.git_config.test", "id", directory),
					resource.TestCheckResourceAttr("data.git_config.test", "scope", "system"),
					resource.TestCheckResourceAttr("data.git_config.test", "user_name", cfg.User.Name),
					resource.TestCheckResourceAttr("data.git_config.test", "user_email", cfg.User.Email),
					resource.TestCheckResourceAttr("data.git_config.test", "author_name", cfg.Author.Name),
					resource.TestCheckResourceAttr("data.git_config.test", "author_email", cfg.Author.Email),
					resource.TestCheckResourceAttr("data.git_config.test", "committer_name", cfg.Committer.Name),
					resource.TestCheckResourceAttr("data.git_config.test", "committer_email", cfg.Committer.Email),
				),
			},
		},
	})
}

func TestDataSourceGitConfig_MissingRepository(t *testing.T) {
	t.Parallel()
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: `
					data "git_config" "test" {}
				`,
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
		},
	})
}
