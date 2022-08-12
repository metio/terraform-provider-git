/*
 * SPDX-FileCopyrightText: Sebastian Hoß <seb@hoß.de>
 * SPDX-License-Identifier: 0BSD
 */

package provider_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"os"
	"testing"
)

func TestResourceGitInit(t *testing.T) {
	t.Parallel()
	directory := testTemporaryDirectory(t)
	defer os.RemoveAll(directory)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testProviderFactories(),
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
	directory := testTemporaryDirectory(t)
	defer os.RemoveAll(directory)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testProviderFactories(),
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
	directory := testTemporaryDirectory(t)
	defer os.RemoveAll(directory)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testProviderFactories(),
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
	directory := testTemporaryDirectory(t)
	defer os.RemoveAll(directory)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testProviderFactories(),
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
				ImportStateId:     fmt.Sprintf("%s", directory),
				ImportStateVerify: true,
			},
		},
	})
}

func TestResourceGitInit_Delete(t *testing.T) {
	t.Parallel()
	directory := testTemporaryDirectory(t)
	defer os.RemoveAll(directory)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testProviderFactories(),
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
