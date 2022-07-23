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

func TestResourceGitInit(t *testing.T) {
	directory := temporaryDirectory(t)
	defer os.RemoveAll(directory)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: protoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_init" "test" {
						directory = "%s"
					}
				`, directory),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("git_init.test", "directory", directory),
					resource.TestCheckResourceAttr("git_init.test", "id", directory),
					resource.TestCheckResourceAttr("git_init.test", "bare", "false"),
				),
			},
		},
	})
}

func TestResourceGitInit_Bare(t *testing.T) {
	directory := temporaryDirectory(t)
	defer os.RemoveAll(directory)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: protoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_init" "test" {
						directory = "%s"
						bare      = true
					}
				`, directory),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("git_init.test", "directory", directory),
					resource.TestCheckResourceAttr("git_init.test", "id", directory),
					resource.TestCheckResourceAttr("git_init.test", "bare", "true"),
				),
			},
		},
	})
}

func TestResourceGitInit_NonBare(t *testing.T) {
	directory := temporaryDirectory(t)
	defer os.RemoveAll(directory)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: protoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_init" "test" {
						directory = "%s"
						bare      = false
					}
				`, directory),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("git_init.test", "directory", directory),
					resource.TestCheckResourceAttr("git_init.test", "id", directory),
					resource.TestCheckResourceAttr("git_init.test", "bare", "false"),
				),
			},
		},
	})
}
