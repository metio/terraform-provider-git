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

func TestResourceGitRemote(t *testing.T) {
	t.Parallel()
	directory, _ := initializeGitRepository(t)
	defer os.RemoveAll(directory)
	name := "some-name"
	url1 := "https://github.com/some-org/some-repo.git"

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: protoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_remote" "test" {
						directory = "%s"
						name      = "%s"
						urls      = ["%s"]
					}
				`, directory, name, url1),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("git_remote.test", "directory", directory),
					resource.TestCheckResourceAttr("git_remote.test", "id", name),
					resource.TestCheckResourceAttr("git_remote.test", "name", name),
					resource.TestCheckResourceAttr("git_remote.test", "urls.#", "1"),
					resource.TestCheckResourceAttr("git_remote.test", "urls.0", url1),
				),
			},
		},
	})
}

func TestResourceGitRemote_MultipleUrls(t *testing.T) {
	t.Parallel()
	directory, _ := initializeGitRepository(t)
	defer os.RemoveAll(directory)
	name := "some-name"
	url1 := "https://github.com/some-org/some-repo.git"
	url2 := "https://codeberg.org/some-org/some-repo.git"

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: protoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_remote" "test" {
						directory = "%s"
						name      = "%s"
						urls      = ["%s", "%s"]
					}
				`, directory, name, url1, url2),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("git_remote.test", "directory", directory),
					resource.TestCheckResourceAttr("git_remote.test", "id", name),
					resource.TestCheckResourceAttr("git_remote.test", "name", name),
					resource.TestCheckResourceAttr("git_remote.test", "urls.#", "2"),
					resource.TestCheckResourceAttr("git_remote.test", "urls.0", url1),
					resource.TestCheckResourceAttr("git_remote.test", "urls.1", url2),
				),
			},
		},
	})
}

func TestResourceGitRemote_InvalidRepository(t *testing.T) {
	t.Parallel()
	name := "some-name"

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: protoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_remote" "test" {
						directory = "/some/random/path"
						name      = "%s"
						urls      = ["https://github.com/some-org/some-repo.git", "https://gitlab.com/some-org/some-repo.git"]
					}
				`, name),
				ExpectError: regexp.MustCompile(`Cannot open repository`),
			},
		},
	})
}

func TestResourceGitRemote_MissingRepository(t *testing.T) {
	t.Parallel()
	name := "some-name"

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: protoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_remote" "test" {
						name      = "%s"
						urls      = ["https://github.com/some-org/some-repo.git", "https://gitlab.com/some-org/some-repo.git"]
					}
				`, name),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
		},
	})
}

func TestResourceGitRemote_MissingName(t *testing.T) {
	t.Parallel()
	directory, _ := initializeGitRepository(t)
	defer os.RemoveAll(directory)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: protoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_remote" "test" {
						directory = "%s"
						urls      = ["https://github.com/some-org/some-repo.git", "https://gitlab.com/some-org/some-repo.git"]
					}
				`, directory),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
		},
	})
}

func TestResourceGitRemote_MissingUrls(t *testing.T) {
	t.Parallel()
	directory, _ := initializeGitRepository(t)
	defer os.RemoveAll(directory)
	name := "some-name"

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: protoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_remote" "test" {
						directory = "%s"
						name      = "%s"
					}
				`, directory, name),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
		},
	})
}

func TestResourceGitRemote_Import(t *testing.T) {
	t.Parallel()
	directory, _ := initializeGitRepository(t)
	defer os.RemoveAll(directory)
	remote := "some-name"
	url := "https://github.com/some-org/some-repo.git"

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: protoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_remote" "test" {
						directory = "%s"
						name      = "%s"
						urls      = ["%s"]
					}
				`, directory, remote, url),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("git_remote.test", "directory", directory),
					resource.TestCheckResourceAttr("git_remote.test", "id", remote),
					resource.TestCheckResourceAttr("git_remote.test", "name", remote),
					resource.TestCheckResourceAttr("git_remote.test", "urls.#", "1"),
					resource.TestCheckResourceAttr("git_remote.test", "urls.0", url),
				),
			},
			{
				ResourceName:      "git_remote.test",
				ImportState:       true,
				ImportStateId:     fmt.Sprintf("%s|%s|%s", directory, remote, url),
				ImportStateVerify: true,
			},
		},
	})
}

func TestResourceGitRemote_ImportMultipleUrls(t *testing.T) {
	t.Parallel()
	directory, _ := initializeGitRepository(t)
	defer os.RemoveAll(directory)
	remote := "some-name"
	url1 := "https://github.com/some-org/some-repo.git"
	url2 := "https://codeberg.org/some-org/some-repo.git"

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: protoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_remote" "test" {
						directory = "%s"
						name      = "%s"
						urls      = ["%s", "%s"]
					}
				`, directory, remote, url1, url2),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("git_remote.test", "directory", directory),
					resource.TestCheckResourceAttr("git_remote.test", "id", remote),
					resource.TestCheckResourceAttr("git_remote.test", "name", remote),
					resource.TestCheckResourceAttr("git_remote.test", "urls.#", "2"),
					resource.TestCheckResourceAttr("git_remote.test", "urls.0", url1),
					resource.TestCheckResourceAttr("git_remote.test", "urls.1", url2),
				),
			},
			{
				ResourceName:      "git_remote.test",
				ImportState:       true,
				ImportStateId:     fmt.Sprintf("%s|%s|%s,%s", directory, remote, url1, url2),
				ImportStateVerify: true,
			},
		},
	})
}

func TestResourceGitRemote_Update_Urls(t *testing.T) {
	t.Parallel()
	directory, _ := initializeGitRepository(t)
	defer os.RemoveAll(directory)
	remote := "some-name"
	url1 := "https://github.com/some-org/some-repo.git"
	url2 := "https://codeberg.org/some-org/some-repo.git"

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: protoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_remote" "test" {
						directory = "%s"
						name      = "%s"
						urls      = ["%s"]
					}
				`, directory, remote, url1),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("git_remote.test", "directory", directory),
					resource.TestCheckResourceAttr("git_remote.test", "id", remote),
					resource.TestCheckResourceAttr("git_remote.test", "name", remote),
					resource.TestCheckResourceAttr("git_remote.test", "urls.#", "1"),
					resource.TestCheckResourceAttr("git_remote.test", "urls.0", url1),
				),
			},
			{
				Config: fmt.Sprintf(`
					resource "git_remote" "test" {
						directory = "%s"
						name      = "%s"
						urls      = ["%s", "%s"]
					}
				`, directory, remote, url1, url2),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("git_remote.test", "directory", directory),
					resource.TestCheckResourceAttr("git_remote.test", "id", remote),
					resource.TestCheckResourceAttr("git_remote.test", "name", remote),
					resource.TestCheckResourceAttr("git_remote.test", "urls.#", "2"),
					resource.TestCheckResourceAttr("git_remote.test", "urls.0", url1),
					resource.TestCheckResourceAttr("git_remote.test", "urls.1", url2),
				),
			},
		},
	})
}
