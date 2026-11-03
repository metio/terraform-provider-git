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

func TestDataSourceGitConfig_ScopedLocal(t *testing.T) {
	directory, _ := testutils.CreateRepository(t)
	scope := "local"

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../data-sources/git_config/scoped",
		Vars: map[string]interface{}{
			"directory": directory,
			"scope":     scope,
		},
		NoColor: true,
	})

	defer terraform.DestroyContext(t, t.Context(), terraformOptions)
	terraform.InitAndApplyAndIdempotentContext(t, t.Context(), terraformOptions)

	actualDirectory := terraform.OutputContext(t, t.Context(), terraformOptions, "directory")
	actualId := terraform.OutputContext(t, t.Context(), terraformOptions, "id")
	actualScope := terraform.OutputContext(t, t.Context(), terraformOptions, "scope")
	authorEmail := terraform.OutputContext(t, t.Context(), terraformOptions, "author_email")
	authorName := terraform.OutputContext(t, t.Context(), terraformOptions, "author_name")
	committerEmail := terraform.OutputContext(t, t.Context(), terraformOptions, "committer_email")
	committerName := terraform.OutputContext(t, t.Context(), terraformOptions, "committer_name")
	userEmail := terraform.OutputContext(t, t.Context(), terraformOptions, "user_email")
	userName := terraform.OutputContext(t, t.Context(), terraformOptions, "user_name")

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
	scope := "global"

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../data-sources/git_config/scoped",
		Vars: map[string]interface{}{
			"directory": directory,
			"scope":     scope,
		},
		NoColor: true,
	})

	defer terraform.DestroyContext(t, t.Context(), terraformOptions)
	terraform.InitAndApplyAndIdempotentContext(t, t.Context(), terraformOptions)

	actualDirectory := terraform.OutputContext(t, t.Context(), terraformOptions, "directory")
	actualId := terraform.OutputContext(t, t.Context(), terraformOptions, "id")
	actualScope := terraform.OutputContext(t, t.Context(), terraformOptions, "scope")
	authorEmail := terraform.OutputContext(t, t.Context(), terraformOptions, "author_email")
	authorName := terraform.OutputContext(t, t.Context(), terraformOptions, "author_name")
	committerEmail := terraform.OutputContext(t, t.Context(), terraformOptions, "committer_email")
	committerName := terraform.OutputContext(t, t.Context(), terraformOptions, "committer_name")
	userEmail := terraform.OutputContext(t, t.Context(), terraformOptions, "user_email")
	userName := terraform.OutputContext(t, t.Context(), terraformOptions, "user_name")

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
	scope := "system"

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../data-sources/git_config/scoped",
		Vars: map[string]interface{}{
			"directory": directory,
			"scope":     scope,
		},
		NoColor: true,
	})

	defer terraform.DestroyContext(t, t.Context(), terraformOptions)
	terraform.InitAndApplyAndIdempotentContext(t, t.Context(), terraformOptions)

	actualDirectory := terraform.OutputContext(t, t.Context(), terraformOptions, "directory")
	actualId := terraform.OutputContext(t, t.Context(), terraformOptions, "id")
	actualScope := terraform.OutputContext(t, t.Context(), terraformOptions, "scope")
	authorEmail := terraform.OutputContext(t, t.Context(), terraformOptions, "author_email")
	authorName := terraform.OutputContext(t, t.Context(), terraformOptions, "author_name")
	committerEmail := terraform.OutputContext(t, t.Context(), terraformOptions, "committer_email")
	committerName := terraform.OutputContext(t, t.Context(), terraformOptions, "committer_name")
	userEmail := terraform.OutputContext(t, t.Context(), terraformOptions, "user_email")
	userName := terraform.OutputContext(t, t.Context(), terraformOptions, "user_name")

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

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../data-sources/git_config/using_defaults",
		Vars: map[string]interface{}{
			"directory": directory,
		},
		NoColor: true,
	})

	defer terraform.DestroyContext(t, t.Context(), terraformOptions)
	terraform.InitAndApplyAndIdempotentContext(t, t.Context(), terraformOptions)

	actualDirectory := terraform.OutputContext(t, t.Context(), terraformOptions, "directory")
	actualId := terraform.OutputContext(t, t.Context(), terraformOptions, "id")
	actualScope := terraform.OutputContext(t, t.Context(), terraformOptions, "scope")
	authorEmail := terraform.OutputContext(t, t.Context(), terraformOptions, "author_email")
	authorName := terraform.OutputContext(t, t.Context(), terraformOptions, "author_name")
	committerEmail := terraform.OutputContext(t, t.Context(), terraformOptions, "committer_email")
	committerName := terraform.OutputContext(t, t.Context(), terraformOptions, "committer_name")
	userEmail := terraform.OutputContext(t, t.Context(), terraformOptions, "user_email")
	userName := terraform.OutputContext(t, t.Context(), terraformOptions, "user_name")

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
