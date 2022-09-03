/*
 * SPDX-FileCopyrightText: The terraform-provider-git Authors
 * SPDX-License-Identifier: 0BSD
 */

package provider_test

import (
	"github.com/go-git/go-git/v5/config"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/metio/terraform-provider-git/internal/testutils"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestDataSourceGitBranch_BranchOnly(t *testing.T) {
	t.Parallel()
	directory, repository := testutils.CreateRepository(t)
	defer os.RemoveAll(directory)
	worktree := testutils.GetRepositoryWorktree(t, repository)
	testutils.AddAndCommitNewFile(t, worktree, "some-file")
	name := "name-of-branch"
	remote := "origin"
	rebase := "true"
	testutils.CreateBranch(t, repository, &config.Branch{
		Name:   name,
		Remote: remote,
		Rebase: rebase,
	})

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../data-sources/git_branch/branch_only",
		Vars: map[string]interface{}{
			"directory": directory,
			"name":      name,
		},
		NoColor: true,
	})

	defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApplyAndIdempotent(t, terraformOptions)

	actualDirectory := terraform.Output(t, terraformOptions, "directory")
	actualId := terraform.Output(t, terraformOptions, "id")
	actualName := terraform.Output(t, terraformOptions, "name")
	actualRemote := terraform.Output(t, terraformOptions, "remote")

	assert.Equal(t, directory, actualDirectory)
	assert.Equal(t, name, actualId)
	assert.Equal(t, name, actualName)
	assert.Equal(t, remote, actualRemote)
}

func TestDataSourceGitBranch_BranchFromRepo(t *testing.T) {
	t.Parallel()
	directory, repository := testutils.CreateRepository(t)
	defer os.RemoveAll(directory)
	worktree := testutils.GetRepositoryWorktree(t, repository)
	testutils.AddAndCommitNewFile(t, worktree, "some-file")

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../data-sources/git_branch/branch_from_repo",
		Vars: map[string]interface{}{
			"directory": directory,
		},
		NoColor: true,
	})

	defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApplyAndIdempotent(t, terraformOptions)

	actualDirectory := terraform.Output(t, terraformOptions, "directory")
	actualId := terraform.Output(t, terraformOptions, "id")
	actualName := terraform.Output(t, terraformOptions, "name")

	assert.Equal(t, directory, actualDirectory)
	assert.Equal(t, "master", actualId)
	assert.Equal(t, "master", actualName)
}

func TestDataSourceGitBranch_EveryBranch(t *testing.T) {
	t.Parallel()
	directory, repository := testutils.CreateRepository(t)
	defer os.RemoveAll(directory)
	worktree := testutils.GetRepositoryWorktree(t, repository)
	testutils.AddAndCommitNewFile(t, worktree, "some-file")
	name := "name-of-branch"
	remote := "origin"
	rebase := "true"
	testutils.CreateBranch(t, repository, &config.Branch{
		Name:   name,
		Remote: remote,
		Rebase: rebase,
	})

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../data-sources/git_branch/every_branch",
		Vars: map[string]interface{}{
			"directory": directory,
		},
		NoColor: true,
	})

	defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApplyAndIdempotent(t, terraformOptions)

	everyBranch := terraform.OutputMap(t, terraformOptions, "every_branch")

	assert.NotNil(t, everyBranch["master"], "master")
	assert.NotNil(t, everyBranch[name], name)
}
