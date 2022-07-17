/*
 * This file is part of terraform-provider-git. It is subject to the license terms in the LICENSE file found in the top-level
 * directory of this distribution and at https://creativecommons.org/publicdomain/zero/1.0/. No part of terraform-provider-git,
 * including this file, may be copied, modified, propagated, or distributed except according to the terms contained
 * in the LICENSE file.
 */

package provider

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
)

func TestDataSourceGitBranch(t *testing.T) {
	directory, err := ioutil.TempDir("", "data_source_git_branch")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(directory)
	repository, err := git.PlainInit(directory, false)
	if err != nil {
		t.Fatal(err)
	}
	branch := "name-of-branch"
	err = repository.CreateBranch(&config.Branch{
		Name:   branch,
		Remote: "origin",
		Rebase: "true",
	})
	if err != nil {
		t.Fatal(err)
	}

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: protoV6ProviderFactories(),

		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "git_branch" "branch" {
						directory = "%s"
						branch    = "%s"
					}
				`, directory, branch),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.git_branch.branch", "directory", directory),
					resource.TestCheckResourceAttr("data.git_branch.branch", "id", directory),
					resource.TestCheckResourceAttr("data.git_branch.branch", "branch", branch),
					resource.TestCheckResourceAttr("data.git_branch.branch", "remote", "origin"),
					resource.TestCheckResourceAttr("data.git_branch.branch", "rebase", "true"),
					//resource.TestCheckResourceAttrWith("data.git_branch.branch", "sha1", testCheckMinLen(6)),
				),
			},
		},
	})
}

func testCheckMinLen(minLen int) func(input string) error {
	return func(input string) error {
		if len(input) < minLen {
			return fmt.Errorf("minimum length %d, actual length %d", minLen, len(input))
		}

		return nil
	}
}
