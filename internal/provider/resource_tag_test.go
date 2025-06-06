/*
 * SPDX-FileCopyrightText: The terraform-provider-git Authors
 * SPDX-License-Identifier: 0BSD
 */

package provider_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/metio/terraform-provider-git/internal/testutils"
	"regexp"
	"testing"
)

func TestResourceGitTag(t *testing.T) {
	t.Parallel()
	directory, repository := testutils.CreateRepository(t)
	testutils.TestConfig(t, repository)
	worktree := testutils.GetRepositoryWorktree(t, repository)
	testutils.AddAndCommitNewFile(t, worktree, "some-file")
	name := "some-name"

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
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
					resource.TestCheckResourceAttrWith("git_tag.test", "sha1", testutils.CheckMinLength(1)),
				),
			},
		},
	})
}

func TestResourceGitTag_Annotated(t *testing.T) {
	t.Parallel()
	directory, repository := testutils.CreateRepository(t)
	testutils.TestConfig(t, repository)
	worktree := testutils.GetRepositoryWorktree(t, repository)
	testutils.AddAndCommitNewFile(t, worktree, "some-file")
	name := "some-name"
	message := "some message for the tag"

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
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
					resource.TestCheckResourceAttr("git_tag.test", "revision", "HEAD"),
				),
			},
		},
	})
}

func TestResourceGitTag_Revision_Hash(t *testing.T) {
	t.Parallel()
	directory, repository := testutils.CreateRepository(t)
	testutils.TestConfig(t, repository)
	worktree := testutils.GetRepositoryWorktree(t, repository)
	testutils.AddAndCommitNewFile(t, worktree, "some-file")
	head, err := repository.Head()
	if err != nil {
		t.Fatal(err)
	}
	name := "some-name"

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
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
	directory, repository := testutils.CreateRepository(t)
	testutils.TestConfig(t, repository)
	worktree := testutils.GetRepositoryWorktree(t, repository)
	testutils.AddAndCommitNewFile(t, worktree, "some-file")
	head, err := repository.Head()
	if err != nil {
		t.Fatal(err)
	}
	name := "some-name"

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
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
	directory, repository := testutils.CreateRepository(t)
	testutils.TestConfig(t, repository)
	worktree := testutils.GetRepositoryWorktree(t, repository)
	testutils.AddAndCommitNewFile(t, worktree, "some-file")
	head, err := repository.Head()
	if err != nil {
		t.Fatal(err)
	}
	name := "some-name"

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
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
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
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
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
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
	directory, _ := testutils.CreateRepository(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
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
	directory, repository := testutils.CreateRepository(t)
	testutils.TestConfig(t, repository)
	worktree := testutils.GetRepositoryWorktree(t, repository)
	testutils.AddAndCommitNewFile(t, worktree, "some-file")
	name := "some-name"
	testutils.CreateTag(t, repository, name)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_tag" "test" {
						directory = "%s"
						name      = "%s"
					}
				`, directory, name),
				ResourceName:       "git_tag.test",
				ImportState:        true,
				ImportStateId:      fmt.Sprintf("%s|%s", directory, name),
				ImportStatePersist: true,
				ImportStateCheck: testutils.ComposeImportStateCheck(
					testutils.CheckResourceAttrInstanceState("directory", directory),
					testutils.CheckResourceAttrInstanceState("id", fmt.Sprintf("%s|%s", directory, name)),
					testutils.CheckResourceAttrInstanceState("name", name),
					testutils.CheckResourceAttrInstanceState("revision", "HEAD"),
				),
			},
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
		},
	})
}

func TestResourceGitTag_ImportWithRevision(t *testing.T) {
	t.Parallel()
	directory, repository := testutils.CreateRepository(t)
	testutils.TestConfig(t, repository)
	worktree := testutils.GetRepositoryWorktree(t, repository)
	testutils.AddAndCommitNewFile(t, worktree, "some-file")
	name := "some-name"
	testutils.CreateTag(t, repository, name)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_tag" "test" {
						directory = "%s"
						name      = "%s"
						revision  = "master" 
					}
				`, directory, name),
				ResourceName:       "git_tag.test",
				ImportState:        true,
				ImportStateId:      fmt.Sprintf("%s|%s|master", directory, name),
				ImportStatePersist: true,
				ImportStateCheck: testutils.ComposeImportStateCheck(
					testutils.CheckResourceAttrInstanceState("directory", directory),
					testutils.CheckResourceAttrInstanceState("id", fmt.Sprintf("%s|%s", directory, name)),
					testutils.CheckResourceAttrInstanceState("name", name),
					testutils.CheckResourceAttrInstanceState("revision", "master"),
				),
			},
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
				),
			},
		},
	})
}

func TestResourceGitTag_Import_NonExistingRepo(t *testing.T) {
	t.Parallel()
	directory := testutils.TemporaryDirectory(t)
	name := "some-name"

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_tag" "test" {
						directory = "%s"
						name      = "%s"
					}
				`, directory, name),
				ResourceName:       "git_tag.test",
				ImportState:        true,
				ImportStateId:      fmt.Sprintf("%s|%s", directory, name),
				ImportStatePersist: false,
				ExpectError:        regexp.MustCompile(`Cannot open repository`),
			},
		},
	})
}

func TestResourceGitTag_Name_Update(t *testing.T) {
	t.Parallel()
	directory, repository := testutils.CreateRepository(t)
	testutils.TestConfig(t, repository)
	worktree := testutils.GetRepositoryWorktree(t, repository)
	testutils.AddAndCommitNewFile(t, worktree, "some-file")
	name := "some-name"
	newName := "other-name"

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
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
	directory, repository := testutils.CreateRepository(t)
	testutils.TestConfig(t, repository)
	worktree := testutils.GetRepositoryWorktree(t, repository)
	testutils.AddAndCommitNewFile(t, worktree, "some-file")
	newDirectory, newRepository := testutils.CreateRepository(t)
	testutils.TestConfig(t, newRepository)
	newWorktree := testutils.GetRepositoryWorktree(t, newRepository)
	testutils.AddAndCommitNewFile(t, newWorktree, "other-file")
	name := "some-name"

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
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
