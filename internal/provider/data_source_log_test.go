/*
 * SPDX-FileCopyrightText: The terraform-provider-git Authors
 * SPDX-License-Identifier: 0BSD
 */

package provider_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/metio/terraform-provider-git/internal/testutils"
)

func TestDataSourceGitLog(t *testing.T) {
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
					data "git_log" "test" {
						directory = "%s"
					}
				`, directory),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.git_log.test", "directory", directory),
					resource.TestCheckResourceAttr("data.git_log.test", "id", directory),
					resource.TestCheckResourceAttr("data.git_log.test", "commits.#", "1"),
					resource.TestCheckResourceAttr("data.git_log.test", "commits.0", commit.String()),
				),
			},
		},
	})
}

func TestDataSourceGitLog_All(t *testing.T) {
	t.Parallel()
	directory, repository := testutils.CreateRepository(t)
	testutils.TestConfig(t, repository)
	worktree := testutils.GetRepositoryWorktree(t, repository)
	testutils.AddAndCommitNewFile(t, worktree, "some-file")
	testutils.AddAndCommitNewFile(t, worktree, "other-file")

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "git_log" "test" {
						directory = "%s"
						all       = true
					}
				`, directory),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.git_log.test", "directory", directory),
					resource.TestCheckResourceAttr("data.git_log.test", "id", directory),
					resource.TestCheckResourceAttr("data.git_log.test", "commits.#", "2"),
				),
			},
		},
	})
}

func TestDataSourceGitLog_MaxCount(t *testing.T) {
	t.Parallel()
	directory, repository := testutils.CreateRepository(t)
	testutils.TestConfig(t, repository)
	worktree := testutils.GetRepositoryWorktree(t, repository)
	testutils.AddAndCommitNewFile(t, worktree, "some-file")
	testutils.AddAndCommitNewFile(t, worktree, "other-file")

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "git_log" "test" {
						directory = "%s"
						max_count = 1
					}
				`, directory),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.git_log.test", "directory", directory),
					resource.TestCheckResourceAttr("data.git_log.test", "id", directory),
					resource.TestCheckResourceAttr("data.git_log.test", "commits.#", "1"),
				),
			},
		},
	})
}

func TestDataSourceGitLog_FilterPaths(t *testing.T) {
	t.Parallel()
	directory, repository := testutils.CreateRepository(t)
	testutils.TestConfig(t, repository)
	worktree := testutils.GetRepositoryWorktree(t, repository)
	testutils.AddAndCommitNewFile(t, worktree, "some-file")
	testutils.AddAndCommitNewFile(t, worktree, "other-file")

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "git_log" "test" {
						directory    = "%s"
						filter_paths = ["some-file"]
					}
				`, directory),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.git_log.test", "directory", directory),
					resource.TestCheckResourceAttr("data.git_log.test", "id", directory),
					resource.TestCheckResourceAttr("data.git_log.test", "commits.#", "1"),
				),
			},
		},
	})
}

func TestDataSourceGitLog_From(t *testing.T) {
	t.Parallel()
	directory, repository := testutils.CreateRepository(t)
	testutils.TestConfig(t, repository)
	worktree := testutils.GetRepositoryWorktree(t, repository)
	testutils.AddAndCommitNewFile(t, worktree, "some-file")
	testutils.AddAndCommitNewFile(t, worktree, "other-file")

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "git_log" "test" {
						directory = "%s"
						from      = "HEAD"
					}
				`, directory),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.git_log.test", "directory", directory),
					resource.TestCheckResourceAttr("data.git_log.test", "id", directory),
					resource.TestCheckResourceAttr("data.git_log.test", "commits.#", "2"),
				),
			},
		},
	})
}

func TestDataSourceGitLog_Skip(t *testing.T) {
	t.Parallel()
	directory, repository := testutils.CreateRepository(t)
	testutils.TestConfig(t, repository)
	worktree := testutils.GetRepositoryWorktree(t, repository)
	testutils.AddAndCommitNewFile(t, worktree, "some-file")
	testutils.AddAndCommitNewFile(t, worktree, "other-file")

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "git_log" "test" {
						directory = "%s"
						skip      = 1
					}
				`, directory),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.git_log.test", "directory", directory),
					resource.TestCheckResourceAttr("data.git_log.test", "id", directory),
					resource.TestCheckResourceAttr("data.git_log.test", "commits.#", "1"),
				),
			},
		},
	})
}

func TestDataSourceGitLog_Since(t *testing.T) {
	t.Parallel()
	directory, repository := testutils.CreateRepository(t)
	testutils.TestConfig(t, repository)
	worktree := testutils.GetRepositoryWorktree(t, repository)
	testutils.AddAndCommitNewFile(t, worktree, "some-file")
	testutils.AddAndCommitNewFile(t, worktree, "other-file")

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "git_log" "test" {
						directory = "%s"
						since     = "2017-11-22T00:00:00Z"
					}
				`, directory),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.git_log.test", "directory", directory),
					resource.TestCheckResourceAttr("data.git_log.test", "id", directory),
					resource.TestCheckResourceAttr("data.git_log.test", "commits.#", "2"),
				),
			},
		},
	})
}

func TestDataSourceGitLog_Until(t *testing.T) {
	t.Parallel()
	directory, repository := testutils.CreateRepository(t)
	testutils.TestConfig(t, repository)
	worktree := testutils.GetRepositoryWorktree(t, repository)
	testutils.AddAndCommitNewFile(t, worktree, "some-file")
	testutils.AddAndCommitNewFile(t, worktree, "other-file")

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "git_log" "test" {
						directory = "%s"
						until     = "2017-11-22T00:00:00Z"
					}
				`, directory),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.git_log.test", "directory", directory),
					resource.TestCheckResourceAttr("data.git_log.test", "id", directory),
					resource.TestCheckResourceAttr("data.git_log.test", "commits.#", "0"),
				),
			},
		},
	})
}

func TestDataSourceGitLog_Order_Time(t *testing.T) {
	t.Parallel()
	directory, repository := testutils.CreateRepository(t)
	testutils.TestConfig(t, repository)
	worktree := testutils.GetRepositoryWorktree(t, repository)
	testutils.AddAndCommitNewFile(t, worktree, "some-file")
	testutils.AddAndCommitNewFile(t, worktree, "other-file")

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "git_log" "test" {
						directory = "%s"
						order     = "time"
					}
				`, directory),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.git_log.test", "directory", directory),
					resource.TestCheckResourceAttr("data.git_log.test", "id", directory),
					resource.TestCheckResourceAttr("data.git_log.test", "commits.#", "2"),
				),
			},
		},
	})
}

func TestDataSourceGitLog_Order_Depth(t *testing.T) {
	t.Parallel()
	directory, repository := testutils.CreateRepository(t)
	testutils.TestConfig(t, repository)
	worktree := testutils.GetRepositoryWorktree(t, repository)
	testutils.AddAndCommitNewFile(t, worktree, "some-file")
	testutils.AddAndCommitNewFile(t, worktree, "other-file")

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "git_log" "test" {
						directory = "%s"
						order     = "depth"
					}
				`, directory),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.git_log.test", "directory", directory),
					resource.TestCheckResourceAttr("data.git_log.test", "id", directory),
					resource.TestCheckResourceAttr("data.git_log.test", "commits.#", "2"),
				),
			},
		},
	})
}

func TestDataSourceGitLog_Order_Breadth(t *testing.T) {
	t.Parallel()
	directory, repository := testutils.CreateRepository(t)
	testutils.TestConfig(t, repository)
	worktree := testutils.GetRepositoryWorktree(t, repository)
	testutils.AddAndCommitNewFile(t, worktree, "some-file")
	testutils.AddAndCommitNewFile(t, worktree, "other-file")

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "git_log" "test" {
						directory = "%s"
						order     = "breadth"
					}
				`, directory),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.git_log.test", "directory", directory),
					resource.TestCheckResourceAttr("data.git_log.test", "id", directory),
					resource.TestCheckResourceAttr("data.git_log.test", "commits.#", "2"),
				),
			},
		},
	})
}

func TestDataSourceGitLog_MaxCount_Zero(t *testing.T) {
	t.Parallel()
	directory, repository := testutils.CreateRepository(t)
	testutils.TestConfig(t, repository)
	worktree := testutils.GetRepositoryWorktree(t, repository)
	testutils.AddAndCommitNewFile(t, worktree, "some-file")
	testutils.AddAndCommitNewFile(t, worktree, "other-file")

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "git_log" "test" {
						directory = "%s"
						max_count = 0
					}
				`, directory),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.git_log.test", "directory", directory),
					resource.TestCheckResourceAttr("data.git_log.test", "id", directory),
					resource.TestCheckResourceAttr("data.git_log.test", "commits.#", "0"),
				),
			},
		},
	})
}

func TestDataSourceGitLog_MaxCount_Negative(t *testing.T) {
	t.Parallel()
	directory, repository := testutils.CreateRepository(t)
	testutils.TestConfig(t, repository)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "git_log" "test" {
						directory = "%s"
						max_count = -1
					}
				`, directory),
				ExpectError: regexp.MustCompile(`Invalid Attribute Value`),
			},
		},
	})
}

func TestDataSourceGitLog_Skip_Negative(t *testing.T) {
	t.Parallel()
	directory, repository := testutils.CreateRepository(t)
	testutils.TestConfig(t, repository)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "git_log" "test" {
						directory = "%s"
						skip      = -1
					}
				`, directory),
				ExpectError: regexp.MustCompile(`Invalid Attribute Value`),
			},
		},
	})
}

func TestDataSourceGitLog_Order_Invalid(t *testing.T) {
	t.Parallel()
	directory, repository := testutils.CreateRepository(t)
	testutils.TestConfig(t, repository)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "git_log" "test" {
						directory = "%s"
						order     = "something-else"
					}
				`, directory),
				ExpectError: regexp.MustCompile(`Invalid Attribute Value Match`),
			},
		},
	})
}
