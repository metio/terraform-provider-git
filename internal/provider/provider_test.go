/*
 * SPDX-FileCopyrightText: The terraform-provider-git Authors
 * SPDX-License-Identifier: 0BSD
 */

package provider_test

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	internal "github.com/metio/terraform-provider-git/internal/provider"
	"github.com/metio/terraform-provider-git/internal/testutils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGitProvider_Metadata(t *testing.T) {
	t.Parallel()
	p := &internal.GitProvider{}
	request := provider.MetadataRequest{}
	response := &provider.MetadataResponse{}
	p.Metadata(context.TODO(), request, response)

	assert.Equal(t, "git", response.TypeName, "TypeName")
}

func TestGitProvider_GetSchema(t *testing.T) {
	t.Parallel()
	p := &internal.GitProvider{}
	schema, _ := p.GetSchema(context.TODO())

	testutils.VerifySchemaDescriptions(t, schema)
	assert.Nil(t, schema.Attributes, "should require no configuration")
}
