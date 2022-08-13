/*
 * SPDX-FileCopyrightText: The terraform-provider-git Authors
 * SPDX-License-Identifier: 0BSD
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
	t.Parallel()
	directory, repository := testRepository(t)
	defer os.RemoveAll(directory)
	branch := "name-of-branch"
	remote := "origin"
	rebase := "true"
	testCreateBranch(t, repository, &config.Branch{
		Name:   branch,
		Remote: remote,
		Rebase: rebase,
	})

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testProviderFactories(),
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
	t.Parallel()
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testProviderFactories(),
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
	t.Parallel()
	directory, _ := testRepository(t)
	defer os.RemoveAll(directory)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testProviderFactories(),
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
	t.Parallel()
	directory, _ := testRepository(t)
	defer os.RemoveAll(directory)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testProviderFactories(),
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
	t.Parallel()
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testProviderFactories(),
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
