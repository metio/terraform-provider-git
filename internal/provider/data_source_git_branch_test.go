/*
 * This file is part of terraform-provider-git. It is subject to the license terms in the LICENSE file found in the top-level
 * directory of this distribution and at https://creativecommons.org/publicdomain/zero/1.0/. No part of terraform-provider-git,
 * including this file, may be copied, modified, propagated, or distributed except according to the terms contained
 * in the LICENSE file.
 */

package provider_test

import (
	"fmt"
	"github.com/go-git/go-git/v5/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"os"
	"regexp"
	"testing"
)

func TestDataSourceGitBranch(t *testing.T) {
	directory, repository := initializeGitRepository(t)
	defer os.RemoveAll(directory)
	branch := "name-of-branch"
	remote := "origin"
	rebase := "true"
	createBranch(t, repository, &config.Branch{
		Name:   branch,
		Remote: remote,
		Rebase: rebase,
	})

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: protoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "git_branch" "test" {
						directory = "%s"
						branch    = "%s"
					}
				`, directory, branch),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.git_branch.test", "directory", directory),
					resource.TestCheckResourceAttr("data.git_branch.test", "id", directory),
					resource.TestCheckResourceAttr("data.git_branch.test", "branch", branch),
					resource.TestCheckResourceAttr("data.git_branch.test", "remote", remote),
					resource.TestCheckResourceAttr("data.git_branch.test", "rebase", rebase),
				),
			},
		},
	})
}

func TestDataSourceGitBranch_InvalidRepository(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: protoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: `
					data "git_branch" "test" {
						directory = "/some/random/path"
						branch    = "this-does-not-exist"
					}
				`,
				ExpectError: regexp.MustCompile(`Cannot open repository`),
			},
		},
	})
}

func TestDataSourceGitBranch_InvalidBranch(t *testing.T) {
	directory, _ := initializeGitRepository(t)
	defer os.RemoveAll(directory)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: protoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "git_branch" "test" {
						directory = "%s"
						branch    = "does-not-exist"
					}
				`, directory),
				ExpectError: regexp.MustCompile(`Cannot read branch`),
			},
		},
	})
}

func TestDataSourceGitBranch_MissingBranch(t *testing.T) {
	directory, _ := initializeGitRepository(t)
	defer os.RemoveAll(directory)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: protoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "git_branch" "test" {
						directory = "%s"
					}
				`, directory),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
		},
	})
}

func TestDataSourceGitBranch_MissingRepository(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: protoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: `
					data "git_branch" "test" {}
				`,
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
		},
	})
}
