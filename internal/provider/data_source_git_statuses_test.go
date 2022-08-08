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
	"testing"
)

func TestDataSourceGitStatuses_Unclean(t *testing.T) {
	t.Parallel()
	directory, repository := initializeGitRepository(t)
	defer os.RemoveAll(directory)
	worktree := createWorktree(t, repository)
	addFile(t, worktree, "some-file")

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: protoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "git_statuses" "test" {
						directory = "%s"
					}
				`, directory),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.git_statuses.test", "directory", directory),
					resource.TestCheckResourceAttr("data.git_statuses.test", "id", directory),
					resource.TestCheckResourceAttr("data.git_statuses.test", "is_clean", "false"),
					resource.TestCheckResourceAttr("data.git_statuses.test", "files.%", "1"),
				),
			},
		},
	})
}

func TestDataSourceGitStatuses_Clean(t *testing.T) {
	t.Parallel()
	directory, repository := initializeGitRepository(t)
	defer os.RemoveAll(directory)
	worktree := createWorktree(t, repository)
	addFile(t, worktree, "some-file")
	commitStaged(t, worktree)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: protoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "git_statuses" "test" {
						directory = "%s"
					}
				`, directory),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.git_statuses.test", "directory", directory),
					resource.TestCheckResourceAttr("data.git_statuses.test", "id", directory),
					resource.TestCheckResourceAttr("data.git_statuses.test", "is_clean", "true"),
					resource.TestCheckResourceAttr("data.git_statuses.test", "files.%", "0"),
				),
			},
		},
	})
}

func TestDataSourceGitStatuses_BareRepository(t *testing.T) {
	t.Parallel()
	directory := temporaryDirectory(t)
	gitInit(t, directory, true)
	defer os.RemoveAll(directory)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: protoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "git_statuses" "test" {
						directory = "%s"
					}
				`, directory),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.git_statuses.test", "directory", directory),
					resource.TestCheckResourceAttr("data.git_statuses.test", "id", directory),
					resource.TestCheckResourceAttr("data.git_statuses.test", "is_clean", "true"),
					resource.TestCheckResourceAttr("data.git_statuses.test", "files.%", "0"),
				),
			},
		},
	})
}
