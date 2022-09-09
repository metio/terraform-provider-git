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

func TestResourceGitAdd_SingleFile_Exact_WriteOnce(t *testing.T) {
	directory, repository := testutils.CreateRepository(t)
	defer os.RemoveAll(directory)
	worktree := testutils.GetRepositoryWorktree(t, repository)
	filename := "some-file"
	testutils.WriteFileInWorktree(t, worktree, filename)

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../resources/git_add/single_file",
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

	assert.Equal(t, "A", actualStaging, "actualStaging")
	assert.Equal(t, " ", actualWorktree, "actualWorktree")
}

func TestResourceGitAdd_SingleFile_Glob_WriteOnce(t *testing.T) {
	directory, repository := testutils.CreateRepository(t)
	defer os.RemoveAll(directory)
	worktree := testutils.GetRepositoryWorktree(t, repository)
	filename := "some-file"
	testutils.WriteFileInWorktree(t, worktree, filename)

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../resources/git_add/single_file",
		Vars: map[string]interface{}{
			"directory": directory,
			"add_paths": []string{"some*"},
			"file":      filename,
		},
		NoColor: true,
	})

	defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApplyAndIdempotent(t, terraformOptions)

	actualStaging := terraform.Output(t, terraformOptions, "staging")
	actualWorktree := terraform.Output(t, terraformOptions, "worktree")

	assert.Equal(t, "A", actualStaging, "actualStaging")
	assert.Equal(t, " ", actualWorktree, "actualWorktree")
}

func TestResourceGitAdd_SingleFile_Exact_NoMatch(t *testing.T) {
	directory, repository := testutils.CreateRepository(t)
	defer os.RemoveAll(directory)
	worktree := testutils.GetRepositoryWorktree(t, repository)
	filename := "some-file"
	testutils.WriteFileInWorktree(t, worktree, filename)

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../resources/git_add/single_file",
		Vars: map[string]interface{}{
			"directory": directory,
			"add_paths": []string{"other-file"},
			"file":      filename,
		},
		NoColor: true,
	})

	defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApplyAndIdempotent(t, terraformOptions)

	actualStaging := terraform.Output(t, terraformOptions, "staging")
	actualWorktree := terraform.Output(t, terraformOptions, "worktree")

	assert.Equal(t, "?", actualStaging, "actualStaging")
	assert.Equal(t, "?", actualWorktree, "actualWorktree")
}

func TestResourceGitAdd_SingleFile_Glob_NoMatch(t *testing.T) {
	directory, repository := testutils.CreateRepository(t)
	defer os.RemoveAll(directory)
	worktree := testutils.GetRepositoryWorktree(t, repository)
	filename := "some-file"
	testutils.WriteFileInWorktree(t, worktree, filename)

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../resources/git_add/single_file",
		Vars: map[string]interface{}{
			"directory": directory,
			"add_paths": []string{"other*"},
			"file":      filename,
		},
		NoColor: true,
	})

	defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApplyAndIdempotent(t, terraformOptions)

	actualStaging := terraform.Output(t, terraformOptions, "staging")
	actualWorktree := terraform.Output(t, terraformOptions, "worktree")

	assert.Equal(t, "?", actualStaging, "actualStaging")
	assert.Equal(t, "?", actualWorktree, "actualWorktree")
}

func TestResourceGitAdd_SingleFile_Exact_WriteMultiple(t *testing.T) {
	directory, repository := testutils.CreateRepository(t)
	defer os.RemoveAll(directory)
	worktree := testutils.GetRepositoryWorktree(t, repository)
	filename := "some-file"
	testutils.WriteFileInWorktree(t, worktree, filename)

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../resources/git_add/single_file",
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

	assert.Equal(t, "A", actualStaging, "actualStaging")
	assert.Equal(t, " ", actualWorktree, "actualWorktree")

	// Modify file in order to trigger new git-add call
	testutils.WriteFileContent(t, testutils.FileInWorktree(worktree, filename), "new content")
	terraform.ApplyAndIdempotent(t, terraformOptions)

	actualStaging2 := terraform.Output(t, terraformOptions, "staging")
	actualWorktree2 := terraform.Output(t, terraformOptions, "worktree")

	assert.Equal(t, "A", actualStaging2, "actualStaging2")
	assert.Equal(t, " ", actualWorktree2, "actualWorktree2")
}

func TestResourceGitAdd_SingleFile_Glob_WriteMultiple(t *testing.T) {
	directory, repository := testutils.CreateRepository(t)
	defer os.RemoveAll(directory)
	worktree := testutils.GetRepositoryWorktree(t, repository)
	filename := "some-file"
	testutils.WriteFileInWorktree(t, worktree, filename)

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../resources/git_add/single_file",
		Vars: map[string]interface{}{
			"directory": directory,
			"add_paths": []string{"some*"},
			"file":      filename,
		},
		NoColor: true,
	})

	defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApplyAndIdempotent(t, terraformOptions)

	actualStaging := terraform.Output(t, terraformOptions, "staging")
	actualWorktree := terraform.Output(t, terraformOptions, "worktree")

	assert.Equal(t, "A", actualStaging, "actualStaging")
	assert.Equal(t, " ", actualWorktree, "actualWorktree")

	// Modify file in order to trigger new git-add call
	testutils.WriteFileContent(t, testutils.FileInWorktree(worktree, filename), "new content")
	terraform.ApplyAndIdempotent(t, terraformOptions)

	actualStaging2 := terraform.Output(t, terraformOptions, "staging")
	actualWorktree2 := terraform.Output(t, terraformOptions, "worktree")

	assert.Equal(t, "A", actualStaging2, "actualStaging2")
	assert.Equal(t, " ", actualWorktree2, "actualWorktree2")
}

func TestResourceGitAdd_SingleFile_Exact_WriteDelete(t *testing.T) {
	directory, repository := testutils.CreateRepository(t)
	defer os.RemoveAll(directory)
	worktree := testutils.GetRepositoryWorktree(t, repository)
	filename := "some-file"
	testutils.WriteFileInWorktree(t, worktree, filename)

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../resources/git_add/single_file",
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

	assert.Equal(t, "A", actualStaging, "actualStaging")
	assert.Equal(t, " ", actualWorktree, "actualWorktree")

	os.Remove(testutils.FileInWorktree(worktree, filename))
	terraform.ApplyAndIdempotent(t, terraformOptions)

	actualStaging2 := terraform.Output(t, terraformOptions, "staging")
	actualWorktree2 := terraform.Output(t, terraformOptions, "worktree")

	assert.Equal(t, "?", actualStaging2, "actualStaging2")
	assert.Equal(t, "?", actualWorktree2, "actualWorktree2")
}

func TestResourceGitAdd_SingleFile_Glob_WriteDelete(t *testing.T) {
	directory, repository := testutils.CreateRepository(t)
	defer os.RemoveAll(directory)
	worktree := testutils.GetRepositoryWorktree(t, repository)
	filename := "some-file"
	testutils.WriteFileInWorktree(t, worktree, filename)

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../resources/git_add/single_file",
		Vars: map[string]interface{}{
			"directory": directory,
			"add_paths": []string{"some*"},
			"file":      filename,
		},
		NoColor: true,
	})

	defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApplyAndIdempotent(t, terraformOptions)

	actualStaging := terraform.Output(t, terraformOptions, "staging")
	actualWorktree := terraform.Output(t, terraformOptions, "worktree")

	assert.Equal(t, "A", actualStaging, "actualStaging")
	assert.Equal(t, " ", actualWorktree, "actualWorktree")

	os.Remove(testutils.FileInWorktree(worktree, filename))
	terraform.ApplyAndIdempotent(t, terraformOptions)

	actualStaging2 := terraform.Output(t, terraformOptions, "staging")
	actualWorktree2 := terraform.Output(t, terraformOptions, "worktree")

	assert.Equal(t, "?", actualStaging2, "actualStaging2")
	assert.Equal(t, "?", actualWorktree2, "actualWorktree2")
}

func TestResourceGitAdd_SingleFile_Exact_DeleteCommitted(t *testing.T) {
	directory, repository := testutils.CreateRepository(t)
	defer os.RemoveAll(directory)
	worktree := testutils.GetRepositoryWorktree(t, repository)
	filename := "some-file"
	testutils.AddAndCommitNewFile(t, worktree, filename)

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../resources/git_add/single_file",
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

	assert.Equal(t, "?", actualStaging, "actualStaging")
	assert.Equal(t, "?", actualWorktree, "actualWorktree")

	os.Remove(testutils.FileInWorktree(worktree, filename))
	terraform.ApplyAndIdempotent(t, terraformOptions)

	actualStaging2 := terraform.Output(t, terraformOptions, "staging")
	actualWorktree2 := terraform.Output(t, terraformOptions, "worktree")

	assert.Equal(t, "D", actualStaging2, "actualStaging2")
	assert.Equal(t, " ", actualWorktree2, "actualWorktree2")
}

func TestResourceGitAdd_SingleFile_Glob_DeleteCommitted(t *testing.T) {
	directory, repository := testutils.CreateRepository(t)
	defer os.RemoveAll(directory)
	worktree := testutils.GetRepositoryWorktree(t, repository)
	filename := "some-file"
	testutils.AddAndCommitNewFile(t, worktree, filename)

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../resources/git_add/single_file",
		Vars: map[string]interface{}{
			"directory": directory,
			"add_paths": []string{"some*"},
			"file":      filename,
		},
		NoColor: true,
	})

	defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApplyAndIdempotent(t, terraformOptions)

	actualStaging := terraform.Output(t, terraformOptions, "staging")
	actualWorktree := terraform.Output(t, terraformOptions, "worktree")

	assert.Equal(t, "?", actualStaging, "actualStaging")
	assert.Equal(t, "?", actualWorktree, "actualWorktree")

	os.Remove(testutils.FileInWorktree(worktree, filename))
	terraform.ApplyAndIdempotent(t, terraformOptions)

	actualStaging2 := terraform.Output(t, terraformOptions, "staging")
	actualWorktree2 := terraform.Output(t, terraformOptions, "worktree")

	assert.Equal(t, "D", actualStaging2, "actualStaging2")
	assert.Equal(t, " ", actualWorktree2, "actualWorktree2")
}

func TestResourceGitAdd_MultipleFiles_WriteOnce(t *testing.T) {
	directory, repository := testutils.CreateRepository(t)
	defer os.RemoveAll(directory)
	worktree := testutils.GetRepositoryWorktree(t, repository)
	filename := "some-file"
	testutils.WriteFileInWorktree(t, worktree, filename)
	testutils.WriteFileInWorktree(t, worktree, "other-file")

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../resources/git_add/multiple_files",
		Vars: map[string]interface{}{
			"directory": directory,
			"add_paths": []string{filename},
		},
		NoColor: true,
	})

	defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApplyAndIdempotent(t, terraformOptions)

	actualFiles := terraform.OutputMapOfObjects(t, terraformOptions, "files")

	assert.NotNil(t, actualFiles[filename], filename)
	assert.NotNil(t, actualFiles["other-file"], filename)

	stagedFile := actualFiles[filename].(map[string]interface{})
	unStagedFile := actualFiles["other-file"].(map[string]interface{})

	assert.Equal(t, "A", stagedFile["staging"], "stagedFile-staging")
	assert.Equal(t, " ", stagedFile["worktree"], "stagedFile-worktree")
	assert.Equal(t, "?", unStagedFile["staging"], "unStagedFile-staging")
	assert.Equal(t, "?", unStagedFile["worktree"], "unStagedFile-worktree")
}
