/*
 * SPDX-FileCopyrightText: The terraform-provider-git Authors
 * SPDX-License-Identifier: 0BSD
 */

package provider_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/metio/terraform-provider-git/internal/testutils"
	"regexp"
	"testing"
)

func TestDataSourceGitBranches(t *testing.T) {
	t.Parallel()
	directory, repository := testutils.CreateRepository(t)
	worktree := testutils.GetRepositoryWorktree(t, repository)
	testutils.AddAndCommitNewFile(t, worktree, "some-file")

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "git_branches" "test" {
						directory = "%s"
					}
				`, directory),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.git_branches.test", "directory", directory),
					resource.TestCheckResourceAttr("data.git_branches.test", "id", directory),
					resource.TestCheckResourceAttr("data.git_branches.test", "branches.%", "1"),
					resource.TestCheckNoResourceAttr("data.git_branches.test", "branches.master.remote"),
					resource.TestCheckNoResourceAttr("data.git_branches.test", "branches.master.rebase"),
					resource.TestCheckResourceAttrWith("data.git_branches.test", "branches.master.sha1", testutils.CheckExactLength(40)),
				),
			},
		},
	})
}

func TestDataSourceGitBranches_InvalidRepository(t *testing.T) {
	t.Parallel()
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: `
					data "git_branches" "test" {
						directory = "/some/random/path"
					}
				`,
				ExpectError: regexp.MustCompile(`Cannot open repository`),
			},
		},
	})
}

func TestDataSourceGitBranches_MissingRepository(t *testing.T) {
	t.Parallel()
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: `
					data "git_branches" "test" {}
				`,
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
		},
	})
}
