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

func TestDataSourceGitRepository(t *testing.T) {
	t.Parallel()
	directory, repository := testRepository(t)
	defer os.RemoveAll(directory)
	worktree := testWorktree(t, repository)
	testAddAndCommitNewFile(t, worktree, "some-file")

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "git_repository" "test" {
						directory = "%s"
					}
				`, directory),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.git_repository.test", "directory", directory),
					resource.TestCheckResourceAttr("data.git_repository.test", "id", directory),
					resource.TestCheckResourceAttr("data.git_repository.test", "branch", "master"),
					resource.TestCheckResourceAttrWith("data.git_repository.test", "sha1", testCheckMinLength(4)),
				),
			},
		},
	})
}

func TestDataSourceGitRepository_Detached(t *testing.T) {
	t.Parallel()
	directory, repository := testRepository(t)
	defer os.RemoveAll(directory)
	worktree := testWorktree(t, repository)
	testAddAndCommitNewFile(t, worktree, "some-file")
	head := testReadHead(t, repository)
	testGitCheckout(t, worktree, head.Hash())

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "git_repository" "test" {
						directory = "%s"
					}
				`, directory),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.git_repository.test", "directory", directory),
					resource.TestCheckResourceAttr("data.git_repository.test", "id", directory),
					resource.TestCheckNoResourceAttr("data.git_repository.test", "branch"),
					resource.TestCheckResourceAttrWith("data.git_repository.test", "sha1", testCheckMinLength(4)),
				),
			},
		},
	})
}

func TestDataSourceGitRepository_Empty(t *testing.T) {
	t.Parallel()
	directory, _ := testRepository(t)
	defer os.RemoveAll(directory)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "git_repository" "test" {
						directory = "%s"
					}
				`, directory),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.git_repository.test", "directory", directory),
					resource.TestCheckResourceAttr("data.git_repository.test", "id", directory),
					resource.TestCheckNoResourceAttr("data.git_repository.test", "branch"),
					resource.TestCheckNoResourceAttr("data.git_repository.test", "sha1"),
				),
			},
		},
	})
}

func TestDataSourceGitRepository_InvalidRepository(t *testing.T) {
	t.Parallel()
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: `
					data "git_repository" "test" {
						directory = "/some/random/path"
					}
				`,
				ExpectError: regexp.MustCompile(`Cannot open repository`),
			},
		},
	})
}

func TestDataSourceGitRepository_MissingRepository(t *testing.T) {
	t.Parallel()
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: `
					data "git_repository" "test" {}
				`,
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
		},
	})
}
