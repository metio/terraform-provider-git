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

func TestDataSourceGitTags(t *testing.T) {
	t.Parallel()
	directory, repository := testRepository(t)
	defer os.RemoveAll(directory)
	worktree := testWorktree(t, repository)
	testAddAndCommitNewFile(t, worktree, "some-file")
	tag := "some-tag"
	testCreateTag(t, repository, tag)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "git_tags" "test" {
						directory = "%s"
					}
				`, directory),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.git_tags.test", "directory", directory),
					resource.TestCheckResourceAttr("data.git_tags.test", "id", directory),
					resource.TestCheckResourceAttr("data.git_tags.test", "annotated", "true"),
					resource.TestCheckResourceAttr("data.git_tags.test", "lightweight", "true"),
					resource.TestCheckResourceAttr("data.git_tags.test", "tags.%", "1"),
				),
			},
		},
	})
}

func TestDataSourceGitTags_NoAnnotated(t *testing.T) {
	t.Parallel()
	directory, repository := testRepository(t)
	defer os.RemoveAll(directory)
	worktree := testWorktree(t, repository)
	testAddAndCommitNewFile(t, worktree, "some-file")
	tag := "some-tag"
	testCreateTag(t, repository, tag)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "git_tags" "test" {
						directory = "%s"
						annotated = false
					}
				`, directory),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.git_tags.test", "directory", directory),
					resource.TestCheckResourceAttr("data.git_tags.test", "id", directory),
					resource.TestCheckResourceAttr("data.git_tags.test", "annotated", "false"),
					resource.TestCheckResourceAttr("data.git_tags.test", "lightweight", "true"),
					resource.TestCheckResourceAttr("data.git_tags.test", "tags.%", "0"),
				),
			},
		},
	})
}

func TestDataSourceGitTags_NoLightweight(t *testing.T) {
	t.Parallel()
	directory, repository := testRepository(t)
	defer os.RemoveAll(directory)
	worktree := testWorktree(t, repository)
	testAddAndCommitNewFile(t, worktree, "some-file")
	tag := "some-tag"
	testCreateTag(t, repository, tag)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "git_tags" "test" {
						directory   = "%s"
						lightweight = false
					}
				`, directory),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.git_tags.test", "directory", directory),
					resource.TestCheckResourceAttr("data.git_tags.test", "id", directory),
					resource.TestCheckResourceAttr("data.git_tags.test", "annotated", "true"),
					resource.TestCheckResourceAttr("data.git_tags.test", "lightweight", "false"),
					resource.TestCheckResourceAttr("data.git_tags.test", "tags.%", "1"),
				),
			},
		},
	})
}

func TestDataSourceGitTags_InvalidRepository(t *testing.T) {
	t.Parallel()
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: `
					data "git_tags" "test" {
						directory = "/some/random/path"
					}
				`,
				ExpectError: regexp.MustCompile(`Cannot open repository`),
			},
		},
	})
}

func TestDataSourceGitTags_MissingRepository(t *testing.T) {
	t.Parallel()
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: `
					data "git_tags" "test" {}
				`,
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
		},
	})
}
