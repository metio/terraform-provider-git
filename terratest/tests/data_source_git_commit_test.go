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

func TestDataSourceGitCommit_CurrentHead(t *testing.T) {
	t.Parallel()
	directory, repository := testutils.CreateRepository(t)
	worktree := testutils.GetRepositoryWorktree(t, repository)
	testutils.AddAndCommitNewFile(t, worktree, "some-file")

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../data-sources/git_commit/current_head",
		Vars: map[string]any{
			"directory": directory,
		},
		NoColor: true,
	})

	defer terraform.DestroyContext(t, t.Context(), terraformOptions)
	terraform.InitAndApplyAndIdempotentContext(t, t.Context(), terraformOptions)

	actualDirectory := terraform.OutputContext(t, t.Context(), terraformOptions, "directory")
	actualId := terraform.OutputContext(t, t.Context(), terraformOptions, "id")
	actualRevision := terraform.OutputContext(t, t.Context(), terraformOptions, "revision")
	actualMessage := terraform.OutputContext(t, t.Context(), terraformOptions, "message")

	assert.Equal(t, directory, actualDirectory)
	assert.Equal(t, "HEAD", actualId)
	assert.Equal(t, "HEAD", actualRevision)
	assert.Equal(t, "example go-git commit", actualMessage)
}

func TestDataSourceGitCommit_FromLog(t *testing.T) {
	t.Parallel()
	directory, repository := testutils.CreateRepository(t)
	worktree := testutils.GetRepositoryWorktree(t, repository)
	testutils.AddAndCommitNewFile(t, worktree, "some-file")

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../data-sources/git_commit/from_log",
		Vars: map[string]any{
			"directory": directory,
		},
		NoColor: true,
	})

	defer terraform.DestroyContext(t, t.Context(), terraformOptions)
	terraform.InitAndApplyAndIdempotentContext(t, t.Context(), terraformOptions)

	fromLog := terraform.OutputMapContext(t, t.Context(), terraformOptions, "from_log")

	assert.NotNil(t, fromLog, "fromLog")
}
