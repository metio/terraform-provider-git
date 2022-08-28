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

func TestResourceGitTag(t *testing.T) {
	t.Parallel()
	directory, repository := testRepository(t)
	defer os.RemoveAll(directory)
	testConfig(t, repository)
	worktree := testWorktree(t, repository)
	testAddAndCommitNewFile(t, worktree, "some-file")
	name := "some-name"

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_tag" "test" {
						directory = "%s"
						name      = "%s"
					}
				`, directory, name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("git_tag.test", "directory", directory),
					resource.TestCheckResourceAttr("git_tag.test", "id", fmt.Sprintf("%s|%s", directory, name)),
					resource.TestCheckResourceAttr("git_tag.test", "name", name),
					resource.TestCheckResourceAttrWith("git_tag.test", "revision", testCheckMinLength(1)),
					resource.TestCheckResourceAttrWith("git_tag.test", "sha1", testCheckMinLength(1)),
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
	name := "some-name"
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
				`, directory, name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("git_tag.test", "directory", directory),
					resource.TestCheckResourceAttr("git_tag.test", "id", fmt.Sprintf("%s|%s", directory, name)),
					resource.TestCheckResourceAttr("git_tag.test", "name", name),
					resource.TestCheckResourceAttr("git_tag.test", "message", message),
				),
			},
		},
	})
}

func TestResourceGitTag_Revision_Hash(t *testing.T) {
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
	name := "some-name"

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_tag" "test" {
						directory = "%s"
						name      = "%s"
						revision  = "%s"
					}
				`, directory, name, head.Hash().String()),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("git_tag.test", "directory", directory),
					resource.TestCheckResourceAttr("git_tag.test", "id", fmt.Sprintf("%s|%s", directory, name)),
					resource.TestCheckResourceAttr("git_tag.test", "name", name),
					resource.TestCheckResourceAttr("git_tag.test", "revision", head.Hash().String()),
					resource.TestCheckResourceAttr("git_tag.test", "sha1", head.Hash().String()),
				),
			},
		},
	})
}

func TestResourceGitTag_Revision_Head(t *testing.T) {
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
	name := "some-name"

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_tag" "test" {
						directory = "%s"
						name      = "%s"
						revision  = "HEAD"
					}
				`, directory, name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("git_tag.test", "directory", directory),
					resource.TestCheckResourceAttr("git_tag.test", "id", fmt.Sprintf("%s|%s", directory, name)),
					resource.TestCheckResourceAttr("git_tag.test", "name", name),
					resource.TestCheckResourceAttr("git_tag.test", "revision", "HEAD"),
					resource.TestCheckResourceAttr("git_tag.test", "sha1", head.Hash().String()),
				),
			},
		},
	})
}

func TestResourceGitTag_Revision_Master(t *testing.T) {
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
	name := "some-name"

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_tag" "test" {
						directory = "%s"
						name      = "%s"
						revision  = "master"
					}
				`, directory, name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("git_tag.test", "directory", directory),
					resource.TestCheckResourceAttr("git_tag.test", "id", fmt.Sprintf("%s|%s", directory, name)),
					resource.TestCheckResourceAttr("git_tag.test", "name", name),
					resource.TestCheckResourceAttr("git_tag.test", "revision", "master"),
					resource.TestCheckResourceAttr("git_tag.test", "sha1", head.Hash().String()),
				),
			},
		},
	})
}

func TestResourceGitTag_Directory_Invalid(t *testing.T) {
	t.Parallel()
	name := "some-name"

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_tag" "test" {
						directory = "/some/random/path"
						name      = "%s"
					}
				`, name),
				ExpectError: regexp.MustCompile(`Cannot open repository`),
			},
		},
	})
}

func TestResourceGitTag_Directory_Missing(t *testing.T) {
	t.Parallel()
	name := "some-name"

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_tag" "test" {
						name      = "%s"
					}
				`, name),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
		},
	})
}

func TestResourceGitTag_Name_Missing(t *testing.T) {
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
	name := "some-name"

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_tag" "test" {
						directory = "%s"
						name      = "%s"
					}
				`, directory, name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("git_tag.test", "directory", directory),
					resource.TestCheckResourceAttr("git_tag.test", "id", fmt.Sprintf("%s|%s", directory, name)),
					resource.TestCheckResourceAttr("git_tag.test", "name", name),
					resource.TestCheckResourceAttr("git_tag.test", "revision", "HEAD"),
				),
			},
			{
				ResourceName:      "git_tag.test",
				ImportState:       true,
				ImportStateId:     fmt.Sprintf("%s|%s", directory, name),
				ImportStateVerify: true,
			},
		},
	})
}

func TestResourceGitTag_Import_Symbolic(t *testing.T) {
	t.Parallel()
	directory, repository := testRepository(t)
	defer os.RemoveAll(directory)
	testConfig(t, repository)
	worktree := testWorktree(t, repository)
	testAddAndCommitNewFile(t, worktree, "some-file")
	name := "some-name"
	revision := "master"

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_tag" "test" {
						directory = "%s"
						name      = "%s"
						revision  = "%s" 
					}
				`, directory, name, revision),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("git_tag.test", "directory", directory),
					resource.TestCheckResourceAttr("git_tag.test", "id", fmt.Sprintf("%s|%s", directory, name)),
					resource.TestCheckResourceAttr("git_tag.test", "name", name),
					resource.TestCheckResourceAttr("git_tag.test", "revision", revision),
				),
			},
			{
				ResourceName:      "git_tag.test",
				ImportState:       true,
				ImportStateId:     fmt.Sprintf("%s|%s|%s", directory, name, revision),
				ImportStateVerify: true,
			},
		},
	})
}

func TestResourceGitTag_Name_Update(t *testing.T) {
	t.Parallel()
	directory, repository := testRepository(t)
	defer os.RemoveAll(directory)
	testConfig(t, repository)
	worktree := testWorktree(t, repository)
	testAddAndCommitNewFile(t, worktree, "some-file")
	name := "some-name"
	newName := "other-name"

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_tag" "test" {
						directory = "%s"
						name      = "%s"
					}
				`, directory, name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("git_tag.test", "directory", directory),
					resource.TestCheckResourceAttr("git_tag.test", "id", fmt.Sprintf("%s|%s", directory, name)),
					resource.TestCheckResourceAttr("git_tag.test", "name", name),
				),
			},
			{
				Config: fmt.Sprintf(`
					resource "git_tag" "test" {
						directory = "%s"
						name      = "%s"
					}
				`, directory, newName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("git_tag.test", "directory", directory),
					resource.TestCheckResourceAttr("git_tag.test", "id", fmt.Sprintf("%s|%s", directory, newName)),
					resource.TestCheckResourceAttr("git_tag.test", "name", newName),
				),
			},
		},
	})
}

func TestResourceGitTag_Directory_Update(t *testing.T) {
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
	name := "some-name"

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_tag" "test" {
						directory = "%s"
						name      = "%s"
					}
				`, directory, name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("git_tag.test", "directory", directory),
					resource.TestCheckResourceAttr("git_tag.test", "id", fmt.Sprintf("%s|%s", directory, name)),
					resource.TestCheckResourceAttr("git_tag.test", "name", name),
				),
			},
			{
				Config: fmt.Sprintf(`
					resource "git_tag" "test" {
						directory = "%s"
						name      = "%s"
					}
				`, newDirectory, name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("git_tag.test", "directory", newDirectory),
					resource.TestCheckResourceAttr("git_tag.test", "id", fmt.Sprintf("%s|%s", newDirectory, name)),
					resource.TestCheckResourceAttr("git_tag.test", "name", name),
				),
			},
		},
	})
}
