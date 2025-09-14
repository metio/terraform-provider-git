/*
 * SPDX-FileCopyrightText: The terraform-provider-git Authors
 * SPDX-License-Identifier: 0BSD
 */

package provider_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/metio/terraform-provider-git/internal/testutils"
)

func TestDataSourceGitStatus_StagedFile(t *testing.T) {
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
					data "git_status" "test" {
						directory = "%s"
						file      = "%s"
					}
				`, directory, fileName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.git_status.test", "directory", directory),
					resource.TestCheckResourceAttr("data.git_status.test", "id", fileName),
					resource.TestCheckResourceAttr("data.git_status.test", "file", fileName),
					resource.TestCheckResourceAttr("data.git_status.test", "staging", "A"),
					resource.TestCheckResourceAttr("data.git_status.test", "worktree", " "),
				),
			},
		},
	})
}

func TestDataSourceGitStatus_Clean(t *testing.T) {
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
					data "git_status" "test" {
						directory = "%s"
						file      = "%s"
					}
				`, directory, fileName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.git_status.test", "directory", directory),
					resource.TestCheckResourceAttr("data.git_status.test", "id", fileName),
					resource.TestCheckResourceAttr("data.git_status.test", "file", fileName),
					resource.TestCheckResourceAttr("data.git_status.test", "staging", "?"),
					resource.TestCheckResourceAttr("data.git_status.test", "worktree", "?"),
				),
			},
		},
	})
}

func TestDataSourceGitStatus_BareRepository(t *testing.T) {
	t.Parallel()
	directory := testutils.CreateBareRepository(t)
	fileName := "some-file"

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "git_status" "test" {
						directory = "%s"
						file      = "%s"
					}
				`, directory, fileName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.git_status.test", "directory", directory),
					resource.TestCheckResourceAttr("data.git_status.test", "id", fileName),
					resource.TestCheckResourceAttr("data.git_status.test", "file", fileName),
					resource.TestCheckNoResourceAttr("data.git_status.test", "staging"),
					resource.TestCheckNoResourceAttr("data.git_status.test", "worktree"),
				),
			},
		},
	})
}
