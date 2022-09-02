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

func TestDataSourceGitTags(t *testing.T) {
	t.Parallel()
	directory, repository := testutils.CreateRepository(t)
	defer os.RemoveAll(directory)
	worktree := testutils.GetRepositoryWorktree(t, repository)
	testutils.AddAndCommitNewFile(t, worktree, "some-file")
	tag := "some-tag"
	testutils.CreateTag(t, repository, tag)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
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
	directory, repository := testutils.CreateRepository(t)
	defer os.RemoveAll(directory)
	worktree := testutils.GetRepositoryWorktree(t, repository)
	testutils.AddAndCommitNewFile(t, worktree, "some-file")
	tag := "some-tag"
	testutils.CreateTag(t, repository, tag)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
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
	directory, repository := testutils.CreateRepository(t)
	defer os.RemoveAll(directory)
	worktree := testutils.GetRepositoryWorktree(t, repository)
	testutils.AddAndCommitNewFile(t, worktree, "some-file")
	tag := "some-tag"
	testutils.CreateTag(t, repository, tag)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
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
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
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
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
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
