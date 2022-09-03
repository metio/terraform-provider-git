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

func TestDataSourceGitTags_AllTags(t *testing.T) {
	t.Parallel()
	directory, repository := testutils.CreateRepository(t)
	defer os.RemoveAll(directory)
	worktree := testutils.GetRepositoryWorktree(t, repository)
	testutils.AddAndCommitNewFile(t, worktree, "some-file")
	testutils.CreateTag(t, repository, "some-tag")
	testutils.CreateTag(t, repository, "other-tag")

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../data-sources/git_tags/all_tags",
		Vars: map[string]interface{}{
			"directory":   directory,
			"annotated":   true,
			"lightweight": true,
		},
		NoColor: true,
	})

	defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApplyAndIdempotent(t, terraformOptions)

	actualDirectory := terraform.Output(t, terraformOptions, "directory")
	actualId := terraform.Output(t, terraformOptions, "id")
	actualAnnotated := terraform.Output(t, terraformOptions, "annotated")
	actualLightweight := terraform.Output(t, terraformOptions, "lightweight")
	actualTags := terraform.OutputMap(t, terraformOptions, "tags")

	assert.Equal(t, directory, actualDirectory, "actualDirectory")
	assert.Equal(t, actualDirectory, actualId, "actualId")
	assert.Equal(t, "true", actualAnnotated, "actualAnnotated")
	assert.Equal(t, "true", actualLightweight, "actualLightweight")
	assert.NotNil(t, actualTags["some-tag"], "some-tag")
	assert.NotNil(t, actualTags["other-tag"], "other-tag")
}
