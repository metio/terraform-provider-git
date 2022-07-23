/*
 * This file is part of terraform-provider-git. It is subject to the license terms in the LICENSE file found in the top-level
 * directory of this distribution and at https://creativecommons.org/publicdomain/zero/1.0/. No part of terraform-provider-git,
 * including this file, may be copied, modified, propagated, or distributed except according to the terms contained
 * in the LICENSE file.
 */

package provider_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"os"
	"regexp"
	"testing"
)

func TestDataSourceGitBranches(t *testing.T) {
	t.Parallel()
	directory, repository := initializeGitRepository(t)
	defer os.RemoveAll(directory)
	worktree := createWorktree(t, repository)
	addAndCommitNewFile(t, worktree)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: protoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "git_branches" "test" {
						directory = "%s"
					}
				`, directory),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.git_branches.test", "directory", directory),
					resource.TestCheckResourceAttr("data.git_branches.test", "id", directory),
					resource.TestCheckResourceAttr("data.git_branches.test", "branches.%", "1"),
					resource.TestCheckResourceAttr("data.git_branches.test", "branches.master.remote", ""),
					resource.TestCheckResourceAttr("data.git_branches.test", "branches.master.rebase", ""),
					resource.TestCheckResourceAttrWith("data.git_branches.test", "branches.master.sha1", testCheckLen(40)),
				),
			},
		},
	})
}

func TestDataSourceGitBranches_InvalidRepository(t *testing.T) {
	t.Parallel()
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: protoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: `
					data "git_branches" "test" {
						directory = "/some/random/path"
					}
				`,
				ExpectError: regexp.MustCompile(`Cannot open repository`),
			},
		},
	})
}

func TestDataSourceGitBranches_MissingRepository(t *testing.T) {
	t.Parallel()
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: protoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: `
					data "git_branches" "test" {}
				`,
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
		},
	})
}
