/*
 * SPDX-FileCopyrightText: The terraform-provider-git Authors
 * SPDX-License-Identifier: 0BSD
 */

package provider_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"os"
	"regexp"
	"testing"
)

func TestResourceGitCommit(t *testing.T) {
	t.Parallel()
	directory, repository := testRepository(t)
	defer os.RemoveAll(directory)
	cfg := testConfig(t, repository)
	worktree := testWorktree(t, repository)
	name := "some-file"
	testWriteFile(t, worktree, name)
	testGitAdd(t, worktree, name)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testProviderFactories(),
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
					resource.TestCheckResourceAttr("git_commit.test", "id", directory),
					resource.TestCheckResourceAttr("git_commit.test", "all", "false"),
					resource.TestCheckResourceAttr("git_commit.test", "message", "committed with terraform"),
					resource.TestCheckResourceAttr("git_commit.test", "author.name", cfg.Author.Name),
					resource.TestCheckResourceAttr("git_commit.test", "author.email", cfg.Author.Email),
					resource.TestCheckResourceAttr("git_commit.test", "committer.name", cfg.Committer.Name),
					resource.TestCheckResourceAttr("git_commit.test", "committer.email", cfg.Committer.Email),
					resource.TestCheckResourceAttrWith("git_commit.test", "sha1", testCheckExactLength(40)),
				),
			},
		},
	})
}

func TestResourceGitCommit_Author_Missing_WithoutConfig(t *testing.T) {
	t.Parallel()
	directory, repository := testRepository(t)
	defer os.RemoveAll(directory)
	worktree := testWorktree(t, repository)
	name := "some-file"
	testWriteFile(t, worktree, name)
	testGitAdd(t, worktree, name)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testProviderFactories(),
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
	directory, repository := testRepository(t)
	defer os.RemoveAll(directory)
	worktree := testWorktree(t, repository)
	name := "some-file"
	testWriteFile(t, worktree, name)
	testGitAdd(t, worktree, name)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testProviderFactories(),
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
	directory, repository := testRepository(t)
	defer os.RemoveAll(directory)
	worktree := testWorktree(t, repository)
	name := "some-file"
	testWriteFile(t, worktree, name)
	testGitAdd(t, worktree, name)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_commit" "test" {
						message = "committed with terraform"
					}
				`),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
		},
	})
}

func TestResourceGitCommit_Author_Partial_Name(t *testing.T) {
	t.Parallel()
	directory, repository := testRepository(t)
	defer os.RemoveAll(directory)
	worktree := testWorktree(t, repository)
	name := "some-file"
	testWriteFile(t, worktree, name)
	testGitAdd(t, worktree, name)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testProviderFactories(),
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
					resource.TestCheckResourceAttr("git_commit.test", "id", directory),
					resource.TestCheckResourceAttr("git_commit.test", "all", "false"),
					resource.TestCheckResourceAttr("git_commit.test", "message", "committed with terraform"),
					resource.TestCheckResourceAttr("git_commit.test", "author.name", "test"),
					resource.TestCheckResourceAttrWith("git_commit.test", "sha1", testCheckExactLength(40)),
				),
			},
		},
	})
}

func TestResourceGitCommit_Author_Partial_Email(t *testing.T) {
	t.Parallel()
	directory, repository := testRepository(t)
	defer os.RemoveAll(directory)
	worktree := testWorktree(t, repository)
	name := "some-file"
	testWriteFile(t, worktree, name)
	testGitAdd(t, worktree, name)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testProviderFactories(),
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
					resource.TestCheckResourceAttr("git_commit.test", "id", directory),
					resource.TestCheckResourceAttr("git_commit.test", "all", "false"),
					resource.TestCheckResourceAttr("git_commit.test", "message", "committed with terraform"),
					resource.TestCheckResourceAttr("git_commit.test", "author.email", "someone@example.com"),
					resource.TestCheckResourceAttrWith("git_commit.test", "sha1", testCheckExactLength(40)),
				),
			},
		},
	})
}

func TestResourceGitCommit_WithoutChanges(t *testing.T) {
	t.Parallel()
	directory, repository := testRepository(t)
	defer os.RemoveAll(directory)
	testConfig(t, repository)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testProviderFactories(),
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
					resource.TestCheckResourceAttr("git_commit.test", "id", directory),
					resource.TestCheckResourceAttr("git_commit.test", "all", "false"),
					resource.TestCheckResourceAttr("git_commit.test", "message", "committed with terraform"),
					resource.TestCheckResourceAttr("git_commit.test", "sha1", ""),
				),
			},
		},
	})
}

func TestResourceGitCommit_WithoutChanges_AllEnabled(t *testing.T) {
	t.Parallel()
	directory, repository := testRepository(t)
	defer os.RemoveAll(directory)
	testConfig(t, repository)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testProviderFactories(),
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
					resource.TestCheckResourceAttr("git_commit.test", "id", directory),
					resource.TestCheckResourceAttr("git_commit.test", "all", "true"),
					resource.TestCheckResourceAttr("git_commit.test", "message", "committed with terraform"),
					resource.TestCheckResourceAttr("git_commit.test", "sha1", ""),
				),
			},
		},
	})
}
