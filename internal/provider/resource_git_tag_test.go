/*
 * This file is part of terraform-provider-git. It is subject to the license terms in the LICENSE file found in the top-level
 * directory of this distribution and at https://creativecommons.org/publicdomain/zero/1.0/. No part of terraform-provider-git,
 * including this file, may be copied, modified, propagated, or distributed except according to the terms contained
 * in the LICENSE file.
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
	directory, repository := initializeGitRepository(t)
	defer os.RemoveAll(directory)
	initTestConfig(t, repository)
	worktree := createWorktree(t, repository)
	addAndCommitNewFile(t, worktree)
	tag := "some-name"

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: protoV6ProviderFactories(),
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
	directory, repository := initializeGitRepository(t)
	defer os.RemoveAll(directory)
	initTestConfig(t, repository)
	worktree := createWorktree(t, repository)
	addAndCommitNewFile(t, worktree)
	tag := "some-name"

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: protoV6ProviderFactories(),
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
				),
			},
		},
	})
}

func TestResourceGitTag_Commitish(t *testing.T) {
	t.Parallel()
	directory, repository := initializeGitRepository(t)
	defer os.RemoveAll(directory)
	initTestConfig(t, repository)
	worktree := createWorktree(t, repository)
	addAndCommitNewFile(t, worktree)
	head, err := repository.Head()
	if err != nil {
		t.Fatal(err)
	}
	tag := "some-name"

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: protoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_tag" "test" {
						directory   = "%s"
						name        = "%s"
						commit_sha1 = "%s"
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
		ProtoV6ProviderFactories: protoV6ProviderFactories(),
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
		ProtoV6ProviderFactories: protoV6ProviderFactories(),
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
	directory, _ := initializeGitRepository(t)
	defer os.RemoveAll(directory)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: protoV6ProviderFactories(),
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
	directory := temporaryDirectory(t)
	defer os.RemoveAll(directory)
	tag := "some-name"

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: protoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_tag" "test" {
						directory = "%s"
						name      = "%s"
					}
				`, directory, tag),
				ImportState:   true,
				ResourceName:  "git_tag.test",
				ImportStateId: tag,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("git_tag.test", "directory", directory),
					resource.TestCheckResourceAttr("git_tag.test", "id", tag),
				),
			},
		},
	})
}

func TestResourceGitTag_Update_Name(t *testing.T) {
	t.Parallel()
	directory, repository := initializeGitRepository(t)
	defer os.RemoveAll(directory)
	initTestConfig(t, repository)
	worktree := createWorktree(t, repository)
	addAndCommitNewFile(t, worktree)
	tag := "some-name"
	newTag := "other-name"

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: protoV6ProviderFactories(),
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
	directory, repository := initializeGitRepository(t)
	defer os.RemoveAll(directory)
	initTestConfig(t, repository)
	worktree := createWorktree(t, repository)
	addAndCommitNewFile(t, worktree)
	newDirectory, newRepository := initializeGitRepository(t)
	defer os.RemoveAll(newDirectory)
	initTestConfig(t, newRepository)
	newWorktree := createWorktree(t, newRepository)
	addAndCommitNewFile(t, newWorktree)
	tag := "some-name"

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: protoV6ProviderFactories(),
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
