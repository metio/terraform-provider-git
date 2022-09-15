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

func TestResourceGitInit_GetSchema(t *testing.T) {
	t.Parallel()
	r := &provider.InitResource{}
	schema, _ := r.GetSchema(context.TODO())

	testutils.VerifySchemaDescriptions(t, schema)
}

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

func TestResourceGitInit_Import_NonExistingRepo(t *testing.T) {
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
				ResourceName:       "git_init.test",
				ImportState:        true,
				ImportStateId:      directory,
				ImportStatePersist: false,
				ExpectError:        regexp.MustCompile(`Cannot open repository`),
			},
		},
	})
}

func TestResourceGitInit_Import_ExistingRepo(t *testing.T) {
	t.Parallel()
	directory, _ := testutils.CreateRepository(t)
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
				ResourceName:       "git_init.test",
				ImportState:        true,
				ImportStateId:      directory,
				ImportStatePersist: true,
				ImportStateCheck: testutils.ComposeImportStateCheck(
					testutils.CheckResourceAttrInstanceState("directory", directory),
					testutils.CheckResourceAttrInstanceState("id", directory),
					testutils.CheckResourceAttrInstanceState("bare", "false"),
				),
			},
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
