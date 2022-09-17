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

func TestResourceGitCommit_WithAdd(t *testing.T) {
	directory, repository := testutils.CreateRepository(t)
	defer os.RemoveAll(directory)
	worktree := testutils.GetRepositoryWorktree(t, repository)
	filename := "some-file"
	testutils.WriteFileInWorktree(t, worktree, filename)

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../resources/git_commit/with_add",
		Vars: map[string]interface{}{
			"directory": directory,
			"add_paths": []string{filename},
			"file":      filename,
		},
		NoColor: true,
	})

	defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApplyAndIdempotent(t, terraformOptions)

	actualStaging := terraform.Output(t, terraformOptions, "staging")
	actualWorktree := terraform.Output(t, terraformOptions, "worktree")
	actualSha1 := terraform.Output(t, terraformOptions, "sha1")
	actualFiles := terraform.OutputList(t, terraformOptions, "files")

	assert.Equal(t, "?", actualStaging, "actualStaging")
	assert.Equal(t, "?", actualWorktree, "actualWorktree")
	assert.NotNil(t, actualSha1, "actualSha1")
	assert.Equal(t, 1, len(actualFiles), "actualFiles")
}

func TestResourceGitCommit_WithAdd_Repeated(t *testing.T) {
	directory, repository := testutils.CreateRepository(t)
	defer os.RemoveAll(directory)
	worktree := testutils.GetRepositoryWorktree(t, repository)
	filename := "some-file"
	testutils.WriteFileInWorktree(t, worktree, filename)

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../resources/git_commit/with_add",
		Vars: map[string]interface{}{
			"directory": directory,
			"add_paths": []string{filename},
			"file":      filename,
		},
		NoColor: true,
	})

	defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApplyAndIdempotent(t, terraformOptions)

	actualStaging := terraform.Output(t, terraformOptions, "staging")
	actualWorktree := terraform.Output(t, terraformOptions, "worktree")
	actualSha1 := terraform.Output(t, terraformOptions, "sha1")
	actualFiles := terraform.OutputList(t, terraformOptions, "files")

	assert.Equal(t, "?", actualStaging, "actualStaging")
	assert.Equal(t, "?", actualWorktree, "actualWorktree")
	assert.NotNil(t, actualSha1, "actualSha1")
	assert.Equal(t, 1, len(actualFiles), "actualFiles")

	testutils.WriteFileContent(t, testutils.FileInWorktree(worktree, filename), "new content")

	terraform.ApplyAndIdempotent(t, terraformOptions)

	actualStaging = terraform.Output(t, terraformOptions, "staging")
	actualWorktree = terraform.Output(t, terraformOptions, "worktree")
	actualSha1 = terraform.Output(t, terraformOptions, "sha1")
	actualFiles = terraform.OutputList(t, terraformOptions, "files")

	assert.Equal(t, "?", actualStaging, "actualStaging")
	assert.Equal(t, "?", actualWorktree, "actualWorktree")
	assert.NotNil(t, actualSha1, "actualSha1")
	assert.Equal(t, 1, len(actualFiles), "actualFiles")
}

func TestResourceGitCommit_All(t *testing.T) {
	directory, repository := testutils.CreateRepository(t)
	defer os.RemoveAll(directory)
	worktree := testutils.GetRepositoryWorktree(t, repository)
	filename := "some-file"
	testutils.WriteFileInWorktree(t, worktree, filename)

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../resources/git_commit/all",
		Vars: map[string]interface{}{
			"directory": directory,
			"file":      filename,
		},
		NoColor: true,
	})

	defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApplyAndIdempotent(t, terraformOptions)

	actualStaging := terraform.Output(t, terraformOptions, "staging")
	actualWorktree := terraform.Output(t, terraformOptions, "worktree")
	actualSha1 := terraform.Output(t, terraformOptions, "sha1")
	actualFiles := terraform.OutputList(t, terraformOptions, "files")

	assert.Equal(t, "?", actualStaging, "actualStaging")
	assert.Equal(t, "?", actualWorktree, "actualWorktree")
	assert.NotNil(t, actualSha1, "actualSha1")
	assert.Equal(t, 0, len(actualFiles), "actualFiles")
}

func TestResourceGitCommit_All_ModifiedCommitted(t *testing.T) {
	directory, repository := testutils.CreateRepository(t)
	defer os.RemoveAll(directory)
	worktree := testutils.GetRepositoryWorktree(t, repository)
	filename := "some-file"
	testutils.AddAndCommitNewFile(t, worktree, filename)
	testutils.WriteFileContent(t, testutils.FileInWorktree(worktree, filename), "new content")

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../resources/git_commit/all",
		Vars: map[string]interface{}{
			"directory": directory,
			"file":      filename,
		},
		NoColor: true,
	})

	defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApplyAndIdempotent(t, terraformOptions)

	actualStaging := terraform.Output(t, terraformOptions, "staging")
	actualWorktree := terraform.Output(t, terraformOptions, "worktree")
	actualSha1 := terraform.Output(t, terraformOptions, "sha1")
	actualFiles := terraform.OutputList(t, terraformOptions, "files")

	assert.Equal(t, "?", actualStaging, "actualStaging")
	assert.Equal(t, "?", actualWorktree, "actualWorktree")
	assert.NotNil(t, actualSha1, "actualSha1")
	assert.Equal(t, 1, len(actualFiles), "actualFiles")
}

func TestResourceGitCommit_All_Repeated(t *testing.T) {
	directory, repository := testutils.CreateRepository(t)
	defer os.RemoveAll(directory)
	worktree := testutils.GetRepositoryWorktree(t, repository)
	filename := "some-file"
	testutils.WriteFileInWorktree(t, worktree, filename)

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../resources/git_commit/all",
		Vars: map[string]interface{}{
			"directory": directory,
			"file":      filename,
		},
		NoColor: true,
	})

	defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApplyAndIdempotent(t, terraformOptions)

	actualStaging := terraform.Output(t, terraformOptions, "staging")
	actualWorktree := terraform.Output(t, terraformOptions, "worktree")
	actualSha1 := terraform.Output(t, terraformOptions, "sha1")
	actualFiles := terraform.OutputList(t, terraformOptions, "files")

	assert.Equal(t, "?", actualStaging, "actualStaging")
	assert.Equal(t, "?", actualWorktree, "actualWorktree")
	assert.NotNil(t, actualSha1, "actualSha1")
	assert.Equal(t, 0, len(actualFiles), "actualFiles")

	testutils.WriteFileContent(t, testutils.FileInWorktree(worktree, filename), "new content")

	terraform.ApplyAndIdempotent(t, terraformOptions)

	actualStaging = terraform.Output(t, terraformOptions, "staging")
	actualWorktree = terraform.Output(t, terraformOptions, "worktree")
	actualSha1 = terraform.Output(t, terraformOptions, "sha1")
	actualFiles = terraform.OutputList(t, terraformOptions, "files")

	assert.Equal(t, "?", actualStaging, "actualStaging")
	assert.Equal(t, "?", actualWorktree, "actualWorktree")
	assert.NotNil(t, actualSha1, "actualSha1")
	assert.Equal(t, 0, len(actualFiles), "actualFiles")
}

func TestResourceGitCommit_All_ModifiedCommitted_Repeated(t *testing.T) {
	directory, repository := testutils.CreateRepository(t)
	defer os.RemoveAll(directory)
	worktree := testutils.GetRepositoryWorktree(t, repository)
	filename := "some-file"
	testutils.AddAndCommitNewFile(t, worktree, filename)
	testutils.WriteFileContent(t, testutils.FileInWorktree(worktree, filename), "new content")

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../resources/git_commit/all",
		Vars: map[string]interface{}{
			"directory": directory,
			"file":      filename,
		},
		NoColor: true,
	})

	defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApplyAndIdempotent(t, terraformOptions)

	actualStaging := terraform.Output(t, terraformOptions, "staging")
	actualWorktree := terraform.Output(t, terraformOptions, "worktree")
	actualSha1 := terraform.Output(t, terraformOptions, "sha1")
	actualFiles := terraform.OutputList(t, terraformOptions, "files")

	assert.Equal(t, "?", actualStaging, "actualStaging")
	assert.Equal(t, "?", actualWorktree, "actualWorktree")
	assert.NotNil(t, actualSha1, "actualSha1")
	assert.Equal(t, 1, len(actualFiles), "actualFiles")

	testutils.WriteFileContent(t, testutils.FileInWorktree(worktree, filename), "repeated content")

	terraform.ApplyAndIdempotent(t, terraformOptions)

	actualStaging = terraform.Output(t, terraformOptions, "staging")
	actualWorktree = terraform.Output(t, terraformOptions, "worktree")
	actualSha1 = terraform.Output(t, terraformOptions, "sha1")
	actualFiles = terraform.OutputList(t, terraformOptions, "files")

	assert.Equal(t, "?", actualStaging, "actualStaging")
	assert.Equal(t, "?", actualWorktree, "actualWorktree")
	assert.NotNil(t, actualSha1, "actualSha1")
	assert.Equal(t, 1, len(actualFiles), "actualFiles")
}

func TestResourceGitCommit_WithoutAdd(t *testing.T) {
	directory, repository := testutils.CreateRepository(t)
	defer os.RemoveAll(directory)
	worktree := testutils.GetRepositoryWorktree(t, repository)
	filename := "some-file"
	testutils.WriteFileInWorktree(t, worktree, filename)
	testutils.GitAdd(t, worktree, filename)

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../resources/git_commit/without_add",
		Vars: map[string]interface{}{
			"directory": directory,
			"file":      filename,
		},
		NoColor: true,
	})

	defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApplyAndIdempotent(t, terraformOptions)

	actualStaging := terraform.Output(t, terraformOptions, "staging")
	actualWorktree := terraform.Output(t, terraformOptions, "worktree")
	actualSha1 := terraform.Output(t, terraformOptions, "sha1")
	actualFiles := terraform.OutputList(t, terraformOptions, "files")

	assert.Equal(t, "?", actualStaging, "actualStaging")
	assert.Equal(t, "?", actualWorktree, "actualWorktree")
	assert.NotNil(t, actualSha1, "actualSha1")
	assert.Equal(t, 1, len(actualFiles), "actualFiles")
}

func TestResourceGitCommit_WithoutAdd_Repeated(t *testing.T) {
	directory, repository := testutils.CreateRepository(t)
	defer os.RemoveAll(directory)
	worktree := testutils.GetRepositoryWorktree(t, repository)
	filename := "some-file"
	testutils.WriteFileInWorktree(t, worktree, filename)
	testutils.GitAdd(t, worktree, filename)

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../resources/git_commit/without_add",
		Vars: map[string]interface{}{
			"directory": directory,
			"file":      filename,
		},
		NoColor: true,
	})

	defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApplyAndIdempotent(t, terraformOptions)

	actualStaging := terraform.Output(t, terraformOptions, "staging")
	actualWorktree := terraform.Output(t, terraformOptions, "worktree")
	actualSha1 := terraform.Output(t, terraformOptions, "sha1")
	actualFiles := terraform.OutputList(t, terraformOptions, "files")

	assert.Equal(t, "?", actualStaging, "actualStaging")
	assert.Equal(t, "?", actualWorktree, "actualWorktree")
	assert.NotNil(t, actualSha1, "actualSha1")
	assert.Equal(t, 1, len(actualFiles), "actualFiles")

	testutils.WriteFileContent(t, testutils.FileInWorktree(worktree, filename), "new content")
	testutils.GitAdd(t, worktree, filename)

	terraform.ApplyAndIdempotent(t, terraformOptions)

	actualStaging = terraform.Output(t, terraformOptions, "staging")
	actualWorktree = terraform.Output(t, terraformOptions, "worktree")
	actualSha1 = terraform.Output(t, terraformOptions, "sha1")
	actualFiles = terraform.OutputList(t, terraformOptions, "files")

	assert.Equal(t, "?", actualStaging, "actualStaging")
	assert.Equal(t, "?", actualWorktree, "actualWorktree")
	assert.NotNil(t, actualSha1, "actualSha1")
	assert.Equal(t, 1, len(actualFiles), "actualFiles")
}
