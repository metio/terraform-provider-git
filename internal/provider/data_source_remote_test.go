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

func TestDataSourceGitRemote_GetSchema(t *testing.T) {
	t.Parallel()
	r := &provider.RemoteDataSource{}
	schema, _ := r.GetSchema(context.TODO())

	testutils.VerifySchemaDescriptions(t, schema)
}

func TestDataSourceGitRemote(t *testing.T) {
	t.Parallel()
	directory, repository := testutils.CreateRepository(t)
	defer os.RemoveAll(directory)
	remote := "example"
	testutils.CreateRemote(t, repository, remote)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "git_remote" "test" {
						directory = "%s"
						name      = "%s"
					}
				`, directory, remote),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.git_remote.test", "directory", directory),
					resource.TestCheckResourceAttr("data.git_remote.test", "id", remote),
					resource.TestCheckResourceAttr("data.git_remote.test", "name", remote),
					resource.TestCheckResourceAttr("data.git_remote.test", "urls.#", "1"),
				),
			},
		},
	})
}

func TestDataSourceGitRemote_InvalidRepository(t *testing.T) {
	t.Parallel()
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: `
					data "git_remote" "test" {
						directory = "/some/random/path"
						name      = "does-not-exist"
					}
				`,
				ExpectError: regexp.MustCompile(`Cannot open repository`),
			},
		},
	})
}

func TestDataSourceGitRemote_InvalidRemote(t *testing.T) {
	t.Parallel()
	directory, _ := testutils.CreateRepository(t)
	defer os.RemoveAll(directory)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "git_remote" "test" {
						directory = "%s"
						name      = "does-not-exist"
					}
				`, directory),
				ExpectError: regexp.MustCompile(`Cannot read remote`),
			},
		},
	})
}

func TestDataSourceGitRemote_MissingRemote(t *testing.T) {
	t.Parallel()
	directory, _ := testutils.CreateRepository(t)
	defer os.RemoveAll(directory)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "git_remote" "test" {
						directory = "%s"
					}
				`, directory),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
		},
	})
}

func TestDataSourceGitRemote_MissingRepository(t *testing.T) {
	t.Parallel()
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: `
					data "git_remote" "test" {}
				`,
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
		},
	})
}
