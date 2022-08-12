/*
 * SPDX-FileCopyrightText: Sebastian Hoß <seb@hoß.de>
 * SPDX-License-Identifier: 0BSD
 */

package provider_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"os"
	"testing"
)

func TestDataSourceGitStatuses_Unclean(t *testing.T) {
	t.Parallel()
	directory, repository := testRepository(t)
	defer os.RemoveAll(directory)
	worktree := testWorktree(t, repository)
	fileName := "some-file"
	testWriteFile(t, worktree, fileName)
	testGitAdd(t, worktree, fileName)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testProviderFactories(),
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
	directory, repository := testRepository(t)
	defer os.RemoveAll(directory)
	worktree := testWorktree(t, repository)
	fileName := "some-file"
	testWriteFile(t, worktree, fileName)
	testGitAdd(t, worktree, fileName)
	testGitCommit(t, worktree)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testProviderFactories(),
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
	directory := testTemporaryDirectory(t)
	testGitInit(t, directory, true)
	defer os.RemoveAll(directory)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testProviderFactories(),
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
