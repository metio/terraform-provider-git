/*
 * SPDX-FileCopyrightText: Sebastian Hoß <seb@hoß.de>
 * SPDX-License-Identifier: BSD0
 */

package provider_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"os"
	"regexp"
	"testing"
)

func TestResourceGitTag(t *testing.T) {
	t.Parallel()
	directory, repository := testRepository(t)
	defer os.RemoveAll(directory)
	testConfig(t, repository)
	worktree := testWorktree(t, repository)
	testAddAndCommitNewFile(t, worktree, "some-file")
	tag := "some-name"

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_tag" "test" {
						directory = "%s"
						name      = "%s"
					}
				`, directory, tag),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("git_tag.test", "directory", directory),
					resource.TestCheckResourceAttr("git_tag.test", "id", tag),
					resource.TestCheckResourceAttr("git_tag.test", "name", tag),
				),
			},
		},
	})
}

func TestResourceGitTag_Annotated(t *testing.T) {
	t.Parallel()
	directory, repository := testRepository(t)
	defer os.RemoveAll(directory)
	testConfig(t, repository)
	worktree := testWorktree(t, repository)
	testAddAndCommitNewFile(t, worktree, "some-file")
	tag := "some-name"
	message := "some message for the tag"

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_tag" "test" {
						directory = "%s"
						name      = "%s"
						message   = "some message for the tag"
					}
				`, directory, tag),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("git_tag.test", "directory", directory),
					resource.TestCheckResourceAttr("git_tag.test", "id", tag),
					resource.TestCheckResourceAttr("git_tag.test", "name", tag),
					resource.TestCheckResourceAttr("git_tag.test", "message", message),
				),
			},
		},
	})
}

func TestResourceGitTag_Commitish(t *testing.T) {
	t.Parallel()
	directory, repository := testRepository(t)
	defer os.RemoveAll(directory)
	testConfig(t, repository)
	worktree := testWorktree(t, repository)
	testAddAndCommitNewFile(t, worktree, "some-file")
	head, err := repository.Head()
	if err != nil {
		t.Fatal(err)
	}
	tag := "some-name"

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_tag" "test" {
						directory = "%s"
						name      = "%s"
						sha1      = "%s"
					}
				`, directory, tag, head.Hash().String()),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("git_tag.test", "directory", directory),
					resource.TestCheckResourceAttr("git_tag.test", "id", tag),
					resource.TestCheckResourceAttr("git_tag.test", "name", tag),
				),
			},
		},
	})
}

func TestResourceGitTag_InvalidRepository(t *testing.T) {
	t.Parallel()
	tag := "some-name"

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_tag" "test" {
						directory = "/some/random/path"
						name      = "%s"
					}
				`, tag),
				ExpectError: regexp.MustCompile(`Cannot open repository`),
			},
		},
	})
}

func TestResourceGitTag_MissingRepository(t *testing.T) {
	t.Parallel()
	tag := "some-name"

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_tag" "test" {
						name      = "%s"
					}
				`, tag),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
		},
	})
}

func TestResourceGitTag_MissingName(t *testing.T) {
	t.Parallel()
	directory, _ := testRepository(t)
	defer os.RemoveAll(directory)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_tag" "test" {
						directory = "%s"
					}
				`, directory),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
		},
	})
}

func TestResourceGitTag_Import(t *testing.T) {
	t.Parallel()
	directory, repository := testRepository(t)
	defer os.RemoveAll(directory)
	testConfig(t, repository)
	worktree := testWorktree(t, repository)
	testAddAndCommitNewFile(t, worktree, "some-file")
	tag := "some-name"

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_tag" "test" {
						directory = "%s"
						name      = "%s"
					}
				`, directory, tag),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("git_tag.test", "directory", directory),
					resource.TestCheckResourceAttr("git_tag.test", "id", tag),
					resource.TestCheckResourceAttr("git_tag.test", "name", tag),
				),
			},
			{
				ResourceName:      "git_tag.test",
				ImportState:       true,
				ImportStateId:     fmt.Sprintf("%s|%s", directory, tag),
				ImportStateVerify: true,
			},
		},
	})
}

func TestResourceGitTag_Update_Name(t *testing.T) {
	t.Parallel()
	directory, repository := testRepository(t)
	defer os.RemoveAll(directory)
	testConfig(t, repository)
	worktree := testWorktree(t, repository)
	testAddAndCommitNewFile(t, worktree, "some-file")
	tag := "some-name"
	newTag := "other-name"

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_tag" "test" {
						directory = "%s"
						name      = "%s"
					}
				`, directory, tag),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("git_tag.test", "directory", directory),
					resource.TestCheckResourceAttr("git_tag.test", "id", tag),
					resource.TestCheckResourceAttr("git_tag.test", "name", tag),
				),
			},
			{
				Config: fmt.Sprintf(`
					resource "git_tag" "test" {
						directory = "%s"
						name      = "%s"
					}
				`, directory, newTag),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("git_tag.test", "directory", directory),
					resource.TestCheckResourceAttr("git_tag.test", "id", newTag),
					resource.TestCheckResourceAttr("git_tag.test", "name", newTag),
				),
			},
		},
	})
}

func TestResourceGitTag_Update_Directory(t *testing.T) {
	t.Parallel()
	directory, repository := testRepository(t)
	defer os.RemoveAll(directory)
	testConfig(t, repository)
	worktree := testWorktree(t, repository)
	testAddAndCommitNewFile(t, worktree, "some-file")
	newDirectory, newRepository := testRepository(t)
	defer os.RemoveAll(newDirectory)
	testConfig(t, newRepository)
	newWorktree := testWorktree(t, newRepository)
	testAddAndCommitNewFile(t, newWorktree, "other-file")
	tag := "some-name"

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_tag" "test" {
						directory = "%s"
						name      = "%s"
					}
				`, directory, tag),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("git_tag.test", "directory", directory),
					resource.TestCheckResourceAttr("git_tag.test", "id", tag),
					resource.TestCheckResourceAttr("git_tag.test", "name", tag),
				),
			},
			{
				Config: fmt.Sprintf(`
					resource "git_tag" "test" {
						directory = "%s"
						name      = "%s"
					}
				`, newDirectory, tag),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("git_tag.test", "directory", newDirectory),
					resource.TestCheckResourceAttr("git_tag.test", "id", tag),
					resource.TestCheckResourceAttr("git_tag.test", "name", tag),
				),
			},
		},
	})
}
