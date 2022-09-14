/*
 * SPDX-FileCopyrightText: The terraform-provider-git Authors
 * SPDX-License-Identifier: 0BSD
 */

package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGitProvider_Metadata(t *testing.T) {
	p := &gitProvider{}
	request := provider.MetadataRequest{}
	response := &provider.MetadataResponse{}
	p.Metadata(context.TODO(), request, response)

	assert.Equal(t, "git", response.TypeName)
}

func TestGitProvider_GetSchema(t *testing.T) {
	p := &gitProvider{}
	schema, _ := p.GetSchema(context.TODO())

	assert.NotNil(t, schema.Description)
	assert.NotNil(t, schema.MarkdownDescription)
	assert.Nil(t, schema.Attributes)
}
