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

func TestDataSourceGitStatus_SingleFile_Committed(t *testing.T) {
	directory, repository := testutils.CreateRepository(t)
	defer os.RemoveAll(directory)
	worktree := testutils.GetRepositoryWorktree(t, repository)
	filename := "some-file"
	testutils.AddAndCommitNewFile(t, worktree, filename)

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../data-sources/git_status/single_file",
		Vars: map[string]interface{}{
			"directory": directory,
			"file":      filename,
		},
		NoColor: true,
	})

	defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApplyAndIdempotent(t, terraformOptions)

	actualDirectory := terraform.Output(t, terraformOptions, "directory")
	actualId := terraform.Output(t, terraformOptions, "id")
	actualFile := terraform.Output(t, terraformOptions, "file")
	actualStaging := terraform.Output(t, terraformOptions, "staging")
	actualWorktree := terraform.Output(t, terraformOptions, "worktree")

	assert.Equal(t, directory, actualDirectory, "actualDirectory")
	assert.Equal(t, filename, actualId, "actualId")
	assert.Equal(t, filename, actualFile, "actualFile")
	assert.Equal(t, "?", actualStaging, "actualStaging")
	assert.Equal(t, "?", actualWorktree, "actualWorktree")
}

func TestDataSourceGitStatus_SingleFile_Written(t *testing.T) {
	directory, repository := testutils.CreateRepository(t)
	defer os.RemoveAll(directory)
	worktree := testutils.GetRepositoryWorktree(t, repository)
	filename := "some-file"
	testutils.WriteFileInWorktree(t, worktree, filename)

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../data-sources/git_status/single_file",
		Vars: map[string]interface{}{
			"directory": directory,
			"file":      filename,
		},
		NoColor: true,
	})

	defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApplyAndIdempotent(t, terraformOptions)

	actualDirectory := terraform.Output(t, terraformOptions, "directory")
	actualId := terraform.Output(t, terraformOptions, "id")
	actualFile := terraform.Output(t, terraformOptions, "file")
	actualStaging := terraform.Output(t, terraformOptions, "staging")
	actualWorktree := terraform.Output(t, terraformOptions, "worktree")

	assert.Equal(t, directory, actualDirectory, "actualDirectory")
	assert.Equal(t, filename, actualId, "actualId")
	assert.Equal(t, filename, actualFile, "actualFile")
	assert.Equal(t, "?", actualStaging, "actualStaging")
	assert.Equal(t, "?", actualWorktree, "actualWorktree")
}

func TestDataSourceGitStatus_SingleFile_Added(t *testing.T) {
	directory, repository := testutils.CreateRepository(t)
	defer os.RemoveAll(directory)
	worktree := testutils.GetRepositoryWorktree(t, repository)
	filename := "some-file"
	testutils.WriteFileInWorktree(t, worktree, filename)
	testutils.GitAdd(t, worktree, filename)

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../data-sources/git_status/single_file",
		Vars: map[string]interface{}{
			"directory": directory,
			"file":      filename,
		},
		NoColor: true,
	})

	defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApplyAndIdempotent(t, terraformOptions)

	actualDirectory := terraform.Output(t, terraformOptions, "directory")
	actualId := terraform.Output(t, terraformOptions, "id")
	actualFile := terraform.Output(t, terraformOptions, "file")
	actualStaging := terraform.Output(t, terraformOptions, "staging")
	actualWorktree := terraform.Output(t, terraformOptions, "worktree")

	assert.Equal(t, directory, actualDirectory, "actualDirectory")
	assert.Equal(t, filename, actualId, "actualId")
	assert.Equal(t, filename, actualFile, "actualFile")
	assert.Equal(t, "A", actualStaging, "actualStaging")
	assert.Equal(t, " ", actualWorktree, "actualWorktree")
}

func TestDataSourceGitStatus_EveryFile(t *testing.T) {
	t.Parallel()
	directory, repository := testutils.CreateRepository(t)
	defer os.RemoveAll(directory)
	worktree := testutils.GetRepositoryWorktree(t, repository)
	testutils.AddAndCommitNewFile(t, worktree, "some-file")
	testutils.AddAndCommitNewFile(t, worktree, "other-file")
	testutils.WriteFileInWorktree(t, worktree, "new-file")

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../data-sources/git_status/every_file",
		Vars: map[string]interface{}{
			"directory": directory,
		},
		NoColor: true,
	})

	defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApplyAndIdempotent(t, terraformOptions)

	everyFile := terraform.OutputMap(t, terraformOptions, "every_file")

	assert.NotNil(t, everyFile["some-file"], "some-file")
	assert.NotNil(t, everyFile["other-file"], "other-file")
	assert.NotNil(t, everyFile["new-file"], "new-file")
}
