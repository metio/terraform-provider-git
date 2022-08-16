/*
 * SPDX-FileCopyrightText: The terraform-provider-git Authors
 * SPDX-License-Identifier: 0BSD
 */

package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type dataSourceGitRemotesType struct{}

type dataSourceGitRemotes struct {
	p gitProvider
}

type dataSourceGitRemotesSchema struct {
	Directory types.String `tfsdk:"directory"`
	Id        types.String `tfsdk:"id"`
	Remotes   types.Map    `tfsdk:"remotes"`
}

func (r *dataSourceGitRemotesType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: "Reads all configured remotes of a Git repository.",
		Attributes: map[string]tfsdk.Attribute{
			"directory": {
				Description: "The path to the local Git repository.",
				Type:        types.StringType,
				Required:    true,
				Validators: []tfsdk.AttributeValidator{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"id": {
				MarkdownDescription: "The same value as the `directory` attribute.",
				Type:                types.StringType,
				Computed:            true,
			},
			"remotes": {
				Description: "All configured remotes of the given Git repository.",
				Computed:    true,
				Attributes: tfsdk.MapNestedAttributes(map[string]tfsdk.Attribute{
					"urls": {
						Description: "The URLs for the remote.",
						Type: types.ListType{
							ElemType: types.StringType,
						},
						Computed: true,
					},
				}),
			},
		},
	}, nil
}

func (r *dataSourceGitRemotesType) NewDataSource(_ context.Context, p provider.Provider) (datasource.DataSource, diag.Diagnostics) {
	return &dataSourceGitRemotes{
		p: *(p.(*gitProvider)),
	}, nil
}

func (r *dataSourceGitRemotes) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Reading Git repository remotes")

	var inputs dataSourceGitRemotesSchema
	var state dataSourceGitRemotesSchema

	diags := req.Config.Get(ctx, &inputs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	directory := inputs.Directory.Value

	repository := openRepository(ctx, directory, &resp.Diagnostics)
	if repository == nil {
		return
	}

	remotes, err := repository.Remotes()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading remotes",
			"Could not read remotes of ["+directory+"] because of: "+err.Error(),
		)
		return
	}

	tflog.Trace(ctx, "read remotes", map[string]interface{}{
		"directory": directory,
	})

	remoteType := map[string]attr.Type{
		"urls": types.ListType{ElemType: types.StringType},
	}

	allRemotes := make(map[string]attr.Value)
	for _, remote := range remotes {
		allRemotes[remote.Config().Name] = types.Object{
			AttrTypes: remoteType,
			Attrs: map[string]attr.Value{
				"urls": stringsToList(remote.Config().URLs),
			},
		}
	}

	state.Directory = inputs.Directory
	state.Id = inputs.Directory
	state.Remotes = types.Map{
		ElemType: types.ObjectType{
			AttrTypes: remoteType,
		},
		Elems: allRemotes,
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
