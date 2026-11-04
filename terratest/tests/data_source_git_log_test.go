/*
 * SPDX-FileCopyrightText: The terraform-provider-git Authors
 * SPDX-License-Identifier: 0BSD
 */

package provider_test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/metio/terraform-provider-git/internal/testutils"
	"github.com/stretchr/testify/assert"
)

func TestDataSourceGitLog_CurrentLog(t *testing.T) {
	t.Parallel()
	directory, repository := testutils.CreateRepository(t)
	worktree := testutils.GetRepositoryWorktree(t, repository)
	testutils.AddAndCommitNewFile(t, worktree, "some-file")

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../data-sources/git_log/current_log",
		Vars: map[string]any{
			"directory": directory,
		},
		NoColor: true,
	})

	defer terraform.DestroyContext(t, t.Context(), terraformOptions)
	terraform.InitAndApplyAndIdempotentContext(t, t.Context(), terraformOptions)

	actualDirectory := terraform.OutputContext(t, t.Context(), terraformOptions, "directory")
	actualId := terraform.OutputContext(t, t.Context(), terraformOptions, "id")
	commits := terraform.OutputListContext(t, t.Context(), terraformOptions, "commits")

	assert.Equal(t, directory, actualDirectory, "actualDirectory")
	assert.Equal(t, directory, actualId, "actualId")
	assert.Equal(t, 1, len(commits), "commits")
}

func TestDataSourceGitLog_Filtered(t *testing.T) {
	t.Parallel()
	directory, repository := testutils.CreateRepository(t)
	worktree := testutils.GetRepositoryWorktree(t, repository)
	filename := "some-file"
	testutils.AddAndCommitNewFile(t, worktree, filename)

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../data-sources/git_log/filtered",
		Vars: map[string]any{
			"directory":   directory,
			"filter_path": filename,
		},
		NoColor: true,
	})

	defer terraform.DestroyContext(t, t.Context(), terraformOptions)
	terraform.InitAndApplyAndIdempotentContext(t, t.Context(), terraformOptions)

	actualDirectory := terraform.OutputContext(t, t.Context(), terraformOptions, "directory")
	actualId := terraform.OutputContext(t, t.Context(), terraformOptions, "id")
	commits := terraform.OutputListContext(t, t.Context(), terraformOptions, "commits")

	assert.Equal(t, directory, actualDirectory, "actualDirectory")
	assert.Equal(t, directory, actualId, "actualId")
	assert.Equal(t, 1, len(commits), "commits")
}
