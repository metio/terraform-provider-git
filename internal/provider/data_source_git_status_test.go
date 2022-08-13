/*
 * SPDX-FileCopyrightText: The terraform-provider-git Authors
 * SPDX-License-Identifier: 0BSD
 */

package provider_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"os"
	"testing"
)

func TestDataSourceGitStatus_StagedFile(t *testing.T) {
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
	directory := testTemporaryDirectory(t)
	testGitInit(t, directory, true)
	defer os.RemoveAll(directory)
	fileName := "some-file"

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testProviderFactories(),
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
					resource.TestCheckResourceAttr("data.git_status.test", "staging", ""),
					resource.TestCheckResourceAttr("data.git_status.test", "worktree", ""),
				),
			},
		},
	})
}
