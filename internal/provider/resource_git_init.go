/*
 * This file is part of terraform-provider-git. It is subject to the license terms in the LICENSE file found in the top-level
 * directory of this distribution and at https://creativecommons.org/publicdomain/zero/1.0/. No part of terraform-provider-git,
 * including this file, may be copied, modified, propagated, or distributed except according to the terms contained
 * in the LICENSE file.
 */

package provider

import (
	"context"
	"github.com/go-git/go-git/v5"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/metio/terraform-provider-git/internal/modifiers"
)

type resourceGitInitType struct{}

type resourceGitInit struct {
	p gitProvider
}

type resourceGitInitSchema struct {
	Directory types.String `tfsdk:"directory"`
	Id        types.String `tfsdk:"id"`
	Bare      types.Bool   `tfsdk:"bare"`
}

func (c *resourceGitInitType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"directory": {
				Description: "The path to the local Git repository.",
				Type:        types.StringType,
				Required:    true,
				Validators: []tfsdk.AttributeValidator{
					stringvalidator.LengthAtLeast(1),
				},
				PlanModifiers: []tfsdk.AttributePlanModifier{
					tfsdk.RequiresReplace(),
				},
			},
			"id": {
				MarkdownDescription: "`DEPRECATED`: Only added in order to use the sdkv2 test framework. The path to the local Git repository.",
				Type:                types.StringType,
				Computed:            true,
			},
			"bare": {
				Description: "Whether the created Git repository is bare or not. Defaults to `false`.",
				Type:        types.BoolType,
				Computed:    true,
				Optional:    true,
				PlanModifiers: []tfsdk.AttributePlanModifier{
					modifiers.DefaultValue(types.Bool{Value: false}),
					tfsdk.RequiresReplace(),
				},
			},
		},
	}, nil
}

func (r *resourceGitInitType) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	return &resourceGitInit{
		p: *(p.(*gitProvider)),
	}, nil
}

func (r *resourceGitInit) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	tflog.Debug(ctx, "Creating Git repository")

	var state resourceGitInitSchema
	var output resourceGitInitSchema

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	directory := state.Directory.Value
	bare := state.Bare.Value

	_, err := git.PlainInit(directory, bare)
	if err != nil {
		resp.Diagnostics.AddError(
			"Cannot create repository",
			"Could not create repository ["+directory+"] because of: "+err.Error(),
		)
	}

	tflog.Trace(ctx, "created repository", map[string]interface{}{
		"directory": directory,
		"bare":      bare,
	})

	output.Directory = types.String{Value: directory}
	output.Id = types.String{Value: directory}
	output.Bare = types.Bool{Value: bare}

	diags = resp.State.Set(ctx, &output)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceGitInit) Read(ctx context.Context, _ tfsdk.ReadResourceRequest, _ *tfsdk.ReadResourceResponse) {
	// NO-OP: all there is to read is in the State, and response is already populated with that.
	tflog.Debug(ctx, "Reading Git repository from state")
}

func (r *resourceGitInit) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	tflog.Debug(ctx, "Updating Git repository")

	updatedUsingPlan(ctx, &req, resp, &resourceGitInitSchema{})
}

func (r *resourceGitInit) Delete(ctx context.Context, _ tfsdk.DeleteResourceRequest, _ *tfsdk.DeleteResourceResponse) {
	tflog.Debug(ctx, "Removing Git repository")
}
