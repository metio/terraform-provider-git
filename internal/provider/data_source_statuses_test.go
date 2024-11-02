/*
 * SPDX-FileCopyrightText: The terraform-provider-git Authors
 * SPDX-License-Identifier: 0BSD
 */

package provider_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/metio/terraform-provider-git/internal/testutils"
	"testing"
)

func TestDataSourceGitStatuses_Unclean(t *testing.T) {
	t.Parallel()
	directory, repository := testutils.CreateRepository(t)
	worktree := testutils.GetRepositoryWorktree(t, repository)
	fileName := "some-file"
	testutils.WriteFileInWorktree(t, worktree, fileName)
	testutils.GitAdd(t, worktree, fileName)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "git_statuses" "test" {
						directory = "%s"
					}
				`, directory),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.git_statuses.test", "directory", directory),
					resource.TestCheckResourceAttr("data.git_statuses.test", "id", directory),
					resource.TestCheckResourceAttr("data.git_statuses.test", "is_clean", "false"),
					resource.TestCheckResourceAttr("data.git_statuses.test", "files.%", "1"),
				),
			},
		},
	})
}

func TestDataSourceGitStatuses_Clean(t *testing.T) {
	t.Parallel()
	directory, repository := testutils.CreateRepository(t)
	worktree := testutils.GetRepositoryWorktree(t, repository)
	fileName := "some-file"
	testutils.WriteFileInWorktree(t, worktree, fileName)
	testutils.GitAdd(t, worktree, fileName)
	testutils.GitCommit(t, worktree)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "git_statuses" "test" {
						directory = "%s"
					}
				`, directory),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.git_statuses.test", "directory", directory),
					resource.TestCheckResourceAttr("data.git_statuses.test", "id", directory),
					resource.TestCheckResourceAttr("data.git_statuses.test", "is_clean", "true"),
					resource.TestCheckResourceAttr("data.git_statuses.test", "files.%", "0"),
				),
			},
		},
	})
}

func TestDataSourceGitStatuses_BareRepository(t *testing.T) {
	t.Parallel()
	directory := testutils.CreateBareRepository(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "git_statuses" "test" {
						directory = "%s"
					}
				`, directory),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.git_statuses.test", "directory", directory),
					resource.TestCheckResourceAttr("data.git_statuses.test", "id", directory),
					resource.TestCheckResourceAttr("data.git_statuses.test", "is_clean", "true"),
					resource.TestCheckResourceAttr("data.git_statuses.test", "files.%", "0"),
				),
			},
		},
	})
}
