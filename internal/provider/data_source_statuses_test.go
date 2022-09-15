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
	"testing"
)

func TestDataSourceGitStatuses_GetSchema(t *testing.T) {
	t.Parallel()
	r := &provider.StatusesDataSource{}
	schema, _ := r.GetSchema(context.TODO())

	testutils.VerifySchemaDescriptions(t, schema)
}

func TestDataSourceGitStatuses_Unclean(t *testing.T) {
	t.Parallel()
	directory, repository := testutils.CreateRepository(t)
	defer os.RemoveAll(directory)
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
	defer os.RemoveAll(directory)
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
	defer os.RemoveAll(directory)

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