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

func TestDataSourceGitCommit(t *testing.T) {
	t.Parallel()
	directory, repository := testRepository(t)
	defer os.RemoveAll(directory)
	testConfig(t, repository)
	worktree := testWorktree(t, repository)
	fileName := "some-file"
	testWriteFile(t, worktree, fileName)
	testGitAdd(t, worktree, fileName)
	commit := testGitCommit(t, worktree)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "git_commit" "test" {
						directory = "%s"
						sha1      = "%s"
					}
				`, directory, commit.String()),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.git_commit.test", "directory", directory),
					resource.TestCheckResourceAttr("data.git_commit.test", "id", commit.String()),
					resource.TestCheckResourceAttr("data.git_commit.test", "sha1", commit.String()),
					resource.TestCheckResourceAttr("data.git_commit.test", "signature", ""),
					resource.TestCheckResourceAttrWith("data.git_commit.test", "message", testCheckMinLength(1)),
					resource.TestCheckResourceAttrWith("data.git_commit.test", "tree_sha1", testCheckMinLength(1)),
					resource.TestCheckResourceAttrWith("data.git_commit.test", "author.name", testCheckMinLength(1)),
					resource.TestCheckResourceAttrWith("data.git_commit.test", "author.email", testCheckMinLength(1)),
					resource.TestCheckResourceAttrWith("data.git_commit.test", "author.timestamp", testCheckMinLength(1)),
					resource.TestCheckResourceAttrWith("data.git_commit.test", "committer.name", testCheckMinLength(1)),
					resource.TestCheckResourceAttrWith("data.git_commit.test", "committer.email", testCheckMinLength(1)),
					resource.TestCheckResourceAttrWith("data.git_commit.test", "committer.timestamp", testCheckMinLength(1)),
				),
			},
		},
	})
}
