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

func TestDataSourceGitRemotes_AllRemotes(t *testing.T) {
	t.Parallel()
	directory, repository := testutils.CreateRepository(t)
	testutils.CreateRemote(t, repository, "origin")
	testutils.CreateRemote(t, repository, "upstream")

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../data-sources/git_remotes/all_remotes",
		Vars: map[string]any{
			"directory": directory,
		},
		NoColor: true,
	})

	defer terraform.DestroyContext(t, t.Context(), terraformOptions)
	terraform.InitAndApplyAndIdempotentContext(t, t.Context(), terraformOptions)

	actualDirectory := terraform.OutputContext(t, t.Context(), terraformOptions, "directory")
	actualId := terraform.OutputContext(t, t.Context(), terraformOptions, "id")
	actualRemotes := terraform.OutputMapContext(t, t.Context(), terraformOptions, "remotes")

	assert.Equal(t, directory, actualDirectory, "actualDirectory")
	assert.Equal(t, directory, actualId, "actualId")
	assert.NotNil(t, actualRemotes["origin"], "origin")
	assert.NotNil(t, actualRemotes["upstream"], "upstream")
}
