/*
 * SPDX-FileCopyrightText: The terraform-provider-git Authors
 * SPDX-License-Identifier: 0BSD
 */

package provider_test

import (
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/metio/terraform-provider-git/internal/testutils"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestDataSourceGitTag_SingleTag(t *testing.T) {
	t.Parallel()
	directory, repository := testutils.CreateRepository(t)
	defer os.RemoveAll(directory)
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

	defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApplyAndIdempotent(t, terraformOptions)

	actualDirectory := terraform.Output(t, terraformOptions, "directory")
	actualId := terraform.Output(t, terraformOptions, "id")
	actualName := terraform.Output(t, terraformOptions, "name")

	assert.Equal(t, directory, actualDirectory, "actualDirectory")
	assert.Equal(t, name, actualId, "actualId")
	assert.Equal(t, name, actualName, "actualName")
}

func TestDataSourceGitTag_EveryTag(t *testing.T) {
	t.Parallel()
	directory, repository := testutils.CreateRepository(t)
	defer os.RemoveAll(directory)
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

	defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApplyAndIdempotent(t, terraformOptions)

	everyTag := terraform.OutputMap(t, terraformOptions, "every_tag")

	assert.NotNil(t, everyTag["some-tag"], "some-tag")
	assert.NotNil(t, everyTag["other-tag"], "other-tag")
}
