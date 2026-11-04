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

func TestDataSourceGitRemote_SingleRemote(t *testing.T) {
	t.Parallel()
	directory, repository := testutils.CreateRepository(t)
	name := "origin"
	repositoryName := "some-repository-name"
	testutils.CreateRemoteWithName(t, repository, name, repositoryName)

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../data-sources/git_remote/single_remote",
		Vars: map[string]any{
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
	actualUrls := terraform.OutputListContext(t, t.Context(), terraformOptions, "urls")
	actualUpstreamRepositoryName := terraform.OutputContext(t, t.Context(), terraformOptions, "upstream_repository_name")

	assert.Equal(t, directory, actualDirectory, "actualDirectory")
	assert.Equal(t, name, actualId, "actualId")
	assert.Equal(t, name, actualName, "actualName")
	assert.Equal(t, 1, len(actualUrls), "actualUrls")
	assert.Equal(t, repositoryName, actualUpstreamRepositoryName, "actualUpstreamRepositoryName")
}
