/*
 * SPDX-FileCopyrightText: The terraform-provider-git Authors
 * SPDX-License-Identifier: 0BSD
 */

package provider

import (
	"context"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/metio/terraform-provider-git/internal/modifiers"
	"time"
)

type PushResource struct{}

var (
	_ resource.Resource = (*PushResource)(nil)
)

type pushResourceModel struct {
	Directory types.String `tfsdk:"directory"`
	Id        types.Int64  `tfsdk:"id"`
	Remote    types.String `tfsdk:"remote"`
	RefSpecs  types.List   `tfsdk:"refspecs"`
	Prune     types.Bool   `tfsdk:"prune"`
	Force     types.Bool   `tfsdk:"force"`
}

func NewPushResource() resource.Resource {
	return &PushResource{}
}

func (r *PushResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_push"
}

func (r *PushResource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description:         "Push changes to a Git remote similar to 'git push'",
		MarkdownDescription: "Push changes to a Git remote similar to `git push`",
		Attributes: map[string]tfsdk.Attribute{
			"directory": {
				Description:         "The path to the local Git repository.",
				MarkdownDescription: "The path to the local Git repository.",
				Type:                types.StringType,
				Required:            true,
				Validators: []tfsdk.AttributeValidator{
					stringvalidator.LengthAtLeast(1),
				},
				PlanModifiers: []tfsdk.AttributePlanModifier{
					resource.RequiresReplace(),
				},
			},
			"id": {
				Description:         "The timestamp of the last push in Unix nanoseconds.",
				MarkdownDescription: "The timestamp of the last push in Unix nanoseconds.",
				Type:                types.Int64Type,
				Computed:            true,
			},
			"remote": {
				Description:         "The name of the remote to push into. Defaults to 'origin'.",
				MarkdownDescription: "The name of the remote to push into. Defaults to `origin`.",
				Type:                types.StringType,
				Computed:            true,
				Optional:            true,
				Validators: []tfsdk.AttributeValidator{
					stringvalidator.LengthAtLeast(1),
				},
				PlanModifiers: []tfsdk.AttributePlanModifier{
					modifiers.DefaultValue(types.String{Value: "origin"}),
					resource.RequiresReplace(),
				},
			},
			"refspecs": {
				Description:         "Specify what destination refs to update with what source objects. Note that these must be fully qualified refspecs, e.g. 'refs/heads/master' instead of just 'master'.",
				MarkdownDescription: "Specify what destination refs to update with what source objects. Note that these must be fully qualified refspecs, e.g. `refs/heads/master` instead of just `master`.",
				Type: types.ListType{
					ElemType: types.StringType,
				},
				Required: true,
				PlanModifiers: []tfsdk.AttributePlanModifier{
					resource.RequiresReplace(),
				},
			},
			"prune": {
				Description:         "Remove remote branches that don’t have a local counterpart. Defaults to 'false'.",
				MarkdownDescription: "Remove remote branches that don’t have a local counterpart. Defaults to `false`.",
				Type:                types.BoolType,
				Computed:            true,
				Optional:            true,
				PlanModifiers: []tfsdk.AttributePlanModifier{
					modifiers.DefaultValue(types.Bool{Value: false}),
					resource.RequiresReplace(),
				},
			},
			"force": {
				Description:         "Allow updating a remote ref that is not an ancestor of the local ref used to overwrite it. Can cause the remote repository to lose commits; use it with care. Defaults to 'false'.",
				MarkdownDescription: "Allow updating a remote ref that is not an ancestor of the local ref used to overwrite it. Can cause the remote repository to lose commits; use it with care. Defaults to `false`.",
				Type:                types.BoolType,
				Computed:            true,
				Optional:            true,
				PlanModifiers: []tfsdk.AttributePlanModifier{
					modifiers.DefaultValue(types.Bool{Value: false}),
					resource.RequiresReplace(),
				},
			},
		},
	}, nil
}

func (r *PushResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Create resource git_push")

	var inputs pushResourceModel
	diags := req.Config.Get(ctx, &inputs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// NOTE: It seems default values are not working?
	if inputs.Remote.IsNull() {
		inputs.Remote = types.String{Value: "origin"}
	}
	if inputs.Prune.IsNull() {
		inputs.Prune = types.Bool{Value: false}
	}
	if inputs.Force.IsNull() {
		inputs.Force = types.Bool{Value: false}
	}

	directory := inputs.Directory.Value

	repository := openRepository(ctx, directory, &resp.Diagnostics)
	if repository == nil {
		return
	}

	refspecs := make([]config.RefSpec, len(inputs.RefSpecs.Elems))
	resp.Diagnostics.Append(inputs.RefSpecs.ElementsAs(ctx, &refspecs, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := repository.PushContext(ctx, &git.PushOptions{
		RemoteName: inputs.Remote.Value,
		RefSpecs:   refspecs,
		Prune:      inputs.Prune.Value,
		Force:      inputs.Force.Value,
	})
	if err != git.NoErrAlreadyUpToDate && err != nil {
		resp.Diagnostics.AddError(
			"Cannot push commits",
			"Could not push commits because of: "+err.Error(),
		)
		return
	}

	var state pushResourceModel
	state.Directory = inputs.Directory
	state.Id = types.Int64{Value: time.Now().UnixNano()}
	state.Remote = inputs.Remote
	state.RefSpecs = inputs.RefSpecs
	state.Prune = inputs.Prune
	state.Force = inputs.Force

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *PushResource) Read(ctx context.Context, _ resource.ReadRequest, _ *resource.ReadResponse) {
	tflog.Debug(ctx, "Read resource git_push")
	// NO-OP: All data is already in Terraform state
}

func (r *PushResource) Update(ctx context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
	tflog.Debug(ctx, "Update resource git_push")
	// NO-OP: All attributes require replacement, thus delete/create will be called
}

func (r *PushResource) Delete(ctx context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
	tflog.Debug(ctx, "Delete resource git_push")
	// NO-OP: Terraform removes the state automatically for us
}
