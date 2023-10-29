/*
 * SPDX-FileCopyrightText: The terraform-provider-git Authors
 * SPDX-License-Identifier: 0BSD
 */

package provider_test

import (
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/metio/terraform-provider-git/internal/testutils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDataSourceGitStatuses_AllFiles(t *testing.T) {
	t.Parallel()
	directory, repository := testutils.CreateRepository(t)
	worktree := testutils.GetRepositoryWorktree(t, repository)
	testutils.AddAndCommitNewFile(t, worktree, "some-file")
	testutils.AddAndCommitNewFile(t, worktree, "other-file")
	testutils.WriteFileInWorktree(t, worktree, "new-file")

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../data-sources/git_statuses/all_files",
		Vars: map[string]interface{}{
			"directory": directory,
		},
		NoColor: true,
	})

	defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApplyAndIdempotent(t, terraformOptions)

	isClean := terraform.Output(t, terraformOptions, "is_clean")
	files := terraform.OutputMap(t, terraformOptions, "files")

	assert.Equal(t, "false", isClean)
	assert.NotNil(t, files["some-file"], "some-file")
	assert.NotNil(t, files["other-file"], "other-file")
	assert.NotNil(t, files["new-file"], "new-file")
}
