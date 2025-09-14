/*
 * SPDX-FileCopyrightText: The terraform-provider-git Authors
 * SPDX-License-Identifier: 0BSD
 */

package provider

import (
	"context"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func CreatePushOptions(ctx context.Context, inputs *PushResourceModel, diag *diag.Diagnostics) *git.PushOptions {
	options := &git.PushOptions{}

	if len(inputs.RefSpecs.Elements()) > 0 {
		refSpecs := make([]config.RefSpec, len(inputs.RefSpecs.Elements()))
		diag.Append(inputs.RefSpecs.ElementsAs(ctx, &refSpecs, false)...)
		if diag.HasError() {
			return nil
		}
		options.RefSpecs = refSpecs
		tflog.Trace(ctx, "using 'RefSpecs'", map[string]interface{}{
			"RefSpecs": refSpecs,
		})
	} else {
		return nil
	}

	options.RemoteName = inputs.Remote.ValueString()
	tflog.Trace(ctx, "using 'RemoteName'", map[string]interface{}{
		"RemoteName": inputs.Remote.ValueString(),
	})

	options.Prune = inputs.Prune.ValueBool()
	tflog.Trace(ctx, "using 'Prune'", map[string]interface{}{
		"Prune": inputs.Prune.ValueBool(),
	})

	options.Force = inputs.Force.ValueBool()
	tflog.Trace(ctx, "using 'Force'", map[string]interface{}{
		"Force": inputs.Force.ValueBool(),
	})

	if !inputs.Auth.IsNull() {
		options.Auth = authOptions(ctx, inputs.Auth, diag)
	}

	return options
}
