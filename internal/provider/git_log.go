/*
 * SPDX-FileCopyrightText: The terraform-provider-git Authors
 * SPDX-License-Identifier: 0BSD
 */

package provider

import (
	"context"
	"path"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func createLogOptions(ctx context.Context, repository *git.Repository, inputs *logDataSourceModel, diag *diag.Diagnostics) *git.LogOptions {
	logOptions := &git.LogOptions{}

	logOptions.All = inputs.All.ValueBool()
	tflog.Trace(ctx, "using 'All'", map[string]interface{}{
		"all": logOptions.All,
	})

	if !inputs.From.IsNull() && !inputs.From.IsUnknown() {
		hash := resolveRevision(ctx, repository, inputs.From.ValueString(), diag)
		if hash == nil {
			return nil
		}
		logOptions.From = *hash
		tflog.Trace(ctx, "using 'From'", map[string]interface{}{
			"from": logOptions.From,
		})
	}

	if !inputs.Since.IsNull() && !inputs.Since.IsUnknown() {
		since, err := time.Parse(time.RFC3339, inputs.Since.ValueString())
		if err != nil {
			diag.AddError(
				"Cannot parse given time",
				"Could not parse 'since' with value ["+inputs.Since.ValueString()+"] because of: "+err.Error(),
			)
			return nil
		}
		logOptions.Since = &since
		tflog.Trace(ctx, "using 'Since'", map[string]interface{}{
			"since": logOptions.Since,
		})
	}

	if !inputs.Until.IsNull() && !inputs.Until.IsUnknown() {
		until, err := time.Parse(time.RFC3339, inputs.Until.ValueString())
		if err != nil {
			diag.AddError(
				"Cannot parse given time",
				"Could not parse 'until' with value ["+inputs.Until.ValueString()+"] because of: "+err.Error(),
			)
			return nil
		}
		logOptions.Until = &until
		tflog.Trace(ctx, "using 'Until'", map[string]interface{}{
			"until": logOptions.Until,
		})
	}

	if !inputs.Order.IsNull() && !inputs.Order.IsUnknown() {
		switch inputs.Order.ValueString() {
		case "depth":
			logOptions.Order = git.LogOrderDFS
		case "breadth":
			logOptions.Order = git.LogOrderBSF
		default:
			logOptions.Order = git.LogOrderCommitterTime
		}
		tflog.Trace(ctx, "using 'Order'", map[string]interface{}{
			"order": logOptions.Order,
		})
	}

	if !inputs.FilterPaths.IsNull() && !inputs.FilterPaths.IsUnknown() {
		paths := make([]string, len(inputs.FilterPaths.Elements()))
		diag.Append(inputs.FilterPaths.ElementsAs(ctx, &paths, false)...)
		if diag.HasError() {
			return nil
		}
		tflog.Trace(ctx, "using 'FilterPaths'", map[string]interface{}{
			"filter_paths": paths,
		})
		logOptions.PathFilter = func(file string) bool {
			for _, pattern := range paths {
				match, err := path.Match(pattern, file)
				if err != nil {
					return false
				}
				if match {
					return true
				}
			}
			return false
		}
	}

	return logOptions
}
