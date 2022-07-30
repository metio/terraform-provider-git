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
	"regexp"
	"testing"
)

func TestDataSourceGitRemote(t *testing.T) {
	t.Parallel()
	directory, repository := initializeGitRepository(t)
	defer os.RemoveAll(directory)
	remote := "example"
	createRemote(t, repository, remote)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: protoV6ProviderFactories(),
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
		ProtoV6ProviderFactories: protoV6ProviderFactories(),
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
	directory, _ := initializeGitRepository(t)
	defer os.RemoveAll(directory)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: protoV6ProviderFactories(),
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
	directory, _ := initializeGitRepository(t)
	defer os.RemoveAll(directory)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: protoV6ProviderFactories(),
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
		ProtoV6ProviderFactories: protoV6ProviderFactories(),
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
