/*
 * SPDX-FileCopyrightText: The terraform-provider-git Authors
 * SPDX-License-Identifier: 0BSD
 */

package provider

import (
	"context"
	"errors"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/metio/terraform-provider-git/internal/modifiers"
)

type PushResource struct{}

var _ resource.Resource = (*PushResource)(nil)

type PushResourceModel struct {
	Directory        types.String `tfsdk:"directory"`
	Id               types.Int64  `tfsdk:"id"`
	Remote           types.String `tfsdk:"remote"`
	RefSpecs         types.List   `tfsdk:"refspecs"`
	Prune            types.Bool   `tfsdk:"prune"`
	Force            types.Bool   `tfsdk:"force"`
	InsecureSkipTls  types.Bool   `tfsdk:"insecure_skip_tls"`
	CaBundleFilePath types.String `tfsdk:"ca_bundle_file_path"`
	Auth             types.Object `tfsdk:"auth"`
}

func NewPushResource() resource.Resource {
	return &PushResource{}
}

func (r *PushResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_push"
}

func (r *PushResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Push changes to a Git remote similar to 'git push'",
		MarkdownDescription: "Push changes to a Git remote similar to `git push`",
		Attributes: map[string]schema.Attribute{
			"directory": schema.StringAttribute{
				Description:         "The path to the local Git repository.",
				MarkdownDescription: "The path to the local Git repository.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"id": schema.Int64Attribute{
				Description:         "The timestamp of the last push in Unix nanoseconds.",
				MarkdownDescription: "The timestamp of the last push in Unix nanoseconds.",
				Computed:            true,
			},
			"remote": schema.StringAttribute{
				Description:         "The name of the remote to push into. Defaults to 'origin'.",
				MarkdownDescription: "The name of the remote to push into. Defaults to `origin`.",
				Computed:            true,
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				PlanModifiers: []planmodifier.String{
					modifiers.DefaultString("origin"),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"refspecs": schema.ListAttribute{
				Description:         "Specify what destination refs to update with what source objects. Note that these must be fully qualified refspecs, e.g. 'refs/heads/master' instead of just 'master'.",
				MarkdownDescription: "Specify what destination refs to update with what source objects. Note that these must be fully qualified refspecs, e.g. `refs/heads/master` instead of just `master`.",
				ElementType:         types.StringType,
				Required:            true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
			},
			"prune": schema.BoolAttribute{
				Description:         "Remove remote branches that don’t have a local counterpart. Defaults to 'false'.",
				MarkdownDescription: "Remove remote branches that don’t have a local counterpart. Defaults to `false`.",
				Computed:            true,
				Optional:            true,
				PlanModifiers: []planmodifier.Bool{
					modifiers.DefaultBool(false),
					boolplanmodifier.RequiresReplace(),
				},
			},
			"force": schema.BoolAttribute{
				Description:         "Allow updating a remote ref that is not an ancestor of the local ref used to overwrite it. Can cause the remote repository to lose commits; use it with care. Defaults to 'false'.",
				MarkdownDescription: "Allow updating a remote ref that is not an ancestor of the local ref used to overwrite it. Can cause the remote repository to lose commits; use it with care. Defaults to `false`.",
				Computed:            true,
				Optional:            true,
				PlanModifiers: []planmodifier.Bool{
					modifiers.DefaultBool(false),
					boolplanmodifier.RequiresReplace(),
				},
			},
			"insecure_skip_tls": schema.BoolAttribute{
				Description:         "Skip SSL verification if protocol is HTTPS.",
				MarkdownDescription: "Skip SSL verification if protocol is HTTPS.",
				Computed:            true,
				Optional:            true,
				PlanModifiers: []planmodifier.Bool{
					modifiers.DefaultBool(false),
					boolplanmodifier.RequiresReplace(),
				},
			},
			"ca_bundle_file_path": schema.StringAttribute{
				Description:         "File system path to an additional CA bundle to use together with the system cert pool.",
				MarkdownDescription: "File system path to an additional CA bundle to use together with the system cert pool.",
				Computed:            true,
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					modifiers.DefaultString(""),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"auth": authResourceSchema(),
		},
	}
}

func (r *PushResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Create resource git_push")

	var inputs PushResourceModel
	diags := req.Plan.Get(ctx, &inputs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	directory := inputs.Directory.ValueString()

	repository := openRepository(ctx, directory, &resp.Diagnostics)
	if repository == nil {
		return
	}

	options := CreatePushOptions(ctx, &inputs, &resp.Diagnostics)
	if options == nil {
		return
	}

	err := repository.PushContext(ctx, options)
	if !errors.Is(err, git.NoErrAlreadyUpToDate) && err != nil {
		resp.Diagnostics.AddError(
			"Cannot push commits",
			"Could not push commits because of: "+err.Error(),
		)
		return
	}

	var state PushResourceModel
	state.Directory = inputs.Directory
	state.Id = types.Int64Value(time.Now().UnixNano())
	state.Remote = inputs.Remote
	state.RefSpecs = inputs.RefSpecs
	state.Prune = inputs.Prune
	state.Force = inputs.Force
	state.InsecureSkipTls = inputs.InsecureSkipTls
	state.CaBundleFilePath = inputs.CaBundleFilePath
	state.Auth = inputs.Auth

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
