/*
 * SPDX-FileCopyrightText: The terraform-provider-git Authors
 * SPDX-License-Identifier: 0BSD
 */

package provider

import (
	"context"
	"os"
	"runtime"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func CreateCloneOptions(ctx context.Context, inputs *cloneResourceModel, diag *diag.Diagnostics) *git.CloneOptions {
	options := &git.CloneOptions{}

	if runtime.GOOS == "windows" {
		options.URL = strings.ReplaceAll(inputs.URL.ValueString(), "/", `\`)
	} else {
		options.URL = inputs.URL.ValueString()
	}

	tflog.Trace(ctx, "using 'URL'", map[string]any{
		"URL": inputs.URL.ValueString(),
	})

	options.RemoteName = inputs.RemoteName.ValueString()
	tflog.Trace(ctx, "using 'RemoteName'", map[string]any{
		"RemoteName": inputs.RemoteName.ValueString(),
	})

	options.InsecureSkipTLS = inputs.InsecureSkipTls.ValueBool()
	tflog.Trace(ctx, "using 'InsecureSkipTls'", map[string]any{
		"InsecureSkipTls": inputs.InsecureSkipTls.ValueBool(),
	})

	if len(inputs.CaBundleFilePath.ValueString()) > 0 {
		caBundle, err := os.ReadFile(inputs.CaBundleFilePath.ValueString())
		if err != nil {
			diag.AddError(
				"Invalid CA bundle file path",
				err.Error(),
			)
			return nil
		}
		options.CABundle = caBundle
		tflog.Trace(ctx, "using 'CaBundleFilePath'", map[string]any{
			"CaBundleFilePath": inputs.CaBundleFilePath.ValueString(),
		})
	}

	if !inputs.ReferenceName.IsNull() {
		options.ReferenceName = plumbing.ReferenceName(inputs.ReferenceName.ValueString())
		tflog.Trace(ctx, "using 'ReferenceName'", map[string]any{
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
