/*
 * SPDX-FileCopyrightText: The terraform-provider-git Authors
 * SPDX-License-Identifier: 0BSD
 */

package provider_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/metio/terraform-provider-git/internal/testutils"
	"os"
	"regexp"
	"testing"
)

func TestResourceGitCommit(t *testing.T) {
	t.Parallel()
	directory, repository := testutils.CreateRepository(t)
	defer os.RemoveAll(directory)
	cfg := testutils.TestConfig(t, repository)
	worktree := testutils.GetRepositoryWorktree(t, repository)
	name := "some-file"
	testutils.WriteFileInWorktree(t, worktree, name)
	testutils.GitAdd(t, worktree, name)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_commit" "test" {
						directory = "%s"
						message   = "committed with terraform"
					}
				`, directory),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("git_commit.test", "directory", directory),
					resource.TestCheckResourceAttrWith("git_commit.test", "id", testutils.CheckMinLength(1)),
					resource.TestCheckResourceAttr("git_commit.test", "all", "false"),
					resource.TestCheckResourceAttr("git_commit.test", "allow_empty_commits", "true"),
					resource.TestCheckResourceAttr("git_commit.test", "message", "committed with terraform"),
					resource.TestCheckResourceAttr("git_commit.test", "author.name", cfg.Author.Name),
					resource.TestCheckResourceAttr("git_commit.test", "author.email", cfg.Author.Email),
					resource.TestCheckResourceAttr("git_commit.test", "committer.name", cfg.Committer.Name),
					resource.TestCheckResourceAttr("git_commit.test", "committer.email", cfg.Committer.Email),
					resource.TestCheckResourceAttrWith("git_commit.test", "sha1", testutils.CheckExactLength(40)),
					resource.TestCheckResourceAttr("git_commit.test", "files.#", "1"),
					resource.TestCheckResourceAttr("git_commit.test", "files.0", name),
				),
			},
		},
	})
}

func TestResourceGitCommit_AllowEmptyCommit_True(t *testing.T) {
	t.Parallel()
	directory, repository := testutils.CreateRepository(t)
	defer os.RemoveAll(directory)
	cfg := testutils.TestConfig(t, repository)
	worktree := testutils.GetRepositoryWorktree(t, repository)
	testutils.WriteFileInWorktree(t, worktree, "some-file")

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_commit" "test" {
						directory = "%s"
						message   = "committed with terraform"
						allow_empty_commits = true
					}
				`, directory),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("git_commit.test", "directory", directory),
					resource.TestCheckResourceAttrWith("git_commit.test", "id", testutils.CheckMinLength(1)),
					resource.TestCheckResourceAttr("git_commit.test", "all", "false"),
					resource.TestCheckResourceAttr("git_commit.test", "allow_empty_commits", "true"),
					resource.TestCheckResourceAttr("git_commit.test", "message", "committed with terraform"),
					resource.TestCheckResourceAttr("git_commit.test", "author.name", cfg.Author.Name),
					resource.TestCheckResourceAttr("git_commit.test", "author.email", cfg.Author.Email),
					resource.TestCheckResourceAttr("git_commit.test", "committer.name", cfg.Committer.Name),
					resource.TestCheckResourceAttr("git_commit.test", "committer.email", cfg.Committer.Email),
					resource.TestCheckResourceAttrWith("git_commit.test", "sha1", testutils.CheckExactLength(40)),
					resource.TestCheckResourceAttr("git_commit.test", "files.#", "0"),
				),
			},
		},
	})
}

func TestResourceGitCommit_AllowEmptyCommit_False(t *testing.T) {
	t.Parallel()
	directory, repository := testutils.CreateRepository(t)
	defer os.RemoveAll(directory)
	testutils.TestConfig(t, repository)
	worktree := testutils.GetRepositoryWorktree(t, repository)
	testutils.WriteFileInWorktree(t, worktree, "some-file")

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_commit" "test" {
						directory = "%s"
						message   = "committed with terraform"
						allow_empty_commits = false
					}
				`, directory),
				ExpectError: regexp.MustCompile(`cannot create empty commit`),
			},
		},
	})
}

func TestResourceGitCommit_Author_Missing_WithoutConfig(t *testing.T) {
	t.Parallel()
	directory, repository := testutils.CreateRepository(t)
	defer os.RemoveAll(directory)
	worktree := testutils.GetRepositoryWorktree(t, repository)
	name := "some-file"
	testutils.WriteFileInWorktree(t, worktree, name)
	testutils.GitAdd(t, worktree, name)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_commit" "test" {
						directory = "%s"
						message   = "committed with terraform"
					}
				`, directory),
				ExpectError: regexp.MustCompile(`Cannot create commit`),
			},
		},
	})
}

func TestResourceGitCommit_Message_Missing(t *testing.T) {
	t.Parallel()
	directory, repository := testutils.CreateRepository(t)
	defer os.RemoveAll(directory)
	worktree := testutils.GetRepositoryWorktree(t, repository)
	name := "some-file"
	testutils.WriteFileInWorktree(t, worktree, name)
	testutils.GitAdd(t, worktree, name)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_commit" "test" {
						directory = "%s"
					}
				`, directory),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
		},
	})
}

func TestResourceGitCommit_Directory_Missing(t *testing.T) {
	t.Parallel()
	directory, repository := testutils.CreateRepository(t)
	defer os.RemoveAll(directory)
	worktree := testutils.GetRepositoryWorktree(t, repository)
	name := "some-file"
	testutils.WriteFileInWorktree(t, worktree, name)
	testutils.GitAdd(t, worktree, name)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: `
					resource "git_commit" "test" {
						message = "committed with terraform"
					}
				`,
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
		},
	})
}

func TestResourceGitCommit_Author_Partial_Name(t *testing.T) {
	t.Parallel()
	directory, repository := testutils.CreateRepository(t)
	defer os.RemoveAll(directory)
	worktree := testutils.GetRepositoryWorktree(t, repository)
	name := "some-file"
	testutils.WriteFileInWorktree(t, worktree, name)
	testutils.GitAdd(t, worktree, name)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_commit" "test" {
						directory = "%s"
						message   = "committed with terraform"
						author    = {
							name = "test"
						}
					}
				`, directory),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("git_commit.test", "directory", directory),
					resource.TestCheckResourceAttrWith("git_commit.test", "id", testutils.CheckMinLength(1)),
					resource.TestCheckResourceAttr("git_commit.test", "all", "false"),
					resource.TestCheckResourceAttr("git_commit.test", "message", "committed with terraform"),
					resource.TestCheckResourceAttr("git_commit.test", "author.name", "test"),
					resource.TestCheckResourceAttr("git_commit.test", "author.email", ""),
					resource.TestCheckResourceAttr("git_commit.test", "committer.name", "test"),
					resource.TestCheckResourceAttr("git_commit.test", "committer.email", ""),
					resource.TestCheckResourceAttrWith("git_commit.test", "sha1", testutils.CheckExactLength(40)),
					resource.TestCheckResourceAttr("git_commit.test", "files.#", "1"),
				),
			},
		},
	})
}

func TestResourceGitCommit_Author_Name_Update(t *testing.T) {
	t.Parallel()
	directory, repository := testutils.CreateRepository(t)
	defer os.RemoveAll(directory)
	worktree := testutils.GetRepositoryWorktree(t, repository)
	name := "some-file"
	testutils.WriteFileInWorktree(t, worktree, name)
	testutils.GitAdd(t, worktree, name)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_commit" "test" {
						directory = "%s"
						message   = "committed with terraform"
						author    = {
							name = "test"
						}
					}
				`, directory),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("git_commit.test", "directory", directory),
					resource.TestCheckResourceAttrWith("git_commit.test", "id", testutils.CheckMinLength(1)),
					resource.TestCheckResourceAttr("git_commit.test", "all", "false"),
					resource.TestCheckResourceAttr("git_commit.test", "message", "committed with terraform"),
					resource.TestCheckResourceAttr("git_commit.test", "author.name", "test"),
					resource.TestCheckResourceAttr("git_commit.test", "author.email", ""),
					resource.TestCheckResourceAttr("git_commit.test", "committer.name", "test"),
					resource.TestCheckResourceAttr("git_commit.test", "committer.email", ""),
					resource.TestCheckResourceAttrWith("git_commit.test", "sha1", testutils.CheckExactLength(40)),
					resource.TestCheckResourceAttr("git_commit.test", "files.#", "1"),
				),
			},
			{
				Config: fmt.Sprintf(`
					resource "git_commit" "test" {
						directory = "%s"
						message   = "committed with terraform"
						author    = {
							name = "different"
						}
					}
				`, directory),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("git_commit.test", "directory", directory),
					resource.TestCheckResourceAttrWith("git_commit.test", "id", testutils.CheckMinLength(1)),
					resource.TestCheckResourceAttr("git_commit.test", "all", "false"),
					resource.TestCheckResourceAttr("git_commit.test", "message", "committed with terraform"),
					resource.TestCheckResourceAttr("git_commit.test", "author.name", "different"),
					resource.TestCheckNoResourceAttr("git_commit.test", "author.email"),
					resource.TestCheckNoResourceAttr("git_commit.test", "committer.name"),
					resource.TestCheckNoResourceAttr("git_commit.test", "committer.email"),
					resource.TestCheckNoResourceAttr("git_commit.test", "sha1"),
					resource.TestCheckResourceAttr("git_commit.test", "files.#", "0"),
				),
			},
		},
	})
}

func TestResourceGitCommit_Author_Partial_Email(t *testing.T) {
	t.Parallel()
	directory, repository := testutils.CreateRepository(t)
	defer os.RemoveAll(directory)
	worktree := testutils.GetRepositoryWorktree(t, repository)
	name := "some-file"
	testutils.WriteFileInWorktree(t, worktree, name)
	testutils.GitAdd(t, worktree, name)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_commit" "test" {
						directory = "%s"
						message   = "committed with terraform"
						author    = {
							email = "someone@example.com"
						}
					}
				`, directory),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("git_commit.test", "directory", directory),
					resource.TestCheckResourceAttrWith("git_commit.test", "id", testutils.CheckMinLength(1)),
					resource.TestCheckResourceAttr("git_commit.test", "all", "false"),
					resource.TestCheckResourceAttr("git_commit.test", "message", "committed with terraform"),
					resource.TestCheckResourceAttr("git_commit.test", "author.name", ""),
					resource.TestCheckResourceAttr("git_commit.test", "author.email", "someone@example.com"),
					resource.TestCheckResourceAttr("git_commit.test", "committer.name", ""),
					resource.TestCheckResourceAttr("git_commit.test", "committer.email", "someone@example.com"),
					resource.TestCheckResourceAttrWith("git_commit.test", "sha1", testutils.CheckExactLength(40)),
					resource.TestCheckResourceAttr("git_commit.test", "files.#", "1"),
				),
			},
		},
	})
}

func TestResourceGitCommit_Author_Email_Update(t *testing.T) {
	t.Parallel()
	directory, repository := testutils.CreateRepository(t)
	defer os.RemoveAll(directory)
	worktree := testutils.GetRepositoryWorktree(t, repository)
	name := "some-file"
	testutils.WriteFileInWorktree(t, worktree, name)
	testutils.GitAdd(t, worktree, name)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_commit" "test" {
						directory = "%s"
						message   = "committed with terraform"
						author    = {
							email = "someone@example.com"
						}
					}
				`, directory),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("git_commit.test", "directory", directory),
					resource.TestCheckResourceAttrWith("git_commit.test", "id", testutils.CheckMinLength(1)),
					resource.TestCheckResourceAttr("git_commit.test", "all", "false"),
					resource.TestCheckResourceAttr("git_commit.test", "message", "committed with terraform"),
					resource.TestCheckResourceAttr("git_commit.test", "author.name", ""),
					resource.TestCheckResourceAttr("git_commit.test", "author.email", "someone@example.com"),
					resource.TestCheckResourceAttr("git_commit.test", "committer.name", ""),
					resource.TestCheckResourceAttr("git_commit.test", "committer.email", "someone@example.com"),
					resource.TestCheckResourceAttrWith("git_commit.test", "sha1", testutils.CheckExactLength(40)),
					resource.TestCheckResourceAttr("git_commit.test", "files.#", "1"),
				),
			},
			{
				Config: fmt.Sprintf(`
					resource "git_commit" "test" {
						directory = "%s"
						message   = "committed with terraform"
						author    = {
							email = "different@example.com"
						}
					}
				`, directory),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("git_commit.test", "directory", directory),
					resource.TestCheckResourceAttrWith("git_commit.test", "id", testutils.CheckMinLength(1)),
					resource.TestCheckResourceAttr("git_commit.test", "all", "false"),
					resource.TestCheckResourceAttr("git_commit.test", "message", "committed with terraform"),
					resource.TestCheckNoResourceAttr("git_commit.test", "author.name"),
					resource.TestCheckResourceAttr("git_commit.test", "author.email", "different@example.com"),
					resource.TestCheckNoResourceAttr("git_commit.test", "committer.name"),
					resource.TestCheckNoResourceAttr("git_commit.test", "committer.email"),
					resource.TestCheckNoResourceAttr("git_commit.test", "sha1"),
					resource.TestCheckResourceAttr("git_commit.test", "files.#", "0"),
				),
			},
		},
	})
}

func TestResourceGitCommit_Committer_Partial_Name(t *testing.T) {
	t.Parallel()
	directory, repository := testutils.CreateRepository(t)
	defer os.RemoveAll(directory)
	cfg := testutils.TestConfig(t, repository)
	worktree := testutils.GetRepositoryWorktree(t, repository)
	name := "some-file"
	testutils.WriteFileInWorktree(t, worktree, name)
	testutils.GitAdd(t, worktree, name)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_commit" "test" {
						directory = "%s"
						message   = "committed with terraform"
						committer    = {
							name = "test"
						}
					}
				`, directory),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("git_commit.test", "directory", directory),
					resource.TestCheckResourceAttrWith("git_commit.test", "id", testutils.CheckMinLength(1)),
					resource.TestCheckResourceAttr("git_commit.test", "all", "false"),
					resource.TestCheckResourceAttr("git_commit.test", "message", "committed with terraform"),
					resource.TestCheckResourceAttr("git_commit.test", "author.name", cfg.Author.Name),
					resource.TestCheckResourceAttr("git_commit.test", "author.email", cfg.Author.Email),
					resource.TestCheckResourceAttr("git_commit.test", "committer.name", "test"),
					resource.TestCheckResourceAttr("git_commit.test", "committer.email", ""),
					resource.TestCheckResourceAttrWith("git_commit.test", "sha1", testutils.CheckExactLength(40)),
					resource.TestCheckResourceAttr("git_commit.test", "files.#", "1"),
				),
			},
		},
	})
}

func TestResourceGitCommit_Committer_Partial_Email(t *testing.T) {
	t.Parallel()
	directory, repository := testutils.CreateRepository(t)
	defer os.RemoveAll(directory)
	cfg := testutils.TestConfig(t, repository)
	worktree := testutils.GetRepositoryWorktree(t, repository)
	name := "some-file"
	testutils.WriteFileInWorktree(t, worktree, name)
	testutils.GitAdd(t, worktree, name)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_commit" "test" {
						directory = "%s"
						message   = "committed with terraform"
						committer    = {
							email = "someone@example.com"
						}
					}
				`, directory),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("git_commit.test", "directory", directory),
					resource.TestCheckResourceAttrWith("git_commit.test", "id", testutils.CheckMinLength(1)),
					resource.TestCheckResourceAttr("git_commit.test", "all", "false"),
					resource.TestCheckResourceAttr("git_commit.test", "message", "committed with terraform"),
					resource.TestCheckResourceAttr("git_commit.test", "author.name", cfg.Author.Name),
					resource.TestCheckResourceAttr("git_commit.test", "author.email", cfg.Author.Email),
					resource.TestCheckResourceAttr("git_commit.test", "committer.name", ""),
					resource.TestCheckResourceAttr("git_commit.test", "committer.email", "someone@example.com"),
					resource.TestCheckResourceAttrWith("git_commit.test", "sha1", testutils.CheckExactLength(40)),
					resource.TestCheckResourceAttr("git_commit.test", "files.#", "1"),
				),
			},
		},
	})
}

func TestResourceGitCommit_Committer_Name_Update(t *testing.T) {
	t.Parallel()
	directory, repository := testutils.CreateRepository(t)
	defer os.RemoveAll(directory)
	testutils.TestConfig(t, repository)
	worktree := testutils.GetRepositoryWorktree(t, repository)
	name := "some-file"
	testutils.WriteFileInWorktree(t, worktree, name)
	testutils.GitAdd(t, worktree, name)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_commit" "test" {
						directory = "%s"
						message   = "committed with terraform"
						committer    = {
							name = "test"
						}
					}
				`, directory),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("git_commit.test", "directory", directory),
					resource.TestCheckResourceAttrWith("git_commit.test", "id", testutils.CheckMinLength(1)),
					resource.TestCheckResourceAttr("git_commit.test", "all", "false"),
					resource.TestCheckResourceAttr("git_commit.test", "message", "committed with terraform"),
					resource.TestCheckResourceAttr("git_commit.test", "committer.name", "test"),
					resource.TestCheckResourceAttr("git_commit.test", "committer.email", ""),
					resource.TestCheckResourceAttrWith("git_commit.test", "sha1", testutils.CheckExactLength(40)),
					resource.TestCheckResourceAttr("git_commit.test", "files.#", "1"),
				),
			},
			{
				Config: fmt.Sprintf(`
					resource "git_commit" "test" {
						directory = "%s"
						message   = "committed with terraform"
						committer    = {
							name = "different"
						}
					}
				`, directory),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("git_commit.test", "directory", directory),
					resource.TestCheckResourceAttrWith("git_commit.test", "id", testutils.CheckMinLength(1)),
					resource.TestCheckResourceAttr("git_commit.test", "all", "false"),
					resource.TestCheckResourceAttr("git_commit.test", "message", "committed with terraform"),
					resource.TestCheckResourceAttr("git_commit.test", "committer.name", "different"),
					resource.TestCheckNoResourceAttr("git_commit.test", "committer.email"),
					resource.TestCheckNoResourceAttr("git_commit.test", "sha1"),
					resource.TestCheckResourceAttr("git_commit.test", "files.#", "0"),
				),
			},
		},
	})
}

func TestResourceGitCommit_Committer_Email_Update(t *testing.T) {
	t.Parallel()
	directory, repository := testutils.CreateRepository(t)
	defer os.RemoveAll(directory)
	testutils.TestConfig(t, repository)
	worktree := testutils.GetRepositoryWorktree(t, repository)
	name := "some-file"
	testutils.WriteFileInWorktree(t, worktree, name)
	testutils.GitAdd(t, worktree, name)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_commit" "test" {
						directory = "%s"
						message   = "committed with terraform"
						committer    = {
							email = "someone@example.com"
						}
					}
				`, directory),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("git_commit.test", "directory", directory),
					resource.TestCheckResourceAttrWith("git_commit.test", "id", testutils.CheckMinLength(1)),
					resource.TestCheckResourceAttr("git_commit.test", "all", "false"),
					resource.TestCheckResourceAttr("git_commit.test", "message", "committed with terraform"),
					resource.TestCheckResourceAttr("git_commit.test", "committer.name", ""),
					resource.TestCheckResourceAttr("git_commit.test", "committer.email", "someone@example.com"),
					resource.TestCheckResourceAttrWith("git_commit.test", "sha1", testutils.CheckExactLength(40)),
					resource.TestCheckResourceAttr("git_commit.test", "files.#", "1"),
				),
			},
			{
				Config: fmt.Sprintf(`
					resource "git_commit" "test" {
						directory = "%s"
						message   = "committed with terraform"
						committer    = {
							email = "different@example.com"
						}
					}
				`, directory),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("git_commit.test", "directory", directory),
					resource.TestCheckResourceAttrWith("git_commit.test", "id", testutils.CheckMinLength(1)),
					resource.TestCheckResourceAttr("git_commit.test", "all", "false"),
					resource.TestCheckResourceAttr("git_commit.test", "message", "committed with terraform"),
					resource.TestCheckNoResourceAttr("git_commit.test", "committer.name"),
					resource.TestCheckResourceAttr("git_commit.test", "committer.email", "different@example.com"),
					resource.TestCheckNoResourceAttr("git_commit.test", "sha1"),
					resource.TestCheckResourceAttr("git_commit.test", "files.#", "0"),
				),
			},
		},
	})
}

func TestResourceGitCommit_Committer_WithoutAuthor(t *testing.T) {
	t.Parallel()
	directory, repository := testutils.CreateRepository(t)
	defer os.RemoveAll(directory)
	worktree := testutils.GetRepositoryWorktree(t, repository)
	name := "some-file"
	testutils.WriteFileInWorktree(t, worktree, name)
	testutils.GitAdd(t, worktree, name)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_commit" "test" {
						directory = "%s"
						message   = "committed with terraform"
						committer    = {
							email = "someone@example.com"
						}
					}
				`, directory),
				ExpectError: regexp.MustCompile(`author field is required`),
			},
		},
	})
}

func TestResourceGitCommit_WithoutChanges(t *testing.T) {
	t.Parallel()
	directory, repository := testutils.CreateRepository(t)
	defer os.RemoveAll(directory)
	testutils.TestConfig(t, repository)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_commit" "test" {
						directory = "%s"
						message   = "committed with terraform"
					}
				`, directory),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("git_commit.test", "directory", directory),
					resource.TestCheckResourceAttrWith("git_commit.test", "id", testutils.CheckMinLength(1)),
					resource.TestCheckResourceAttr("git_commit.test", "all", "false"),
					resource.TestCheckResourceAttr("git_commit.test", "message", "committed with terraform"),
					resource.TestCheckNoResourceAttr("git_commit.test", "sha1"),
				),
			},
		},
	})
}

func TestResourceGitCommit_WithoutChanges_AllEnabled(t *testing.T) {
	t.Parallel()
	directory, repository := testutils.CreateRepository(t)
	defer os.RemoveAll(directory)
	testutils.TestConfig(t, repository)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_commit" "test" {
						directory = "%s"
						message   = "committed with terraform"
						all       = true
					}
				`, directory),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("git_commit.test", "directory", directory),
					resource.TestCheckResourceAttrWith("git_commit.test", "id", testutils.CheckMinLength(1)),
					resource.TestCheckResourceAttr("git_commit.test", "all", "true"),
					resource.TestCheckResourceAttr("git_commit.test", "message", "committed with terraform"),
					resource.TestCheckNoResourceAttr("git_commit.test", "sha1"),
				),
			},
		},
	})
}

func TestResourceGitCommit_WithUnstagedChanges_AllEnabled(t *testing.T) {
	t.Parallel()
	directory, repository := testutils.CreateRepository(t)
	defer os.RemoveAll(directory)
	testutils.TestConfig(t, repository)
	worktree := testutils.GetRepositoryWorktree(t, repository)
	testutils.WriteFileInWorktree(t, worktree, "some-file")
	testutils.WriteFileInWorktree(t, worktree, "other-file")

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_commit" "test" {
						directory = "%s"
						message   = "committed with terraform"
						all       = true
					}
				`, directory),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("git_commit.test", "directory", directory),
					resource.TestCheckResourceAttrWith("git_commit.test", "id", testutils.CheckMinLength(1)),
					resource.TestCheckResourceAttr("git_commit.test", "all", "true"),
					resource.TestCheckResourceAttr("git_commit.test", "message", "committed with terraform"),
					resource.TestCheckResourceAttrWith("git_commit.test", "sha1", testutils.CheckExactLength(40)),
					resource.TestCheckResourceAttr("git_commit.test", "files.#", "0"),
				),
			},
		},
	})
}

func TestResourceGitCommit_WithStagedChanges_AllEnabled(t *testing.T) {
	t.Parallel()
	directory, repository := testutils.CreateRepository(t)
	defer os.RemoveAll(directory)
	testutils.TestConfig(t, repository)
	worktree := testutils.GetRepositoryWorktree(t, repository)
	testutils.WriteFileInWorktree(t, worktree, "some-file")
	testutils.WriteFileInWorktree(t, worktree, "other-file")
	testutils.GitAdd(t, worktree, "some-file")
	testutils.GitAdd(t, worktree, "other-file")

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_commit" "test" {
						directory = "%s"
						message   = "committed with terraform"
						all       = true
					}
				`, directory),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("git_commit.test", "directory", directory),
					resource.TestCheckResourceAttrWith("git_commit.test", "id", testutils.CheckMinLength(1)),
					resource.TestCheckResourceAttr("git_commit.test", "all", "true"),
					resource.TestCheckResourceAttr("git_commit.test", "message", "committed with terraform"),
					resource.TestCheckResourceAttrWith("git_commit.test", "sha1", testutils.CheckExactLength(40)),
					resource.TestCheckResourceAttr("git_commit.test", "files.#", "2"),
				),
			},
		},
	})
}
