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

func TestDataSourceGitTag(t *testing.T) {
	directory, repository := initializeGitRepository(t)
	defer os.RemoveAll(directory)
	worktree := createWorktree(t, repository)
	addAndCommitNewFile(t, worktree)
	tag := "some-tag"
	createTag(t, repository, tag)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: protoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "git_tag" "test" {
						directory = "%s"
						tag       = "%s"
					}
				`, directory, tag),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.git_tag.test", "directory", directory),
					resource.TestCheckResourceAttr("data.git_tag.test", "id", directory),
					resource.TestCheckResourceAttr("data.git_tag.test", "tag", tag),
					resource.TestCheckResourceAttr("data.git_tag.test", "annotated", "true"),
					resource.TestCheckResourceAttr("data.git_tag.test", "lightweight", "false"),
					resource.TestCheckResourceAttrWith("data.git_tag.test", "sha1", testCheckLen(40)),
				),
			},
		},
	})
}
