/*
 * SPDX-FileCopyrightText: The terraform-provider-git Authors
 * SPDX-License-Identifier: 0BSD
 */

package provider_test

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/go-git/go-git/v5/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/metio/terraform-provider-git/internal/testutils"
)

func TestDataSourceGitRemote(t *testing.T) {
	t.Parallel()
	directory, repository := testutils.CreateRepository(t)
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
					resource.TestCheckNoResourceAttr("data.git_remote.test", "head_branch"),
				),
			},
		},
	})
}

func TestDataSourceGitRemote_HeadBranch(t *testing.T) {
	t.Parallel()
	directory, repository := testutils.CreateRepository(t)
	remote := "example"
	testutils.CreateRemote(t, repository, remote)
	testutils.SetRemoteHead(t, repository, remote, "feature/some-branch")

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
					resource.TestCheckResourceAttr("data.git_remote.test", "name", remote),
					resource.TestCheckResourceAttr("data.git_remote.test", "head_branch", "feature/some-branch"),
				),
			},
		},
	})
}

func TestDataSourceGitRemote_QueryRemote(t *testing.T) {
	t.Parallel()
	upstreamDirectory, upstreamRepository := testutils.CreateRepository(t)
	testutils.TestConfig(t, upstreamRepository)
	worktree := testutils.GetRepositoryWorktree(t, upstreamRepository)
	testutils.AddAndCommitNewFile(t, worktree, "some-file")
	testutils.CreateBranch(t, upstreamRepository, &config.Branch{Name: "develop"})
	testutils.SetHead(t, upstreamRepository, "develop")

	directory, repository := testutils.CreateRepository(t)
	remote := "origin"
	testutils.CreateRemoteWithUrls(t, repository, remote, []string{upstreamDirectory})
	testutils.SetRemoteHead(t, repository, remote, "outdated")

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "git_remote" "test" {
						directory    = "%s"
						name         = "%s"
						query_remote = true
					}
				`, directory, remote),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.git_remote.test", "name", remote),
					resource.TestCheckResourceAttr("data.git_remote.test", "query_remote", "true"),
					resource.TestCheckResourceAttr("data.git_remote.test", "head_branch", "develop"),
				),
			},
		},
	})
}

func TestDataSourceGitRemote_QueryRemote_EmptyRepository(t *testing.T) {
	t.Parallel()
	upstreamDirectory := testutils.CreateBareRepository(t)
	directory, repository := testutils.CreateRepository(t)
	remote := "origin"
	testutils.CreateRemoteWithUrls(t, repository, remote, []string{upstreamDirectory})

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "git_remote" "test" {
						directory    = "%s"
						name         = "%s"
						query_remote = true
					}
				`, directory, remote),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.git_remote.test", "name", remote),
					resource.TestCheckNoResourceAttr("data.git_remote.test", "head_branch"),
				),
			},
		},
	})
}

func TestDataSourceGitRemote_QueryRemote_Unreachable(t *testing.T) {
	t.Parallel()
	directory, repository := testutils.CreateRepository(t)
	remote := "origin"
	testutils.CreateRemoteWithUrls(t, repository, remote, []string{filepath.Join(testutils.TemporaryDirectory(t), "does-not-exist")})

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "git_remote" "test" {
						directory    = "%s"
						name         = "%s"
						query_remote = true
					}
				`, directory, remote),
				ExpectError: regexp.MustCompile(`Cannot list remote`),
			},
		},
	})
}

func TestDataSourceGitRemote_QueryRemote_Auth_Invalid_Empty(t *testing.T) {
	t.Parallel()
	directory, repository := testutils.CreateRepository(t)
	remote := "origin"
	testutils.CreateRemote(t, repository, remote)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "git_remote" "test" {
						directory    = "%s"
						name         = "%s"
						query_remote = true
						auth         = {}
					}
				`, directory, remote),
				ExpectError: regexp.MustCompile(`Invalid Attribute Combination`),
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

func TestDataSourceGitRemote_Issue314(t *testing.T) {
	t.Parallel()
	directory, repository := testutils.CreateRepository(t)
	nested := filepath.Join(directory, "nested")
	err := os.MkdirAll(nested, os.ModePerm)
	if err != nil {
		t.Errorf("cannot create nested directory: %v", err)
	}
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
				`, filepath.ToSlash(nested), remote),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.git_remote.test", "directory", filepath.ToSlash(nested)),
					resource.TestCheckResourceAttr("data.git_remote.test", "id", remote),
					resource.TestCheckResourceAttr("data.git_remote.test", "name", remote),
					resource.TestCheckResourceAttr("data.git_remote.test", "urls.#", "1"),
				),
			},
		},
	})
}
