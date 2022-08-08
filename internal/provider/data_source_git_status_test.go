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

func TestDataSourceGitStatus_StagedFile(t *testing.T) {
	t.Parallel()
	directory, repository := initializeGitRepository(t)
	defer os.RemoveAll(directory)
	worktree := createWorktree(t, repository)
	fileName := "some-file"
	addFile(t, worktree, fileName)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: protoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "git_status" "test" {
						directory = "%s"
						file      = "%s"
					}
				`, directory, fileName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.git_status.test", "directory", directory),
					resource.TestCheckResourceAttr("data.git_status.test", "id", fileName),
					resource.TestCheckResourceAttr("data.git_status.test", "file", fileName),
					resource.TestCheckResourceAttr("data.git_status.test", "staging", "A"),
					resource.TestCheckResourceAttr("data.git_status.test", "worktree", " "),
				),
			},
		},
	})
}

func TestDataSourceGitStatus_Clean(t *testing.T) {
	t.Parallel()
	directory, repository := initializeGitRepository(t)
	defer os.RemoveAll(directory)
	worktree := createWorktree(t, repository)
	fileName := "some-file"
	addFile(t, worktree, fileName)
	commitStaged(t, worktree)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: protoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "git_status" "test" {
						directory = "%s"
						file      = "%s"
					}
				`, directory, fileName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.git_status.test", "directory", directory),
					resource.TestCheckResourceAttr("data.git_status.test", "id", fileName),
					resource.TestCheckResourceAttr("data.git_status.test", "file", fileName),
					resource.TestCheckResourceAttr("data.git_status.test", "staging", "?"),
					resource.TestCheckResourceAttr("data.git_status.test", "worktree", "?"),
				),
			},
		},
	})
}

func TestDataSourceGitStatus_BareRepository(t *testing.T) {
	t.Parallel()
	directory := temporaryDirectory(t)
	gitInit(t, directory, true)
	defer os.RemoveAll(directory)
	fileName := "some-file"

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: protoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "git_status" "test" {
						directory = "%s"
						file      = "%s"
					}
				`, directory, fileName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.git_status.test", "directory", directory),
					resource.TestCheckResourceAttr("data.git_status.test", "id", fileName),
					resource.TestCheckResourceAttr("data.git_status.test", "file", fileName),
					resource.TestCheckResourceAttr("data.git_status.test", "staging", ""),
					resource.TestCheckResourceAttr("data.git_status.test", "worktree", ""),
				),
			},
		},
	})
}
