/*
 * SPDX-FileCopyrightText: The terraform-provider-git Authors
 * SPDX-License-Identifier: 0BSD
 */

package provider

import (
	"context"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func CreateCloneOptions(ctx context.Context, inputs *cloneResourceModel, diag *diag.Diagnostics) *git.CloneOptions {
	options := &git.CloneOptions{}

	options.URL = strings.ReplaceAll(inputs.URL.ValueString(), "/", "\\")
	tflog.Trace(ctx, "using 'URL'", map[string]interface{}{
		"URL": inputs.URL.ValueString(),
	})

	options.RemoteName = inputs.RemoteName.ValueString()
	tflog.Trace(ctx, "using 'RemoteName'", map[string]interface{}{
		"RemoteName": inputs.RemoteName.ValueString(),
	})

	if !inputs.ReferenceName.IsNull() {
		options.ReferenceName = plumbing.ReferenceName(inputs.ReferenceName.ValueString())
		tflog.Trace(ctx, "using 'ReferenceName'", map[string]interface{}{
			"ReferenceName": inputs.ReferenceName.ValueString(),
		})
	}

	if !inputs.Auth.IsNull() {
		options.Auth = authOptions(ctx, inputs.Auth, diag)
		if diag.HasError() {
			return nil
		}
	}

	return options
}
