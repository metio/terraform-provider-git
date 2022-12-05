/*
 * SPDX-FileCopyrightText: The terraform-provider-git Authors
 * SPDX-License-Identifier: 0BSD
 */

package provider

import (
	"context"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/storage/filesystem"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/metio/terraform-provider-git/internal/modifiers"
	"os"
)

type InitResource struct{}

var (
	_ resource.Resource                = (*InitResource)(nil)
	_ resource.ResourceWithSchema      = (*InitResource)(nil)
	_ resource.ResourceWithImportState = (*InitResource)(nil)
)

type initResourceModel struct {
	Directory types.String `tfsdk:"directory"`
	Id        types.String `tfsdk:"id"`
	Bare      types.Bool   `tfsdk:"bare"`
}

func NewInitResource() resource.Resource {
	return &InitResource{}
}

func (r *InitResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_init"
}

func (r *InitResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Initializes a Git repository similar to 'git init'.",
		MarkdownDescription: "Initializes a Git repository similar to `git init`.",
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
			"id": schema.StringAttribute{
				Description:         "The import ID to import this resource which is equal to the value of the 'directory' attribute.",
				MarkdownDescription: "The import ID to import this resource which is equal to the value of the `directory` attribute.",
				Computed:            true,
			},
			"bare": schema.BoolAttribute{
				Description:         "Whether the created Git repository is bare or not. Defaults to 'false'.",
				MarkdownDescription: "Whether the created Git repository is bare or not. Defaults to `false`.",
				Computed:            true,
				Optional:            true,
				PlanModifiers: []planmodifier.Bool{
					modifiers.DefaultBool(false),
					boolplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *InitResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Create resource git_init")

	var inputs initResourceModel
	var state initResourceModel

	diags := req.Plan.Get(ctx, &inputs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	directory := inputs.Directory.ValueString()
	bare := inputs.Bare.ValueBool()

	_, err := git.PlainInit(directory, bare)
	if err != nil {
		resp.Diagnostics.AddError(
			"Cannot create repository",
			"Could not create repository ["+directory+"] because of: "+err.Error(),
		)
		return
	}

	tflog.Trace(ctx, "created repository", map[string]interface{}{
		"directory": directory,
		"bare":      bare,
	})

	state.Directory = inputs.Directory
	state.Id = inputs.Directory
	state.Bare = inputs.Bare

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *InitResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Read resource git_init")

	var state initResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	directory := state.Directory.ValueString()

	var newState initResourceModel
	newState.Directory = state.Directory
	newState.Id = state.Id

	repository := openRepository(ctx, directory, &resp.Diagnostics)
	if repository == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	worktree, err := getWorktree(repository, &resp.Diagnostics)
	if err != nil {
		return
	}
	newState.Bare = types.BoolValue(worktree == nil)

	diags = resp.State.Set(ctx, &newState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *InitResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Update resource git_init")
	updatedUsingPlan(ctx, &req, resp, &initResourceModel{})
}

func (r *InitResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Delete resource git_init")

	var state initResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	directory := state.Directory.ValueString()
	bare := state.Bare.ValueBool()

	if !bare {
		repository := openRepository(ctx, directory, &resp.Diagnostics)

		if repository.Storer != nil {
			storage, ok := repository.Storer.(*filesystem.Storage)

			if ok {
				err := os.RemoveAll(storage.Filesystem().Root())
				if err != nil {
					resp.Diagnostics.AddError(
						"Cannot delete repository",
						"Could not delete git repository ["+directory+"] because of: "+err.Error(),
					)
					return
				}
			}
		}
	}
}

func (r *InitResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Debug(ctx, "ImportState resource git_init")

	if req.ID == "" {
		resp.Diagnostics.AddError(
			"Unexpected import identifier",
			fmt.Sprintf("Expected import identifier with format: 'path/to/git/repository' Got: '%q'", req.ID),
		)
		return
	}

	var state initResourceModel
	state.Directory = types.StringValue(req.ID)
	state.Id = types.StringValue(req.ID)

	repository := openRepository(ctx, req.ID, &resp.Diagnostics)
	if repository == nil {
		return
	}
	worktree, err := getWorktree(repository, &resp.Diagnostics)
	if err != nil {
		return
	}
	state.Bare = types.BoolValue(worktree == nil)

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
