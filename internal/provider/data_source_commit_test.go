/*
 * SPDX-FileCopyrightText: The terraform-provider-git Authors
 * SPDX-License-Identifier: 0BSD
 */

package provider_test

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/metio/terraform-provider-git/internal/testutils"
	"testing"
)

func TestDataSourceGitCommit(t *testing.T) {
	t.Parallel()
	directory, repository := testutils.CreateRepository(t)
	testutils.TestConfig(t, repository)
	worktree := testutils.GetRepositoryWorktree(t, repository)
	fileName := "some-file"
	testutils.WriteFileInWorktree(t, worktree, fileName)
	testutils.GitAdd(t, worktree, fileName)
	commit := testutils.GitCommit(t, worktree)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "git_commit" "test" {
						directory = "%s"
						revision  = "%s"
					}
				`, directory, commit.String()),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.git_commit.test", "directory", directory),
					resource.TestCheckResourceAttr("data.git_commit.test", "id", commit.String()),
					resource.TestCheckResourceAttr("data.git_commit.test", "revision", commit.String()),
					resource.TestCheckResourceAttr("data.git_commit.test", "sha1", commit.String()),
					resource.TestCheckResourceAttr("data.git_commit.test", "signature", ""),
					resource.TestCheckResourceAttr("data.git_commit.test", "files.#", "1"),
					resource.TestCheckResourceAttrWith("data.git_commit.test", "message", testutils.CheckMinLength(1)),
					resource.TestCheckResourceAttrWith("data.git_commit.test", "tree_sha1", testutils.CheckMinLength(1)),
					resource.TestCheckResourceAttrWith("data.git_commit.test", "author.name", testutils.CheckMinLength(1)),
					resource.TestCheckResourceAttrWith("data.git_commit.test", "author.email", testutils.CheckMinLength(1)),
					resource.TestCheckResourceAttrWith("data.git_commit.test", "author.timestamp", testutils.CheckMinLength(1)),
					resource.TestCheckResourceAttrWith("data.git_commit.test", "committer.name", testutils.CheckMinLength(1)),
					resource.TestCheckResourceAttrWith("data.git_commit.test", "committer.email", testutils.CheckMinLength(1)),
					resource.TestCheckResourceAttrWith("data.git_commit.test", "committer.timestamp", testutils.CheckMinLength(1)),
				),
			},
		},
	})
}

func TestDataSourceGitCommit_MultipleFiles(t *testing.T) {
	t.Parallel()
	directory, repository := testutils.CreateRepository(t)
	testutils.TestConfig(t, repository)
	worktree := testutils.GetRepositoryWorktree(t, repository)
	fileName1 := "some-file"
	fileName2 := "other-file"
	testutils.WriteFileInWorktree(t, worktree, fileName1)
	testutils.WriteFileInWorktree(t, worktree, fileName2)
	testutils.GitAdd(t, worktree, fileName1)
	testutils.GitAdd(t, worktree, fileName2)
	commit := testutils.GitCommit(t, worktree)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "git_commit" "test" {
						directory = "%s"
						revision  = "%s"
					}
				`, directory, commit.String()),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.git_commit.test", "directory", directory),
					resource.TestCheckResourceAttr("data.git_commit.test", "id", commit.String()),
					resource.TestCheckResourceAttr("data.git_commit.test", "revision", commit.String()),
					resource.TestCheckResourceAttr("data.git_commit.test", "sha1", commit.String()),
					resource.TestCheckResourceAttr("data.git_commit.test", "signature", ""),
					resource.TestCheckResourceAttr("data.git_commit.test", "files.#", "2"),
					resource.TestCheckResourceAttrWith("data.git_commit.test", "message", testutils.CheckMinLength(1)),
					resource.TestCheckResourceAttrWith("data.git_commit.test", "tree_sha1", testutils.CheckMinLength(1)),
					resource.TestCheckResourceAttrWith("data.git_commit.test", "author.name", testutils.CheckMinLength(1)),
					resource.TestCheckResourceAttrWith("data.git_commit.test", "author.email", testutils.CheckMinLength(1)),
					resource.TestCheckResourceAttrWith("data.git_commit.test", "author.timestamp", testutils.CheckMinLength(1)),
					resource.TestCheckResourceAttrWith("data.git_commit.test", "committer.name", testutils.CheckMinLength(1)),
					resource.TestCheckResourceAttrWith("data.git_commit.test", "committer.email", testutils.CheckMinLength(1)),
					resource.TestCheckResourceAttrWith("data.git_commit.test", "committer.timestamp", testutils.CheckMinLength(1)),
				),
			},
		},
	})
}

func TestDataSourceGitCommit_WithHead(t *testing.T) {
	t.Parallel()
	directory, repository := testutils.CreateRepository(t)
	testutils.TestConfig(t, repository)
	worktree := testutils.GetRepositoryWorktree(t, repository)
	fileName := "some-file"
	testutils.WriteFileInWorktree(t, worktree, fileName)
	testutils.GitAdd(t, worktree, fileName)
	commit := testutils.GitCommit(t, worktree)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "git_commit" "test" {
						directory = "%s"
						revision  = "HEAD"
					}
				`, directory),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.git_commit.test", "directory", directory),
					resource.TestCheckResourceAttr("data.git_commit.test", "id", "HEAD"),
					resource.TestCheckResourceAttr("data.git_commit.test", "revision", "HEAD"),
					resource.TestCheckResourceAttr("data.git_commit.test", "sha1", commit.String()),
					resource.TestCheckResourceAttr("data.git_commit.test", "files.#", "1"),
					resource.TestCheckResourceAttr("data.git_commit.test", "signature", ""),
					resource.TestCheckResourceAttrWith("data.git_commit.test", "message", testutils.CheckMinLength(1)),
					resource.TestCheckResourceAttrWith("data.git_commit.test", "tree_sha1", testutils.CheckMinLength(1)),
					resource.TestCheckResourceAttrWith("data.git_commit.test", "author.name", testutils.CheckMinLength(1)),
					resource.TestCheckResourceAttrWith("data.git_commit.test", "author.email", testutils.CheckMinLength(1)),
					resource.TestCheckResourceAttrWith("data.git_commit.test", "author.timestamp", testutils.CheckMinLength(1)),
					resource.TestCheckResourceAttrWith("data.git_commit.test", "committer.name", testutils.CheckMinLength(1)),
					resource.TestCheckResourceAttrWith("data.git_commit.test", "committer.email", testutils.CheckMinLength(1)),
					resource.TestCheckResourceAttrWith("data.git_commit.test", "committer.timestamp", testutils.CheckMinLength(1)),
				),
			},
		},
	})
}

func TestDataSourceGitCommit_WithBranch(t *testing.T) {
	t.Parallel()
	directory, repository := testutils.CreateRepository(t)
	testutils.TestConfig(t, repository)
	worktree := testutils.GetRepositoryWorktree(t, repository)
	fileName := "some-file"
	testutils.WriteFileInWorktree(t, worktree, fileName)
	testutils.GitAdd(t, worktree, fileName)
	commit := testutils.GitCommit(t, worktree)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "git_commit" "test" {
						directory = "%s"
						revision  = "master"
					}
				`, directory),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.git_commit.test", "directory", directory),
					resource.TestCheckResourceAttr("data.git_commit.test", "id", "master"),
					resource.TestCheckResourceAttr("data.git_commit.test", "revision", "master"),
					resource.TestCheckResourceAttr("data.git_commit.test", "sha1", commit.String()),
					resource.TestCheckResourceAttr("data.git_commit.test", "signature", ""),
					resource.TestCheckResourceAttrWith("data.git_commit.test", "message", testutils.CheckMinLength(1)),
					resource.TestCheckResourceAttrWith("data.git_commit.test", "tree_sha1", testutils.CheckMinLength(1)),
					resource.TestCheckResourceAttrWith("data.git_commit.test", "author.name", testutils.CheckMinLength(1)),
					resource.TestCheckResourceAttrWith("data.git_commit.test", "author.email", testutils.CheckMinLength(1)),
					resource.TestCheckResourceAttrWith("data.git_commit.test", "author.timestamp", testutils.CheckMinLength(1)),
					resource.TestCheckResourceAttrWith("data.git_commit.test", "committer.name", testutils.CheckMinLength(1)),
					resource.TestCheckResourceAttrWith("data.git_commit.test", "committer.email", testutils.CheckMinLength(1)),
					resource.TestCheckResourceAttrWith("data.git_commit.test", "committer.timestamp", testutils.CheckMinLength(1)),
				),
			},
		},
	})
}

func TestDataSourceGitCommit_WithSignature(t *testing.T) {
	t.Parallel()
	directory, repository := testutils.CreateRepository(t)
	worktree := testutils.GetRepositoryWorktree(t, repository)
	fileName := "some-file"
	testutils.WriteFileInWorktree(t, worktree, fileName)
	testutils.GitAdd(t, worktree, fileName)
	signature := testutils.Signature()
	commit := testutils.GitCommitWith(t, worktree, &git.CommitOptions{
		Author: signature,
	})

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "git_commit" "test" {
						directory = "%s"
						revision  = "%s"
					}
				`, directory, commit.String()),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.git_commit.test", "directory", directory),
					resource.TestCheckResourceAttr("data.git_commit.test", "id", commit.String()),
					resource.TestCheckResourceAttr("data.git_commit.test", "revision", commit.String()),
					resource.TestCheckResourceAttr("data.git_commit.test", "sha1", commit.String()),
					resource.TestCheckResourceAttr("data.git_commit.test", "signature", ""),
					resource.TestCheckResourceAttrWith("data.git_commit.test", "message", testutils.CheckMinLength(1)),
					resource.TestCheckResourceAttrWith("data.git_commit.test", "tree_sha1", testutils.CheckMinLength(1)),
					resource.TestCheckResourceAttr("data.git_commit.test", "author.name", signature.Name),
					resource.TestCheckResourceAttr("data.git_commit.test", "author.email", signature.Email),
					resource.TestCheckResourceAttrWith("data.git_commit.test", "author.timestamp", testutils.CheckMinLength(1)),
					resource.TestCheckResourceAttr("data.git_commit.test", "committer.name", signature.Name),
					resource.TestCheckResourceAttr("data.git_commit.test", "committer.email", signature.Email),
					resource.TestCheckResourceAttrWith("data.git_commit.test", "committer.timestamp", testutils.CheckMinLength(1)),
				),
			},
		},
	})
}
