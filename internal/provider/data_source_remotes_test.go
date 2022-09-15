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

func TestDataSourceGitRemotes_GetSchema(t *testing.T) {
	t.Parallel()
	r := &provider.RemotesDataSource{}
	schema, _ := r.GetSchema(context.TODO())

	testutils.VerifySchemaDescriptions(t, schema)
}

func TestDataSourceGitRemotes(t *testing.T) {
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
					data "git_remotes" "test" {
						directory = "%s"
					}
				`, directory),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.git_remotes.test", "directory", directory),
					resource.TestCheckResourceAttr("data.git_remotes.test", "id", directory),
					resource.TestCheckResourceAttr("data.git_remotes.test", "remotes.%", "1"),
					resource.TestCheckResourceAttr("data.git_remotes.test", "remotes.example.urls.#", "1"),
				),
			},
		},
	})
}

func TestDataSourceGitRemotes_InvalidRepository(t *testing.T) {
	t.Parallel()
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: `
					data "git_remotes" "test" {
						directory = "/some/random/path"
					}
				`,
				ExpectError: regexp.MustCompile(`Cannot open repository`),
			},
		},
	})
}

func TestDataSourceGitRemotes_MissingRepository(t *testing.T) {
	t.Parallel()
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: `
					data "git_remotes" "test" {}
				`,
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
		},
	})
}
