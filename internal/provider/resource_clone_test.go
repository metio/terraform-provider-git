/*
 * SPDX-FileCopyrightText: The terraform-provider-git Authors
 * SPDX-License-Identifier: 0BSD
 */

package provider_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/metio/terraform-provider-git/internal/testutils"
)

func TestResourceGitClone(t *testing.T) {
	t.Parallel()
	directory := testutils.TemporaryDirectory(t)
	localRepository, repository := testutils.CreateRepository(t)
	testutils.TestConfig(t, repository)
	worktree := testutils.GetRepositoryWorktree(t, repository)
	name := "some-file"
	testutils.WriteFileInWorktree(t, worktree, name)
	testutils.GitAdd(t, worktree, name)
	hash := testutils.GitCommit(t, worktree)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_clone" "test" {
						directory      = "%s"
						url       	   = "%s"
						reference_name = "master"
					}
				`, directory, localRepository),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("git_clone.test", tfjsonpath.New("directory"), knownvalue.StringExact(directory)),
					statecheck.ExpectKnownValue("git_clone.test", tfjsonpath.New("id"), knownvalue.StringExact(directory)),
					statecheck.ExpectKnownValue("git_clone.test", tfjsonpath.New("url"), knownvalue.StringExact(localRepository)),
					statecheck.ExpectKnownValue("git_clone.test", tfjsonpath.New("sha1"), knownvalue.StringExact(hash.String())),
					statecheck.ExpectKnownValue("git_clone.test", tfjsonpath.New("remote_name"), knownvalue.StringExact("origin")),
					statecheck.ExpectKnownValue("git_clone.test", tfjsonpath.New("bare"), knownvalue.Bool(false)),
				},
			},
		},
	})
}

func TestResourceGitClone_Bare(t *testing.T) {
	t.Parallel()
	directory := testutils.TemporaryDirectory(t)
	localRepository, repository := testutils.CreateRepository(t)
	testutils.TestConfig(t, repository)
	worktree := testutils.GetRepositoryWorktree(t, repository)
	name := "some-file"
	testutils.WriteFileInWorktree(t, worktree, name)
	testutils.GitAdd(t, worktree, name)
	hash := testutils.GitCommit(t, worktree)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_clone" "test" {
						directory 	   = "%s"
						url       	   = "%s"
						bare      	   = true
						reference_name = "master"
					}
				`, directory, localRepository),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("git_clone.test", tfjsonpath.New("bare"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue("git_clone.test", tfjsonpath.New("directory"), knownvalue.StringExact(directory)),
					statecheck.ExpectKnownValue("git_clone.test", tfjsonpath.New("id"), knownvalue.StringExact(directory)),
					statecheck.ExpectKnownValue("git_clone.test", tfjsonpath.New("url"), knownvalue.StringExact(localRepository)),
					statecheck.ExpectKnownValue("git_clone.test", tfjsonpath.New("sha1"), knownvalue.StringExact(hash.String())),
					statecheck.ExpectKnownValue("git_clone.test", tfjsonpath.New("remote_name"), knownvalue.StringExact("origin")),
				},
			},
		},
	})
}

func TestResourceGitClone_MultipleCommits(t *testing.T) {
	t.Parallel()
	directory := testutils.TemporaryDirectory(t)
	localRepository, repository := testutils.CreateRepository(t)
	testutils.TestConfig(t, repository)
	worktree := testutils.GetRepositoryWorktree(t, repository)
	name := "some-file"
	testutils.WriteFileInWorktree(t, worktree, name)
	testutils.GitAdd(t, worktree, name)
	hash := testutils.GitCommit(t, worktree)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_clone" "test" {
						directory      = "%s"
						url            = "%s"
						reference_name = "master"
					}
				`, directory, localRepository),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("git_clone.test", tfjsonpath.New("directory"), knownvalue.StringExact(directory)),
					statecheck.ExpectKnownValue("git_clone.test", tfjsonpath.New("id"), knownvalue.StringExact(directory)),
					statecheck.ExpectKnownValue("git_clone.test", tfjsonpath.New("url"), knownvalue.StringExact(localRepository)),
					statecheck.ExpectKnownValue("git_clone.test", tfjsonpath.New("sha1"), knownvalue.StringExact(hash.String())),
					statecheck.ExpectKnownValue("git_clone.test", tfjsonpath.New("remote_name"), knownvalue.StringExact("origin")),
					statecheck.ExpectKnownValue("git_clone.test", tfjsonpath.New("reference_name"), knownvalue.StringExact("master")),
					statecheck.ExpectKnownValue("git_clone.test", tfjsonpath.New("bare"), knownvalue.Bool(false)),
				},
			},
			{
				PreConfig: func() {
					otherName := "other-file"
					testutils.WriteFileInWorktree(t, worktree, otherName)
					testutils.GitAdd(t, worktree, otherName)
					testutils.GitCommit(t, worktree)
				},
				Config: fmt.Sprintf(`
					resource "git_clone" "test" {
						directory      = "%s"
						url            = "%s"
						reference_name = "master"
					}
				`, directory, localRepository),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("git_clone.test", tfjsonpath.New("sha1"), knownvalue.StringFunc(testutils.CheckExactLength(40))),
				},
			},
		},
	})
}

func TestResourceGitClone_Directory_Missing(t *testing.T) {
	t.Parallel()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: `
					resource "git_clone" "test" {
						url = "some-url"
					}
				`,
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
		},
	})
}

func TestResourceGitClone_URL_Missing(t *testing.T) {
	t.Parallel()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: `
					resource "git_clone" "test" {
						directory = "some-directory"
					}
				`,
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
		},
	})
}
