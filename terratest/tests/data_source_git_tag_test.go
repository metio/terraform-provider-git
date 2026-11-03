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

func TestDataSourceGitTag_SingleTag(t *testing.T) {
	t.Parallel()
	directory, repository := testutils.CreateRepository(t)
	worktree := testutils.GetRepositoryWorktree(t, repository)
	testutils.AddAndCommitNewFile(t, worktree, "some-file")
	name := "some-tag"
	testutils.CreateTag(t, repository, name)

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../data-sources/git_tag/single_tag",
		Vars: map[string]interface{}{
			"directory": directory,
			"name":      name,
		},
		NoColor: true,
	})

	defer terraform.DestroyContext(t, t.Context(), terraformOptions)
	terraform.InitAndApplyAndIdempotentContext(t, t.Context(), terraformOptions)

	actualDirectory := terraform.OutputContext(t, t.Context(), terraformOptions, "directory")
	actualId := terraform.OutputContext(t, t.Context(), terraformOptions, "id")
	actualName := terraform.OutputContext(t, t.Context(), terraformOptions, "name")

	assert.Equal(t, directory, actualDirectory, "actualDirectory")
	assert.Equal(t, name, actualId, "actualId")
	assert.Equal(t, name, actualName, "actualName")
}

func TestDataSourceGitTag_EveryTag(t *testing.T) {
	t.Parallel()
	directory, repository := testutils.CreateRepository(t)
	worktree := testutils.GetRepositoryWorktree(t, repository)
	testutils.AddAndCommitNewFile(t, worktree, "some-file")
	testutils.CreateTag(t, repository, "some-tag")
	testutils.CreateTag(t, repository, "other-tag")

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../data-sources/git_tag/every_tag",
		Vars: map[string]interface{}{
			"directory": directory,
		},
		NoColor: true,
	})

	defer terraform.DestroyContext(t, t.Context(), terraformOptions)
	terraform.InitAndApplyAndIdempotentContext(t, t.Context(), terraformOptions)

	everyTag := terraform.OutputMapContext(t, t.Context(), terraformOptions, "every_tag")

	assert.NotNil(t, everyTag["some-tag"], "some-tag")
	assert.NotNil(t, everyTag["other-tag"], "other-tag")
}
