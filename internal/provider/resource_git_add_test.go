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

func TestResourceGitAdd(t *testing.T) {
	t.Parallel()
	directory, repository := testutils.CreateRepository(t)
	defer os.RemoveAll(directory)
	worktree := testutils.GetRepositoryWorktree(t, repository)
	name := "some-file"
	testutils.WriteFileInWorktree(t, worktree, name)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_add" "test" {
						directory  = "%s"
					}
				`, directory),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("git_add.test", "directory", directory),
					resource.TestCheckResourceAttr("git_add.test", "id", directory),
					resource.TestCheckResourceAttr("git_add.test", "all", "true"),
					resource.TestCheckNoResourceAttr("git_add.test", "exact_path"),
					resource.TestCheckNoResourceAttr("git_add.test", "glob_path"),
				),
			},
		},
	})
}

func TestResourceGitAdd_All_Disabled(t *testing.T) {
	t.Parallel()
	directory, repository := testutils.CreateRepository(t)
	defer os.RemoveAll(directory)
	worktree := testutils.GetRepositoryWorktree(t, repository)
	name := "some-file"
	testutils.WriteFileInWorktree(t, worktree, name)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_add" "test" {
						directory  = "%s"
						all        = "false"
					}
				`, directory),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("git_add.test", "directory", directory),
					resource.TestCheckResourceAttr("git_add.test", "id", directory),
					resource.TestCheckResourceAttr("git_add.test", "all", "false"),
					resource.TestCheckNoResourceAttr("git_add.test", "exact_path"),
					resource.TestCheckNoResourceAttr("git_add.test", "glob_path"),
				),
			},
		},
	})
}

func TestResourceGitAdd_ExactPath(t *testing.T) {
	t.Parallel()
	directory, repository := testutils.CreateRepository(t)
	defer os.RemoveAll(directory)
	worktree := testutils.GetRepositoryWorktree(t, repository)
	name := "some-file"
	testutils.WriteFileInWorktree(t, worktree, name)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_add" "test" {
						directory  = "%s"
						exact_path = "%s"
					}
				`, directory, name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("git_add.test", "directory", directory),
					resource.TestCheckResourceAttr("git_add.test", "id", directory),
					resource.TestCheckResourceAttr("git_add.test", "all", "true"),
					resource.TestCheckResourceAttr("git_add.test", "exact_path", name),
					resource.TestCheckNoResourceAttr("git_add.test", "glob_path"),
				),
			},
		},
	})
}

func TestResourceGitAdd_GlobPath(t *testing.T) {
	t.Parallel()
	directory, repository := testutils.CreateRepository(t)
	defer os.RemoveAll(directory)
	worktree := testutils.GetRepositoryWorktree(t, repository)
	name := "some-file"
	testutils.WriteFileInWorktree(t, worktree, name)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_add" "test" {
						directory  = "%s"
						glob_path  = "some*"
					}
				`, directory),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("git_add.test", "directory", directory),
					resource.TestCheckResourceAttr("git_add.test", "id", directory),
					resource.TestCheckResourceAttr("git_add.test", "all", "true"),
					resource.TestCheckNoResourceAttr("git_add.test", "exact_path"),
					resource.TestCheckResourceAttr("git_add.test", "glob_path", "some*"),
				),
			},
		},
	})
}

func TestResourceGitAdd_ExactPath_NonExistingFile(t *testing.T) {
	t.Parallel()
	directory, _ := testutils.CreateRepository(t)
	defer os.RemoveAll(directory)
	name := "some-file"

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_add" "test" {
						directory  = "%s"
						exact_path = "%s"
					}
				`, directory, name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("git_add.test", "directory", directory),
					resource.TestCheckResourceAttr("git_add.test", "id", directory),
					resource.TestCheckResourceAttr("git_add.test", "all", "true"),
					resource.TestCheckResourceAttr("git_add.test", "exact_path", name),
					resource.TestCheckNoResourceAttr("git_add.test", "glob_path"),
				),
			},
		},
	})
}

func TestResourceGitAdd_ExactPath_Directory(t *testing.T) {
	t.Parallel()
	directory, _ := testutils.CreateRepository(t)
	defer os.RemoveAll(directory)
	err := os.Mkdir(directory+"/nested-folder", 0700)
	if err != nil {
		t.Fatal(err)
	}

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_add" "test" {
						directory  = "%s"
						exact_path = "%s"
					}
				`, directory, "nested-folder"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("git_add.test", "directory", directory),
					resource.TestCheckResourceAttr("git_add.test", "id", directory),
					resource.TestCheckResourceAttr("git_add.test", "all", "true"),
					resource.TestCheckResourceAttr("git_add.test", "exact_path", "nested-folder"),
					resource.TestCheckNoResourceAttr("git_add.test", "glob_path"),
				),
			},
		},
	})
}

func TestResourceGitAdd_ExactPath_GlobPath(t *testing.T) {
	t.Parallel()
	directory, _ := testutils.CreateRepository(t)
	defer os.RemoveAll(directory)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_add" "test" {
						directory  = "%s"
						exact_path = "some-file"
						glob_path  = "some*"
					}
				`, directory),
				ExpectError: regexp.MustCompile(`Invalid Attribute Combination`),
			},
		},
	})
}

func TestResourceGitAdd_ExactPath_EmptyString(t *testing.T) {
	t.Parallel()
	directory, _ := testutils.CreateRepository(t)
	defer os.RemoveAll(directory)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_add" "test" {
						directory  = "%s"
						exact_path = ""
					}
				`, directory),
				ExpectError: regexp.MustCompile(`Invalid Attribute Value Length`),
			},
		},
	})
}

func TestResourceGitAdd_GlobPath_EmptyString(t *testing.T) {
	t.Parallel()
	directory, _ := testutils.CreateRepository(t)
	defer os.RemoveAll(directory)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_add" "test" {
						directory  = "%s"
						glob_path  = ""
					}
				`, directory),
				ExpectError: regexp.MustCompile(`Invalid Attribute Value Length`),
			},
		},
	})
}

func TestResourceGitAdd_BareRepository(t *testing.T) {
	t.Parallel()
	directory := testutils.CreateBareRepository(t)
	defer os.RemoveAll(directory)
	name := "some-file"

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_add" "test" {
						directory  = "%s"
						exact_path = "%s"
					}
				`, directory, name),
				ExpectError: regexp.MustCompile(`Cannot add file to bare repository`),
			},
		},
	})
}

func TestResourceGitAdd_Directory_Invalid(t *testing.T) {
	t.Parallel()
	name := "some-file"

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_add" "test" {
						directory  = "/some/random/path"
						exact_path = "%s"
					}
				`, name),
				ExpectError: regexp.MustCompile(`Cannot open repository`),
			},
		},
	})
}

func TestResourceGitAdd_Directory_Missing(t *testing.T) {
	t.Parallel()
	name := "some-file"

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_add" "test" {
						exact_path = "%s"
					}
				`, name),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
		},
	})
}

func TestResourceGitAdd_ExactPath_Update(t *testing.T) {
	t.Parallel()
	directory, repository := testutils.CreateRepository(t)
	defer os.RemoveAll(directory)
	worktree := testutils.GetRepositoryWorktree(t, repository)
	name1 := "some-file"
	name2 := "other-file"
	testutils.WriteFileInWorktree(t, worktree, name1)
	testutils.WriteFileInWorktree(t, worktree, name2)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_add" "test" {
						directory  = "%s"
						exact_path = "%s"
					}
				`, directory, name1),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("git_add.test", "directory", directory),
					resource.TestCheckResourceAttr("git_add.test", "id", directory),
					resource.TestCheckResourceAttr("git_add.test", "all", "true"),
					resource.TestCheckResourceAttr("git_add.test", "exact_path", name1),
					resource.TestCheckNoResourceAttr("git_add.test", "glob_path"),
				),
			},
			{
				Config: fmt.Sprintf(`
					resource "git_add" "test" {
						directory  = "%s"
						exact_path = "%s"
					}
				`, directory, name2),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("git_add.test", "directory", directory),
					resource.TestCheckResourceAttr("git_add.test", "id", directory),
					resource.TestCheckResourceAttr("git_add.test", "all", "true"),
					resource.TestCheckResourceAttr("git_add.test", "exact_path", name2),
					resource.TestCheckNoResourceAttr("git_add.test", "glob_path"),
				),
			},
		},
	})
}

func TestResourceGitAdd_GlobPath_Update(t *testing.T) {
	t.Parallel()
	directory, repository := testutils.CreateRepository(t)
	defer os.RemoveAll(directory)
	worktree := testutils.GetRepositoryWorktree(t, repository)
	name1 := "some-file"
	name2 := "other-file"
	testutils.WriteFileInWorktree(t, worktree, name1)
	testutils.WriteFileInWorktree(t, worktree, name2)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_add" "test" {
						directory  = "%s"
						glob_path  = "some*"
					}
				`, directory),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("git_add.test", "directory", directory),
					resource.TestCheckResourceAttr("git_add.test", "id", directory),
					resource.TestCheckResourceAttr("git_add.test", "all", "true"),
					resource.TestCheckNoResourceAttr("git_add.test", "exact_path"),
					resource.TestCheckResourceAttr("git_add.test", "glob_path", "some*"),
				),
			},
			{
				Config: fmt.Sprintf(`
					resource "git_add" "test" {
						directory  = "%s"
						glob_path  = "other*"
					}
				`, directory),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("git_add.test", "directory", directory),
					resource.TestCheckResourceAttr("git_add.test", "id", directory),
					resource.TestCheckResourceAttr("git_add.test", "all", "true"),
					resource.TestCheckNoResourceAttr("git_add.test", "exact_path"),
					resource.TestCheckResourceAttr("git_add.test", "glob_path", "other*"),
				),
			},
		},
	})
}
