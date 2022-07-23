/*
 * This file is part of terraform-provider-git. It is subject to the license terms in the LICENSE file found in the top-level
 * directory of this distribution and at https://creativecommons.org/publicdomain/zero/1.0/. No part of terraform-provider-git,
 * including this file, may be copied, modified, propagated, or distributed except according to the terms contained
 * in the LICENSE file.
 */

package provider_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"os"
	"regexp"
	"testing"
)

func TestDataSourceGitConfig(t *testing.T) {
	t.Parallel()
	directory, repository := initializeGitRepository(t)
	defer os.RemoveAll(directory)
	cfg := initTestConfig(t, repository)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: protoV6ProviderFactories(),
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
		ProtoV6ProviderFactories: protoV6ProviderFactories(),
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
	directory, _ := initializeGitRepository(t)
	defer os.RemoveAll(directory)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: protoV6ProviderFactories(),
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
	directory, repository := initializeGitRepository(t)
	defer os.RemoveAll(directory)
	cfg := initTestConfig(t, repository)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: protoV6ProviderFactories(),
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
	directory, repository := initializeGitRepository(t)
	defer os.RemoveAll(directory)
	cfg := initTestConfig(t, repository)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: protoV6ProviderFactories(),
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
	directory, repository := initializeGitRepository(t)
	defer os.RemoveAll(directory)
	cfg := initTestConfig(t, repository)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: protoV6ProviderFactories(),
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
		ProtoV6ProviderFactories: protoV6ProviderFactories(),
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
