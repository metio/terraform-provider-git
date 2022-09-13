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
	"regexp"
	"testing"
)

func TestResourceGitRemote(t *testing.T) {
	t.Parallel()
	directory, _ := testutils.CreateRepository(t)
	defer os.RemoveAll(directory)
	name := "some-name"
	url1 := "https://github.com/some-org/some-repo.git"

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
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
					resource.TestCheckResourceAttr("git_remote.test", "id", fmt.Sprintf("%s|%s", directory, name)),
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
	directory, _ := testutils.CreateRepository(t)
	defer os.RemoveAll(directory)
	name := "some-name"
	url1 := "https://github.com/some-org/some-repo.git"
	url2 := "https://codeberg.org/some-org/some-repo.git"

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
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
					resource.TestCheckResourceAttr("git_remote.test", "id", fmt.Sprintf("%s|%s", directory, name)),
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
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
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
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
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
	directory, _ := testutils.CreateRepository(t)
	defer os.RemoveAll(directory)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
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
	directory, _ := testutils.CreateRepository(t)
	defer os.RemoveAll(directory)
	name := "some-name"

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
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
	directory, _ := testutils.CreateRepository(t)
	defer os.RemoveAll(directory)
	name := "some-name"
	url := "https://github.com/some-org/some-repo.git"

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_remote" "test" {
						directory = "%s"
						name      = "%s"
						urls      = ["%s"]
					}
				`, directory, name, url),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("git_remote.test", "directory", directory),
					resource.TestCheckResourceAttr("git_remote.test", "id", fmt.Sprintf("%s|%s", directory, name)),
					resource.TestCheckResourceAttr("git_remote.test", "name", name),
					resource.TestCheckResourceAttr("git_remote.test", "urls.#", "1"),
					resource.TestCheckResourceAttr("git_remote.test", "urls.0", url),
				),
			},
			{
				ResourceName:      "git_remote.test",
				ImportState:       true,
				ImportStateId:     fmt.Sprintf("%s|%s", directory, name),
				ImportStateVerify: true,
			},
		},
	})
}

func TestResourceGitRemote_ImportMultipleUrls(t *testing.T) {
	t.Parallel()
	directory, _ := testutils.CreateRepository(t)
	defer os.RemoveAll(directory)
	name := "some-name"
	url1 := "https://github.com/some-org/some-repo.git"
	url2 := "https://codeberg.org/some-org/some-repo.git"

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
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
					resource.TestCheckResourceAttr("git_remote.test", "id", fmt.Sprintf("%s|%s", directory, name)),
					resource.TestCheckResourceAttr("git_remote.test", "name", name),
					resource.TestCheckResourceAttr("git_remote.test", "urls.#", "2"),
					resource.TestCheckResourceAttr("git_remote.test", "urls.0", url1),
					resource.TestCheckResourceAttr("git_remote.test", "urls.1", url2),
				),
			},
			{
				ResourceName:      "git_remote.test",
				ImportState:       true,
				ImportStateId:     fmt.Sprintf("%s|%s", directory, name),
				ImportStateVerify: true,
			},
		},
	})
}

func TestResourceGitRemote_Update_Urls(t *testing.T) {
	t.Parallel()
	directory, _ := testutils.CreateRepository(t)
	defer os.RemoveAll(directory)
	name := "some-name"
	url1 := "https://github.com/some-org/some-repo.git"
	url2 := "https://codeberg.org/some-org/some-repo.git"

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
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
					resource.TestCheckResourceAttr("git_remote.test", "id", fmt.Sprintf("%s|%s", directory, name)),
					resource.TestCheckResourceAttr("git_remote.test", "name", name),
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
				`, directory, name, url1, url2),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("git_remote.test", "directory", directory),
					resource.TestCheckResourceAttr("git_remote.test", "id", fmt.Sprintf("%s|%s", directory, name)),
					resource.TestCheckResourceAttr("git_remote.test", "name", name),
					resource.TestCheckResourceAttr("git_remote.test", "urls.#", "2"),
					resource.TestCheckResourceAttr("git_remote.test", "urls.0", url1),
					resource.TestCheckResourceAttr("git_remote.test", "urls.1", url2),
				),
			},
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
					resource.TestCheckResourceAttr("git_remote.test", "id", fmt.Sprintf("%s|%s", directory, name)),
					resource.TestCheckResourceAttr("git_remote.test", "name", name),
					resource.TestCheckResourceAttr("git_remote.test", "urls.#", "1"),
					resource.TestCheckResourceAttr("git_remote.test", "urls.0", url1),
				),
			},
		},
	})
}

func TestResourceGitRemote_Update_Name(t *testing.T) {
	t.Parallel()
	directory, _ := testutils.CreateRepository(t)
	defer os.RemoveAll(directory)
	name := "some-name"
	newName := "other-name"
	url1 := "https://github.com/some-org/some-repo.git"

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
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
					resource.TestCheckResourceAttr("git_remote.test", "id", fmt.Sprintf("%s|%s", directory, name)),
					resource.TestCheckResourceAttr("git_remote.test", "name", name),
					resource.TestCheckResourceAttr("git_remote.test", "urls.#", "1"),
					resource.TestCheckResourceAttr("git_remote.test", "urls.0", url1),
				),
			},
			{
				Config: fmt.Sprintf(`
					resource "git_remote" "test" {
						directory = "%s"
						name      = "%s"
						urls      = ["%s"]
					}
				`, directory, newName, url1),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("git_remote.test", "directory", directory),
					resource.TestCheckResourceAttr("git_remote.test", "id", fmt.Sprintf("%s|%s", directory, newName)),
					resource.TestCheckResourceAttr("git_remote.test", "name", newName),
					resource.TestCheckResourceAttr("git_remote.test", "urls.#", "1"),
					resource.TestCheckResourceAttr("git_remote.test", "urls.0", url1),
				),
			},
		},
	})
}

func TestResourceGitRemote_Update_Directory(t *testing.T) {
	t.Parallel()
	directory1, _ := testutils.CreateRepository(t)
	directory2, _ := testutils.CreateRepository(t)
	defer os.RemoveAll(directory1)
	defer os.RemoveAll(directory2)
	name := "some-name"
	url1 := "https://github.com/some-org/some-repo.git"

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_remote" "test" {
						directory = "%s"
						name      = "%s"
						urls      = ["%s"]
					}
				`, directory1, name, url1),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("git_remote.test", "directory", directory1),
					resource.TestCheckResourceAttr("git_remote.test", "id", fmt.Sprintf("%s|%s", directory1, name)),
					resource.TestCheckResourceAttr("git_remote.test", "name", name),
					resource.TestCheckResourceAttr("git_remote.test", "urls.#", "1"),
					resource.TestCheckResourceAttr("git_remote.test", "urls.0", url1),
				),
			},
			{
				Config: fmt.Sprintf(`
					resource "git_remote" "test" {
						directory = "%s"
						name      = "%s"
						urls      = ["%s"]
					}
				`, directory2, name, url1),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("git_remote.test", "directory", directory2),
					resource.TestCheckResourceAttr("git_remote.test", "id", fmt.Sprintf("%s|%s", directory2, name)),
					resource.TestCheckResourceAttr("git_remote.test", "name", name),
					resource.TestCheckResourceAttr("git_remote.test", "urls.#", "1"),
					resource.TestCheckResourceAttr("git_remote.test", "urls.0", url1),
				),
			},
		},
	})
}
