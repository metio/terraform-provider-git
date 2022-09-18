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

func TestResourceGitPush_GetSchema(t *testing.T) {
	t.Parallel()
	r := &provider.PushResource{}
	schema, _ := r.GetSchema(context.TODO())

	testutils.VerifySchemaDescriptions(t, schema)
}

func TestResourceGitPush(t *testing.T) {
	t.Parallel()
	directory, repository := testutils.CreateRepository(t)
	directory2, _ := testutils.CreateRepository(t)
	defer os.RemoveAll(directory)
	testutils.TestConfig(t, repository)
	worktree := testutils.GetRepositoryWorktree(t, repository)
	testutils.AddAndCommitNewFile(t, worktree, "some-file")
	testutils.CreateRemoteWithUrls(t, repository, "origin", []string{directory2})
	head := testutils.GetRepositoryHead(t, repository)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_push" "test" {
						directory = "%s"
						refspecs  = ["%s"]
					}
					data "git_repository" "second" {
						directory = "%s"
						depends_on = [git_push.test]
					}
				`, directory, "refs/heads/master:refs/heads/master", directory2),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("git_push.test", "directory", directory),
					resource.TestCheckResourceAttrWith("git_push.test", "id", testutils.CheckMinLength(1)),
					resource.TestCheckResourceAttr("git_push.test", "refspecs.#", "1"),
					resource.TestCheckResourceAttr("git_push.test", "refspecs.0", "refs/heads/master:refs/heads/master"),
					resource.TestCheckResourceAttr("git_push.test", "prune", "false"),
					resource.TestCheckResourceAttr("git_push.test", "force", "false"),
					resource.TestCheckResourceAttr("data.git_repository.second", "sha1", head.Hash().String()),
				),
			},
		},
	})
}

func TestResourceGitPush_RefSpecs_Invalid(t *testing.T) {
	t.Parallel()
	directory, _ := testutils.CreateRepository(t)
	defer os.RemoveAll(directory)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_push" "test" {
						directory = "%s"
						refspecs  = ["%s"]
					}
				`, directory, "master"),
				ExpectError: regexp.MustCompile(`malformed refspec, separators are wrong`),
			},
		},
	})
}

func TestResourceGitPush_Remote_Unknown(t *testing.T) {
	t.Parallel()
	directory, _ := testutils.CreateRepository(t)
	defer os.RemoveAll(directory)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_push" "test" {
						directory = "%s"
						refspecs  = ["%s"]
					}
				`, directory, "master:master"),
				ExpectError: regexp.MustCompile(`remote not found`),
			},
		},
	})
}

func TestResourceGitPush_Directory_Invalid(t *testing.T) {
	t.Parallel()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_push" "test" {
						directory = "does/not/exist"
						refspecs  = ["%s"]
					}
				`, "master:master"),
				ExpectError: regexp.MustCompile(`Cannot open repository`),
			},
		},
	})
}
