/*
 * SPDX-FileCopyrightText: The terraform-provider-git Authors
 * SPDX-License-Identifier: 0BSD
 */

package provider_test

import (
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/metio/terraform-provider-git/internal/testutils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDataSourceGitRemote_SingleRemote(t *testing.T) {
	t.Parallel()
	directory, repository := testutils.CreateRepository(t)
	name := "origin"
	testutils.CreateRemote(t, repository, name)

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../data-sources/git_remote/single_remote",
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
	actualUrls := terraform.OutputList(t, terraformOptions, "urls")

	assert.Equal(t, directory, actualDirectory, "actualDirectory")
	assert.Equal(t, name, actualId, "actualId")
	assert.Equal(t, name, actualName, "actualName")
	assert.Equal(t, 1, len(actualUrls), "actualUrls")
}
