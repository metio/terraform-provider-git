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

func TestDataSourceGitBranches_AllBranches(t *testing.T) {
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
		TerraformDir: "../data-sources/git_branches/all_branches",
		Vars: map[string]interface{}{
			"directory": directory,
		},
		NoColor: true,
	})

	defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApplyAndIdempotent(t, terraformOptions)

	everyBranch := terraform.OutputMap(t, terraformOptions, "all_branches")

	assert.NotNil(t, everyBranch["master"], "master")
	assert.NotNil(t, everyBranch[name], name)
}
