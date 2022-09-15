/*
 * SPDX-FileCopyrightText: The terraform-provider-git Authors
 * SPDX-License-Identifier: 0BSD
 */

package provider_test

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/metio/terraform-provider-git/internal/provider"
	"github.com/metio/terraform-provider-git/internal/testutils"
	"os"
	"regexp"
	"testing"
)

func TestResourceGitAdd_GetSchema(t *testing.T) {
	t.Parallel()
	r := &provider.AddResource{}
	schema, _ := r.GetSchema(context.TODO())

	testutils.VerifySchemaDescriptions(t, schema)
}

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
						add_paths  = ["%s"]
					}
				`, directory, name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("git_add.test", "directory", directory),
					resource.TestCheckResourceAttrWith("git_add.test", "id", testutils.CheckMinLength(1)),
					resource.TestCheckResourceAttr("git_add.test", "add_paths.0", name),
				),
			},
		},
	})
}

func TestResourceGitAdd_AddPaths_Multiple(t *testing.T) {
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
						add_paths  = ["%s", "other-file"]
					}
				`, directory, name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("git_add.test", "directory", directory),
					resource.TestCheckResourceAttrWith("git_add.test", "id", testutils.CheckMinLength(1)),
					resource.TestCheckResourceAttr("git_add.test", "add_paths.0", name),
					resource.TestCheckResourceAttr("git_add.test", "add_paths.1", "other-file"),
				),
			},
		},
	})
}

func TestResourceGitAdd_AddPaths_NonExistingFile(t *testing.T) {
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
						add_paths  = ["%s"]
					}
				`, directory, name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("git_add.test", "directory", directory),
					resource.TestCheckResourceAttrWith("git_add.test", "id", testutils.CheckMinLength(1)),
					resource.TestCheckResourceAttr("git_add.test", "add_paths.0", name),
				),
			},
		},
	})
}

func TestResourceGitAdd_AddPaths_Directory(t *testing.T) {
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
						add_paths  = ["%s"]
					}
				`, directory, "nested-folder"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("git_add.test", "directory", directory),
					resource.TestCheckResourceAttrWith("git_add.test", "id", testutils.CheckMinLength(1)),
					resource.TestCheckResourceAttr("git_add.test", "add_paths.0", "nested-folder"),
				),
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
						add_paths  = ["%s"]
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
						add_paths  = ["%s"]
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
						add_paths = ["%s"]
					}
				`, name),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
		},
	})
}

func TestResourceGitAdd_AddPaths_Update(t *testing.T) {
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
						add_paths  = ["%s"]
					}
				`, directory, name1),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("git_add.test", "directory", directory),
					resource.TestCheckResourceAttrWith("git_add.test", "id", testutils.CheckMinLength(1)),
					resource.TestCheckResourceAttr("git_add.test", "add_paths.0", name1),
				),
			},
			{
				Config: fmt.Sprintf(`
					resource "git_add" "test" {
						directory  = "%s"
						add_paths  = ["%s"]
					}
				`, directory, name2),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("git_add.test", "directory", directory),
					resource.TestCheckResourceAttrWith("git_add.test", "id", testutils.CheckMinLength(1)),
					resource.TestCheckResourceAttr("git_add.test", "add_paths.0", name2),
				),
			},
		},
	})
}