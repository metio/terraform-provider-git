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

func TestResourceGitAdd(t *testing.T) {
	t.Parallel()
	directory, repository := testRepository(t)
	defer os.RemoveAll(directory)
	worktree := testWorktree(t, repository)
	name := "some-file"
	testWriteFile(t, worktree, name)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_add" "test" {
						directory = "%s"
						file      = "%s"
					}
				`, directory, name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("git_add.test", "directory", directory),
					resource.TestCheckResourceAttr("git_add.test", "id", fmt.Sprintf("%s|%s", directory, name)),
					resource.TestCheckResourceAttr("git_add.test", "file", name),
				),
			},
		},
	})
}

func TestResourceGitAdd_NonExistingFile(t *testing.T) {
	t.Parallel()
	directory, _ := testRepository(t)
	defer os.RemoveAll(directory)
	name := "some-file"

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_add" "test" {
						directory = "%s"
						file      = "%s"
					}
				`, directory, name),
				ExpectError: regexp.MustCompile(`Cannot get file infos`),
			},
		},
	})
}

func TestResourceGitAdd_AddDirectory(t *testing.T) {
	t.Parallel()
	directory, _ := testRepository(t)
	defer os.RemoveAll(directory)
	err := os.Mkdir(directory+"/nested-folder", 0700)
	if err != nil {
		t.Fatal(err)
	}

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_add" "test" {
						directory = "%s"
						file      = "%s"
					}
				`, directory, "nested-folder"),
				ExpectError: regexp.MustCompile(`Cannot open directory`),
			},
		},
	})
}

func TestResourceGitAdd_BareRepository(t *testing.T) {
	t.Parallel()
	directory := testRepositoryBare(t)
	defer os.RemoveAll(directory)
	name := "some-file"

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_add" "test" {
						directory = "%s"
						file      = "%s"
					}
				`, directory, name),
				ExpectError: regexp.MustCompile(`Cannot add file to bare repository`),
			},
		},
	})
}

func TestResourceGitAdd_InvalidRepository(t *testing.T) {
	t.Parallel()
	name := "some-file"

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_add" "test" {
						directory = "/some/random/path"
						file      = "%s"
					}
				`, name),
				ExpectError: regexp.MustCompile(`Cannot open repository`),
			},
		},
	})
}

func TestResourceGitAdd_MissingDirectory(t *testing.T) {
	t.Parallel()
	name := "some-file"

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_add" "test" {
						file      = "%s"
					}
				`, name),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
		},
	})
}

func TestResourceGitAdd_MissingFile(t *testing.T) {
	t.Parallel()
	directory, _ := testRepository(t)
	defer os.RemoveAll(directory)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_add" "test" {
						directory = "%s"
					}
				`, directory),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
		},
	})
}

func TestResourceGitAdd_Import(t *testing.T) {
	t.Parallel()
	directory, repository := testRepository(t)
	defer os.RemoveAll(directory)
	worktree := testWorktree(t, repository)
	name := "some-file"
	testWriteFile(t, worktree, name)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_add" "test" {
						directory = "%s"
						file      = "%s"
					}
				`, directory, name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("git_add.test", "directory", directory),
					resource.TestCheckResourceAttr("git_add.test", "id", fmt.Sprintf("%s|%s", directory, name)),
					resource.TestCheckResourceAttr("git_add.test", "file", name),
				),
			},
			{
				ResourceName:      "git_add.test",
				ImportState:       true,
				ImportStateId:     fmt.Sprintf("%s|%s", directory, name),
				ImportStateVerify: true,
			},
		},
	})
}

func TestResourceGitAdd_Update_File(t *testing.T) {
	t.Parallel()
	directory, repository := testRepository(t)
	defer os.RemoveAll(directory)
	worktree := testWorktree(t, repository)
	name1 := "some-file"
	name2 := "other-file"
	testWriteFile(t, worktree, name1)
	testWriteFile(t, worktree, name2)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_add" "test" {
						directory = "%s"
						file      = "%s"
					}
				`, directory, name1),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("git_add.test", "directory", directory),
					resource.TestCheckResourceAttr("git_add.test", "id", fmt.Sprintf("%s|%s", directory, name1)),
					resource.TestCheckResourceAttr("git_add.test", "file", name1),
				),
			},
			{
				Config: fmt.Sprintf(`
					resource "git_add" "test" {
						directory = "%s"
						file      = "%s"
					}
				`, directory, name2),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("git_add.test", "directory", directory),
					resource.TestCheckResourceAttr("git_add.test", "id", fmt.Sprintf("%s|%s", directory, name2)),
					resource.TestCheckResourceAttr("git_add.test", "file", name2),
				),
			},
		},
	})
}
