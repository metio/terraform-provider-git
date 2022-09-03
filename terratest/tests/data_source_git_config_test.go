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

func TestDataSourceGitConfig_ScopedLocal(t *testing.T) {
	directory, _ := testutils.CreateRepository(t)
	defer os.RemoveAll(directory)
	scope := "local"

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../data-sources/git_config/scoped",
		Vars: map[string]interface{}{
			"directory": directory,
			"scope":     scope,
		},
		NoColor: true,
	})

	defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApplyAndIdempotent(t, terraformOptions)

	actualDirectory := terraform.Output(t, terraformOptions, "directory")
	actualId := terraform.Output(t, terraformOptions, "id")
	actualScope := terraform.Output(t, terraformOptions, "scope")
	authorEmail := terraform.Output(t, terraformOptions, "author_email")
	authorName := terraform.Output(t, terraformOptions, "author_name")
	committerEmail := terraform.Output(t, terraformOptions, "committer_email")
	committerName := terraform.Output(t, terraformOptions, "committer_name")
	userEmail := terraform.Output(t, terraformOptions, "user_email")
	userName := terraform.Output(t, terraformOptions, "user_name")

	assert.Equal(t, directory, actualDirectory, "actualDirectory")
	assert.Equal(t, directory, actualId, "actualId")
	assert.Equal(t, scope, actualScope, "actualScope")
	assert.NotNil(t, authorEmail, "authorEmail")
	assert.NotNil(t, authorName, "authorName")
	assert.NotNil(t, committerEmail, "committerEmail")
	assert.NotNil(t, committerName, "committerName")
	assert.NotNil(t, userEmail, "userEmail")
	assert.NotNil(t, userName, "userName")
}

func TestDataSourceGitConfig_ScopedGlobal(t *testing.T) {
	directory, _ := testutils.CreateRepository(t)
	defer os.RemoveAll(directory)
	scope := "global"

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../data-sources/git_config/scoped",
		Vars: map[string]interface{}{
			"directory": directory,
			"scope":     scope,
		},
		NoColor: true,
	})

	defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApplyAndIdempotent(t, terraformOptions)

	actualDirectory := terraform.Output(t, terraformOptions, "directory")
	actualId := terraform.Output(t, terraformOptions, "id")
	actualScope := terraform.Output(t, terraformOptions, "scope")
	authorEmail := terraform.Output(t, terraformOptions, "author_email")
	authorName := terraform.Output(t, terraformOptions, "author_name")
	committerEmail := terraform.Output(t, terraformOptions, "committer_email")
	committerName := terraform.Output(t, terraformOptions, "committer_name")
	userEmail := terraform.Output(t, terraformOptions, "user_email")
	userName := terraform.Output(t, terraformOptions, "user_name")

	assert.Equal(t, directory, actualDirectory, "actualDirectory")
	assert.Equal(t, directory, actualId, "actualId")
	assert.Equal(t, scope, actualScope, "actualScope")
	assert.NotNil(t, authorEmail, "authorEmail")
	assert.NotNil(t, authorName, "authorName")
	assert.NotNil(t, committerEmail, "committerEmail")
	assert.NotNil(t, committerName, "committerName")
	assert.NotNil(t, userEmail, "userEmail")
	assert.NotNil(t, userName, "userName")
}

func TestDataSourceGitConfig_ScopedSystem(t *testing.T) {
	directory, _ := testutils.CreateRepository(t)
	defer os.RemoveAll(directory)
	scope := "system"

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../data-sources/git_config/scoped",
		Vars: map[string]interface{}{
			"directory": directory,
			"scope":     scope,
		},
		NoColor: true,
	})

	defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApplyAndIdempotent(t, terraformOptions)

	actualDirectory := terraform.Output(t, terraformOptions, "directory")
	actualId := terraform.Output(t, terraformOptions, "id")
	actualScope := terraform.Output(t, terraformOptions, "scope")
	authorEmail := terraform.Output(t, terraformOptions, "author_email")
	authorName := terraform.Output(t, terraformOptions, "author_name")
	committerEmail := terraform.Output(t, terraformOptions, "committer_email")
	committerName := terraform.Output(t, terraformOptions, "committer_name")
	userEmail := terraform.Output(t, terraformOptions, "user_email")
	userName := terraform.Output(t, terraformOptions, "user_name")

	assert.Equal(t, directory, actualDirectory, "actualDirectory")
	assert.Equal(t, directory, actualId, "actualId")
	assert.Equal(t, scope, actualScope, "actualScope")
	assert.NotNil(t, authorEmail, "authorEmail")
	assert.NotNil(t, authorName, "authorName")
	assert.NotNil(t, committerEmail, "committerEmail")
	assert.NotNil(t, committerName, "committerName")
	assert.NotNil(t, userEmail, "userEmail")
	assert.NotNil(t, userName, "userName")
}

func TestDataSourceGitConfig_UsingDefaults(t *testing.T) {
	t.Parallel()
	directory, _ := testutils.CreateRepository(t)
	defer os.RemoveAll(directory)

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../data-sources/git_config/using_defaults",
		Vars: map[string]interface{}{
			"directory": directory,
		},
		NoColor: true,
	})

	defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApplyAndIdempotent(t, terraformOptions)

	actualDirectory := terraform.Output(t, terraformOptions, "directory")
	actualId := terraform.Output(t, terraformOptions, "id")
	actualScope := terraform.Output(t, terraformOptions, "scope")
	authorEmail := terraform.Output(t, terraformOptions, "author_email")
	authorName := terraform.Output(t, terraformOptions, "author_name")
	committerEmail := terraform.Output(t, terraformOptions, "committer_email")
	committerName := terraform.Output(t, terraformOptions, "committer_name")
	userEmail := terraform.Output(t, terraformOptions, "user_email")
	userName := terraform.Output(t, terraformOptions, "user_name")

	assert.Equal(t, directory, actualDirectory, "actualDirectory")
	assert.Equal(t, directory, actualId, "actualId")
	assert.Equal(t, "global", actualScope, "actualScope")
	assert.NotNil(t, authorEmail, "authorEmail")
	assert.NotNil(t, authorName, "authorName")
	assert.NotNil(t, committerEmail, "committerEmail")
	assert.NotNil(t, committerName, "committerName")
	assert.NotNil(t, userEmail, "userEmail")
	assert.NotNil(t, userName, "userName")
}
