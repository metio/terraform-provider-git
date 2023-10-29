/*
 * SPDX-FileCopyrightText: The terraform-provider-git Authors
 * SPDX-License-Identifier: 0BSD
 */

package provider_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/metio/terraform-provider-git/internal/testutils"
	"path/filepath"
	"regexp"
	"testing"
)

func TestResourceGitPush(t *testing.T) {
	t.Parallel()
	directory, repository := testutils.CreateRepository(t)
	directory2 := testutils.CreateBareRepository(t)
	testutils.TestConfig(t, repository)
	worktree := testutils.GetRepositoryWorktree(t, repository)
	testutils.AddAndCommitNewFile(t, worktree, "some-file")
	testutils.CreateRemoteWithUrls(t, repository, "origin", []string{filepath.FromSlash(directory2)})
	head := testutils.GetRepositoryHead(t, repository)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_push" "test" {
						directory = "%s"
						refspecs  = ["%s"]
					}
					data "git_repository" "second" {
						directory = "%s"
						depends_on = [git_push.test]
					}
				`, directory, "refs/heads/master:refs/heads/master", directory2),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("git_push.test", "directory", directory),
					resource.TestCheckResourceAttrWith("git_push.test", "id", testutils.CheckMinLength(1)),
					resource.TestCheckResourceAttr("git_push.test", "refspecs.#", "1"),
					resource.TestCheckResourceAttr("git_push.test", "refspecs.0", "refs/heads/master:refs/heads/master"),
					resource.TestCheckResourceAttr("git_push.test", "prune", "false"),
					resource.TestCheckResourceAttr("git_push.test", "force", "false"),
					resource.TestCheckResourceAttr("data.git_repository.second", "sha1", head.Hash().String()),
				),
			},
		},
	})
}

func TestResourceGitPush_RefSpecs_Invalid(t *testing.T) {
	t.Parallel()
	directory, _ := testutils.CreateRepository(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_push" "test" {
						directory = "%s"
						refspecs  = ["%s"]
					}
				`, directory, "master"),
				ExpectError: regexp.MustCompile(`malformed refspec, separators are wrong`),
			},
		},
	})
}

func TestResourceGitPush_Remote_Unknown(t *testing.T) {
	t.Parallel()
	directory, _ := testutils.CreateRepository(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_push" "test" {
						directory = "%s"
						refspecs  = ["%s"]
					}
				`, directory, "master:master"),
				ExpectError: regexp.MustCompile(`remote not found`),
			},
		},
	})
}

func TestResourceGitPush_Directory_Invalid(t *testing.T) {
	t.Parallel()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_push" "test" {
						directory = "does/not/exist"
						refspecs  = ["%s"]
					}
				`, "master:master"),
				ExpectError: regexp.MustCompile(`Cannot open repository`),
			},
		},
	})
}

func TestResourceGitPush_Auth_Invalid_BasicAndBearer(t *testing.T) {
	t.Parallel()
	directory, _ := testutils.CreateRepository(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_push" "test" {
						directory = "%s"
						refspecs  = ["%s"]
						auth      = {
							basic  = {
								username = "user"
								password = "secret"
							}
							bearer = "secret-token"
						}
					}
				`, directory, "refs/heads/master:refs/heads/master"),
				ExpectError: regexp.MustCompile(`Invalid Attribute Combination`),
			},
		},
	})
}

func TestResourceGitPush_Auth_Invalid_BasicAndSshKey(t *testing.T) {
	t.Parallel()
	directory, _ := testutils.CreateRepository(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_push" "test" {
						directory = "%s"
						refspecs  = ["%s"]
						auth      = {
							basic  = {
								username = "user"
								password = "secret"
							}
							ssh_key = {
								private_key_path = "/some/path/to/ssh/key"
							}
						}
					}
				`, directory, "refs/heads/master:refs/heads/master"),
				ExpectError: regexp.MustCompile(`Invalid Attribute Combination`),
			},
		},
	})
}

func TestResourceGitPush_Auth_Invalid_BasicAndSshAgent(t *testing.T) {
	t.Parallel()
	directory, _ := testutils.CreateRepository(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_push" "test" {
						directory = "%s"
						refspecs  = ["%s"]
						auth      = {
							basic  = {
								username = "user"
								password = "secret"
							}
							ssh_agent = {}
						}
					}
				`, directory, "refs/heads/master:refs/heads/master"),
				ExpectError: regexp.MustCompile(`Invalid Attribute Combination`),
			},
		},
	})
}

func TestResourceGitPush_Auth_Invalid_BasicAndSshPassword(t *testing.T) {
	t.Parallel()
	directory, _ := testutils.CreateRepository(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_push" "test" {
						directory = "%s"
						refspecs  = ["%s"]
						auth      = {
							basic  = {
								username = "user"
								password = "secret"
							}
							ssh_password = {
								username = "user"
								password = "pass"
							}
						}
					}
				`, directory, "refs/heads/master:refs/heads/master"),
				ExpectError: regexp.MustCompile(`Invalid Attribute Combination`),
			},
		},
	})
}

func TestResourceGitPush_Auth_Invalid_BearerAndSshKey(t *testing.T) {
	t.Parallel()
	directory, _ := testutils.CreateRepository(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_push" "test" {
						directory = "%s"
						refspecs  = ["%s"]
						auth      = {
							bearer  = "token"
							ssh_key = {
								private_key_path = "/some/path/to/ssh/key"
							}
						}
					}
				`, directory, "refs/heads/master:refs/heads/master"),
				ExpectError: regexp.MustCompile(`Invalid Attribute Combination`),
			},
		},
	})
}

func TestResourceGitPush_Auth_Invalid_BearerAndSshAgent(t *testing.T) {
	t.Parallel()
	directory, _ := testutils.CreateRepository(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_push" "test" {
						directory = "%s"
						refspecs  = ["%s"]
						auth      = {
							bearer  = "token"
							ssh_agent = {}
						}
					}
				`, directory, "refs/heads/master:refs/heads/master"),
				ExpectError: regexp.MustCompile(`Invalid Attribute Combination`),
			},
		},
	})
}

func TestResourceGitPush_Auth_Invalid_BearerAndSshPassword(t *testing.T) {
	t.Parallel()
	directory, _ := testutils.CreateRepository(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_push" "test" {
						directory = "%s"
						refspecs  = ["%s"]
						auth      = {
							bearer  = "token"
							ssh_password = {
								username = "user"
								password = "pass"
							}
						}
					}
				`, directory, "refs/heads/master:refs/heads/master"),
				ExpectError: regexp.MustCompile(`Invalid Attribute Combination`),
			},
		},
	})
}

func TestResourceGitPush_Auth_Invalid_SshKeyAndSshAgent(t *testing.T) {
	t.Parallel()
	directory, _ := testutils.CreateRepository(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_push" "test" {
						directory = "%s"
						refspecs  = ["%s"]
						auth      = {
							ssh_key = {
								private_key_path = "/some/path/to/ssh/key"
							}
							ssh_agent = {}
						}
					}
				`, directory, "refs/heads/master:refs/heads/master"),
				ExpectError: regexp.MustCompile(`Invalid Attribute Combination`),
			},
		},
	})
}

func TestResourceGitPush_Auth_Invalid_SshKeyAndSshPassword(t *testing.T) {
	t.Parallel()
	directory, _ := testutils.CreateRepository(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_push" "test" {
						directory = "%s"
						refspecs  = ["%s"]
						auth      = {
							ssh_key = {
								private_key_path = "/some/path/to/ssh/key"
							}
							ssh_password = {
								username = "user"
								password = "pass"
							}
						}
					}
				`, directory, "refs/heads/master:refs/heads/master"),
				ExpectError: regexp.MustCompile(`Invalid Attribute Combination`),
			},
		},
	})
}

func TestResourceGitPush_Auth_Invalid_BasicAuth_Username_Required(t *testing.T) {
	t.Parallel()
	directory, _ := testutils.CreateRepository(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_push" "test" {
						directory = "%s"
						refspecs  = ["%s"]
						auth      = {
							basic = {
								password = "pass"
							}
						}
					}
				`, directory, "refs/heads/master:refs/heads/master"),
				ExpectError: regexp.MustCompile(`"username" is required.`),
			},
		},
	})
}

func TestResourceGitPush_Auth_Invalid_BasicAuth_Password_Required(t *testing.T) {
	t.Parallel()
	directory, _ := testutils.CreateRepository(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_push" "test" {
						directory = "%s"
						refspecs  = ["%s"]
						auth      = {
							basic = {
								username = "user"
							}
						}
					}
				`, directory, "refs/heads/master:refs/heads/master"),
				ExpectError: regexp.MustCompile(`"password" is required.`),
			},
		},
	})
}

func TestResourceGitPush_Auth_Invalid_Bearer_Token_Required(t *testing.T) {
	t.Parallel()
	directory, _ := testutils.CreateRepository(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_push" "test" {
						directory = "%s"
						refspecs  = ["%s"]
						auth      = {
							bearer = ""
						}
					}
				`, directory, "refs/heads/master:refs/heads/master"),
				ExpectError: regexp.MustCompile(`Invalid Attribute Value Length`),
			},
		},
	})
}

func TestResourceGitPush_Auth_Invalid_SshKey_PrivateKey_Required(t *testing.T) {
	t.Parallel()
	directory, _ := testutils.CreateRepository(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_push" "test" {
						directory = "%s"
						refspecs  = ["%s"]
						auth      = {
							ssh_key = {}
						}
					}
				`, directory, "refs/heads/master:refs/heads/master"),
				ExpectError: regexp.MustCompile(`Invalid Attribute Combination`),
			},
		},
	})
}

func TestResourceGitPush_Auth_Invalid_SshKey_PathAndPem(t *testing.T) {
	t.Parallel()
	directory, _ := testutils.CreateRepository(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_push" "test" {
						directory = "%s"
						refspecs  = ["%s"]
						auth      = {
							ssh_key = {
								private_key_path = "/path/to/key"
								private_key_pem  = "PEM data here"
							}
						}
					}
				`, directory, "refs/heads/master:refs/heads/master"),
				ExpectError: regexp.MustCompile(`Invalid Attribute Combination`),
			},
		},
	})
}

func TestResourceGitPush_Auth_Invalid_SshPassword_Username_Required(t *testing.T) {
	t.Parallel()
	directory, _ := testutils.CreateRepository(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_push" "test" {
						directory = "%s"
						refspecs  = ["%s"]
						auth      = {
							ssh_password = {
								password = "pass"
							}
						}
					}
				`, directory, "refs/heads/master:refs/heads/master"),
				ExpectError: regexp.MustCompile(`"username" is required.`),
			},
		},
	})
}

func TestResourceGitPush_Auth_Invalid_SshPassword_Password_Required(t *testing.T) {
	t.Parallel()
	directory, _ := testutils.CreateRepository(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_push" "test" {
						directory = "%s"
						refspecs  = ["%s"]
						auth      = {
							ssh_password = {
								username = "user"
							}
						}
					}
				`, directory, "refs/heads/master:refs/heads/master"),
				ExpectError: regexp.MustCompile(`"password" is required.`),
			},
		},
	})
}

func TestResourceGitPush_Auth_Invalid_Empty(t *testing.T) {
	t.Parallel()
	directory, _ := testutils.CreateRepository(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testutils.ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "git_push" "test" {
						directory = "%s"
						refspecs  = ["%s"]
						auth      = {}
					}
				`, directory, "refs/heads/master:refs/heads/master"),
				ExpectError: regexp.MustCompile(`Invalid Attribute Combination`),
			},
		},
	})
}
