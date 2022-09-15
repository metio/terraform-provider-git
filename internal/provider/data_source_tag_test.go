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

func TestDataSourceGitTag_GetSchema(t *testing.T) {
	t.Parallel()
	r := &provider.TagDataSource{}
	schema, _ := r.GetSchema(context.TODO())

	testutils.VerifySchemaDescriptions(t, schema)
}

func TestDataSourceGitTag(t *testing.T) {
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
					data "git_tag" "test" {
						directory = "%s"
						name      = "%s"
					}
				`, directory, tag),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.git_tag.test", "directory", directory),
					resource.TestCheckResourceAttr("data.git_tag.test", "id", tag),
					resource.TestCheckResourceAttr("data.git_tag.test", "name", tag),
					resource.TestCheckResourceAttr("data.git_tag.test", "annotated", "true"),
					resource.TestCheckResourceAttr("data.git_tag.test", "lightweight", "false"),
					resource.TestCheckResourceAttrWith("data.git_tag.test", "sha1", testutils.CheckExactLength(40)),
				),
			},
		},
	})
}

func TestDataSourceGitTag_InvalidRepository(t *testing.T) {
	t.Parallel()
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: `
					data "git_tag" "test" {
						directory = "/some/random/path"
						name      = "does-not-exist"
					}
				`,
				ExpectError: regexp.MustCompile(`Cannot open repository`),
			},
		},
	})
}

func TestDataSourceGitTag_InvalidTag(t *testing.T) {
	t.Parallel()
	directory, _ := testutils.CreateRepository(t)
	defer os.RemoveAll(directory)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "git_tag" "test" {
						directory = "%s"
						name      = "does-not-exist"
					}
				`, directory),
				ExpectError: regexp.MustCompile(`Cannot read tag`),
			},
		},
	})
}

func TestDataSourceGitTag_MissingName(t *testing.T) {
	t.Parallel()
	directory, _ := testutils.CreateRepository(t)
	defer os.RemoveAll(directory)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "git_tag" "test" {
						directory = "%s"
					}
				`, directory),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
		},
	})
}

func TestDataSourceGitTag_MissingDirectory(t *testing.T) {
	t.Parallel()
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: `
					data "git_tag" "test" {
						name = "some-tag"
					}
				`,
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
		},
	})
}
