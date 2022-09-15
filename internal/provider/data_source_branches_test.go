/*
 * SPDX-FileCopyrightText: The terraform-provider-git Authors
 * SPDX-License-Identifier: 0BSD
 */

package provider_test

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/metio/terraform-provider-git/internal/provider"
	"github.com/metio/terraform-provider-git/internal/testutils"
	"os"
	"regexp"
	"testing"
)

func TestDataSourceGitBranches_GetSchema(t *testing.T) {
	t.Parallel()
	r := &provider.BranchesDataSource{}
	schema, _ := r.GetSchema(context.TODO())

	testutils.VerifySchemaDescriptions(t, schema)
}

func TestDataSourceGitBranches(t *testing.T) {
	t.Parallel()
	directory, repository := testutils.CreateRepository(t)
	defer os.RemoveAll(directory)
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