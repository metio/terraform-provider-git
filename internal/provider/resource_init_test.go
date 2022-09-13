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
	"testing"
)

func TestResourceGitInit(t *testing.T) {
	t.Parallel()
	directory := testutils.TemporaryDirectory(t)
	defer os.RemoveAll(directory)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
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
	t.Parallel()
	directory := testutils.TemporaryDirectory(t)
	defer os.RemoveAll(directory)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
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
	t.Parallel()
	directory := testutils.TemporaryDirectory(t)
	defer os.RemoveAll(directory)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
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

func TestResourceGitInit_Import(t *testing.T) {
	t.Parallel()
	directory := testutils.TemporaryDirectory(t)
	defer os.RemoveAll(directory)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
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
			{
				ResourceName:      "git_init.test",
				ImportState:       true,
				ImportStateId:     directory,
				ImportStateVerify: true,
			},
		},
	})
}

func TestResourceGitInit_Delete(t *testing.T) {
	t.Parallel()
	directory := testutils.TemporaryDirectory(t)
	defer os.RemoveAll(directory)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_init" "test" {
						directory = "%s"
					}
				`, directory),
				Destroy: true,
			},
		},
	})
}
