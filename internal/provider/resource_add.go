/*
 * SPDX-FileCopyrightText: The terraform-provider-git Authors
 * SPDX-License-Identifier: 0BSD
 */

package provider

import (
	"context"
	"github.com/go-git/go-git/v5"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"path/filepath"
	"time"
)

type AddResource struct{}

var (
	_ resource.Resource               = (*AddResource)(nil)
	_ resource.ResourceWithModifyPlan = (*AddResource)(nil)
)

type addResourceModel struct {
	Directory types.String `tfsdk:"directory"`
	Id        types.Int64  `tfsdk:"id"`
	Paths     types.List   `tfsdk:"add_paths"`
}

func NewAddResource() resource.Resource {
	return &AddResource{}
}

func (r *AddResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_add"
}

func (r *AddResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Add file contents to the index similar to 'git add'.",
		MarkdownDescription: "Add file contents to the index similar to `git add`.",
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
				Description:         "The timestamp of the last addition in Unix nanoseconds.",
				MarkdownDescription: "The timestamp of the last addition in Unix nanoseconds.",
				Computed:            true,
			},
			"add_paths": schema.ListAttribute{
				Description:         "The paths to add to the Git index. Values can be exact paths or glob patterns.",
				MarkdownDescription: "The paths to add to the Git index. Values can be exact paths or glob patterns.",
				ElementType:         types.StringType,
				Required:            true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *AddResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Create resource git_add")

	var inputs addResourceModel
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

	worktree, err := getWorktree(repository, &resp.Diagnostics)
	if err != nil {
		return
	}
	if worktree == nil {
		resp.Diagnostics.AddError(
			"Cannot add file to bare repository",
			"The repository at ["+directory+"] is bare. Create a worktree first in order to add files to it.",
		)
		return
	}

	status := getStatus(ctx, worktree, &resp.Diagnostics)
	if status == nil {
		return
	}

	paths := make([]string, len(inputs.Paths.Elements()))
	resp.Diagnostics.Append(inputs.Paths.ElementsAs(ctx, &paths, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	for _, pattern := range paths {
		for file, fileStatus := range status {
			if fileStatus.Worktree != git.Unmodified {
				match, errMatch := filepath.Match(pattern, file)
				if errMatch != nil {
					resp.Diagnostics.AddError(
						"Cannot match file path",
						"Could not match pattern ["+pattern+"] because of: "+errMatch.Error(),
					)
				}
				if match {
					_, errAdd := worktree.Add(file)
					if errAdd != nil {
						resp.Diagnostics.AddError(
							"Cannot add file",
							"Could not add file ["+file+"] because of: "+errAdd.Error(),
						)
					}
				}
			}
		}
	}

	var state addResourceModel
	state.Directory = inputs.Directory
	state.Id = types.Int64Value(time.Now().UnixNano())
	state.Paths = inputs.Paths

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *AddResource) Read(ctx context.Context, _ resource.ReadRequest, _ *resource.ReadResponse) {
	tflog.Debug(ctx, "Read resource git_add")
	// NO-OP: All data is already in Terraform state
}

func (r *AddResource) Update(ctx context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
	tflog.Debug(ctx, "Update resource git_add")
	// NO-OP: All attributes require replacement, thus delete/create will be called
}

func (r *AddResource) Delete(ctx context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
	tflog.Debug(ctx, "Delete resource git_add")
	// NO-OP: Terraform removes the state automatically for us
}

func (r *AddResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	tflog.Debug(ctx, "ModifyPlan resource git_add")

	if req.State.Raw.IsNull() {
		// if we're creating the resource, no need to modify it
		return
	}

	if req.Plan.Raw.IsNull() {
		// if we're deleting the resource, no need to modify it
		return
	}

	var inputs addResourceModel
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

	worktree, err := getWorktree(repository, &resp.Diagnostics)
	if err != nil || worktree == nil {
		return
	}

	status := getStatus(ctx, worktree, &resp.Diagnostics)
	if status == nil {
		return
	}

	paths := make([]string, len(inputs.Paths.Elements()))
	resp.Diagnostics.Append(inputs.Paths.ElementsAs(ctx, &paths, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	for _, pattern := range paths {
		for key, val := range status {
			if val.Worktree != git.Unmodified {
				match, errMatch := filepath.Match(pattern, key)
				if errMatch != nil {
					resp.Diagnostics.AddError(
						"Cannot match file path",
						"Could not match pattern ["+pattern+"] because of: "+errMatch.Error(),
					)
					return
				}
				if match {
					id := path.Root("id")
					resp.Plan.SetAttribute(ctx, id, time.Now().UnixNano())
					resp.RequiresReplace = append(resp.RequiresReplace, id)
					break
				}
			}
		}
	}
}
