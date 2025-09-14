/*
 * SPDX-FileCopyrightText: The terraform-provider-git Authors
 * SPDX-License-Identifier: 0BSD
 */

package provider_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	internal "github.com/metio/terraform-provider-git/internal/provider"
	"github.com/stretchr/testify/assert"
)

func TestGitProvider_Metadata(t *testing.T) {
	t.Parallel()
	p := &internal.GitProvider{}
	request := provider.MetadataRequest{}
	response := &provider.MetadataResponse{}
	p.Metadata(context.TODO(), request, response)

	assert.Equal(t, "git", response.TypeName, "TypeName")
}
