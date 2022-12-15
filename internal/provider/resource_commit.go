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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/metio/terraform-provider-git/internal/modifiers"
	"time"
)

type CommitResource struct{}

var (
	_ resource.Resource               = (*CommitResource)(nil)
	_ resource.ResourceWithModifyPlan = (*CommitResource)(nil)
)

type commitResourceModel struct {
	Directory         types.String `tfsdk:"directory"`
	Id                types.Int64  `tfsdk:"id"`
	Message           types.String `tfsdk:"message"`
	All               types.Bool   `tfsdk:"all"`
	AllowEmptyCommits types.Bool   `tfsdk:"allow_empty_commits"`
	Author            types.Object `tfsdk:"author"`
	Committer         types.Object `tfsdk:"committer"`
	SHA1              types.String `tfsdk:"sha1"`
	Files             types.List   `tfsdk:"files"`
}

func NewCommitResource() resource.Resource {
	return &CommitResource{}
}

func (r *CommitResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_commit"
}

func (r *CommitResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Record changes to the repository similar to 'git commit'. Note that configuration changes to this resource which cause a replacement will create a new commit and keep the previous commit as-is.",
		MarkdownDescription: "Record changes to the repository similar to `git commit`.\n\n" +
			"-> **Note** Configuration changes to this resource which cause a replacement will create a new commit and keep the previous commit as-is.",
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
				Description:         "The timestamp of the last commit in Unix nanoseconds.",
				MarkdownDescription: "The timestamp of the last commit in Unix nanoseconds.",
				Computed:            true,
			},
			"message": schema.StringAttribute{
				Description:         "The commit message to use.",
				MarkdownDescription: "The commit message to use.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"all": schema.BoolAttribute{
				Description:         "Automatically stage files that have been modified and deleted, but new files you have not told Git about are not affected. Defaults to 'false'.",
				MarkdownDescription: "Automatically stage files that have been modified and deleted, but new files you have not told Git about are not affected. Defaults to `false`.",
				Computed:            true,
				Optional:            true,
				PlanModifiers: []planmodifier.Bool{
					modifiers.DefaultBool(false),
					boolplanmodifier.RequiresReplace(),
				},
			},
			"allow_empty_commits": schema.BoolAttribute{
				Description:         "Enable empty commits to be created. Defaults to 'true'.",
				MarkdownDescription: "Enable empty commits to be created. Defaults to `true`.",
				Computed:            true,
				Optional:            true,
				PlanModifiers: []planmodifier.Bool{
					modifiers.DefaultBool(true),
					boolplanmodifier.RequiresReplace(),
				},
			},
			"author": schema.SingleNestedAttribute{
				Description:         "The original author of the commit. If none is specified, the author will be read from the Git configuration.",
				MarkdownDescription: "The original author of the commit. If none is specified, the author will be read from the Git configuration.",
				Computed:            true,
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						Description:         "The name of the author.",
						MarkdownDescription: "The name of the author.",
						Computed:            true,
						Optional:            true,
					},
					"email": schema.StringAttribute{
						Description:         "The email address of the author.",
						MarkdownDescription: "The email address of the author.",
						Computed:            true,
						Optional:            true,
					},
				},
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.RequiresReplace(),
				},
			},
			"committer": schema.SingleNestedAttribute{
				Description:         "The person performing the commit. If none is specified, the author is used as committer.",
				MarkdownDescription: "The person performing the commit. If none is specified, the author is used as committer.",
				Computed:            true,
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						Description:         "The name of the committer.",
						MarkdownDescription: "The name of the committer.",
						Computed:            true,
						Optional:            true,
					},
					"email": schema.StringAttribute{
						Description:         "The email address of the committer.",
						MarkdownDescription: "The email address of the committer.",
						Computed:            true,
						Optional:            true,
					},
				},
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.RequiresReplace(),
				},
			},
			"sha1": schema.StringAttribute{
				Description:         "The SHA1 hash of the created commit.",
				MarkdownDescription: "The SHA1 hash of the created commit.",
				Computed:            true,
			},
			"files": schema.ListAttribute{
				Description:         "The files updated by the commit.",
				MarkdownDescription: "The files updated by the commit.",
				ElementType:         types.StringType,
				Computed:            true,
			},
		},
	}
}

func (r *CommitResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Create resource git_commit")

	var inputs commitResourceModel
	diags := req.Config.Get(ctx, &inputs)
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
			"Cannot commit in bare repository",
			"The repository at ["+directory+"] is bare. Create a worktree first in order to commit changes.",
		)
		return
	}

	// NOTE: It seems default values are not working?
	if inputs.All.IsNull() {
		inputs.All = types.BoolValue(false)
	}
	if inputs.AllowEmptyCommits.IsNull() {
		inputs.AllowEmptyCommits = types.BoolValue(true)
	}

	status, err := worktree.Status()
	if err != nil {
		return
	}

	var state commitResourceModel
	state.Directory = inputs.Directory
	state.Id = types.Int64Value(time.Now().UnixNano())
	state.All = inputs.All
	state.AllowEmptyCommits = inputs.AllowEmptyCommits
	state.Message = inputs.Message
	state.Author = inputs.Author
	state.Committer = inputs.Committer
	state.Files = types.ListNull(types.StringType)
	state.SHA1 = types.StringNull()

	if !status.IsClean() {
		options := createCommitOptions(ctx, inputs)

		hash := createCommit(worktree, inputs.Message.ValueString(), options, &resp.Diagnostics)
		if hash == nil {
			return
		}

		commitObject := getCommit(ctx, repository, hash, &resp.Diagnostics)
		if commitObject == nil {
			return
		}

		state.Author = signatureToObjectWithoutTimestamp(&commitObject.Author)
		state.Committer = signatureToObjectWithoutTimestamp(&commitObject.Committer)
		state.Files, _ = types.ListValueFrom(ctx, types.StringType, extractModifiedFiles(commitObject))
		state.SHA1 = types.StringValue(hash.String())
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *CommitResource) Read(ctx context.Context, _ resource.ReadRequest, _ *resource.ReadResponse) {
	tflog.Debug(ctx, "Read resource git_add")
	// NO-OP: All data is already in Terraform state
}

func (r *CommitResource) Update(ctx context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
	tflog.Debug(ctx, "Update resource git_add")
	// NO-OP: All attributes require replacement, thus delete/create will be called
}

func (r *CommitResource) Delete(ctx context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
	tflog.Debug(ctx, "Delete resource git_add")
	// NO-OP: Terraform removes the state automatically for us
}

func (r *CommitResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	tflog.Debug(ctx, "ModifyPlan resource git_commit")

	if req.State.Raw.IsNull() {
		// if we're creating the resource, no need to modify it
		return
	}

	if req.Plan.Raw.IsNull() {
		// if we're deleting the resource, no need to modify it
		return
	}

	var inputs commitResourceModel
	diags := req.Config.Get(ctx, &inputs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	directory := inputs.Directory.ValueString()
	all := inputs.All.ValueBool()

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

	for _, val := range status {
		if (val.Staging != git.Unmodified && val.Staging != git.Untracked) ||
			(all && (val.Worktree == git.Modified || val.Worktree == git.Deleted)) {
			id := path.Root("id")
			resp.Plan.SetAttribute(ctx, id, time.Now().UnixNano())
			resp.RequiresReplace = append(resp.RequiresReplace, id)
			break
		}
	}
}
