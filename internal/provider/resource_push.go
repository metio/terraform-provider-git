/*
 * SPDX-FileCopyrightText: The terraform-provider-git Authors
 * SPDX-License-Identifier: 0BSD
 */

package provider

import (
	"context"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/metio/terraform-provider-git/internal/modifiers"
)

type PushResource struct{}

var (
	_ resource.Resource = (*PushResource)(nil)
)

type PushResourceModel struct {
	Directory types.String `tfsdk:"directory"`
	Id        types.Int64  `tfsdk:"id"`
	Remote    types.String `tfsdk:"remote"`
	RefSpecs  types.List   `tfsdk:"refspecs"`
	Prune     types.Bool   `tfsdk:"prune"`
	Force     types.Bool   `tfsdk:"force"`
	Auth      types.Object `tfsdk:"auth"`
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
			"auth": schema.SingleNestedAttribute{
				Description:         "The authentication credentials, if required, to use with the remote repository.",
				MarkdownDescription: "The authentication credentials, if required, to use with the remote repository.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"basic": schema.SingleNestedAttribute{
						Description:         "Configure basic auth authentication.",
						MarkdownDescription: "Configure basic auth authentication.",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"username": schema.StringAttribute{
								Description:         "The basic auth username.",
								MarkdownDescription: "The basic auth username.",
								Required:            true,
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.RequiresReplace(),
								},
							},
							"password": schema.StringAttribute{
								Description:         "The basic auth password.",
								MarkdownDescription: "The basic auth password.",
								Required:            true,
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.RequiresReplace(),
								},
							},
						},
						Validators: []validator.Object{
							objectvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("bearer")),
							objectvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("ssh_key")),
							objectvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("ssh_password")),
							objectvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("ssh_agent")),
							objectvalidator.ExactlyOneOf(
								path.MatchRelative().AtParent().AtName("basic"),
								path.MatchRelative().AtParent().AtName("bearer"),
								path.MatchRelative().AtParent().AtName("ssh_key"),
								path.MatchRelative().AtParent().AtName("ssh_password"),
								path.MatchRelative().AtParent().AtName("ssh_agent"),
							),
						},
						PlanModifiers: []planmodifier.Object{
							objectplanmodifier.RequiresReplace(),
						},
					},
					"bearer": schema.StringAttribute{
						Description:         "Configure HTTP bearer token authentication. Note that services like GitHub use basic auth with your OAuth2 personal access token as the password.",
						MarkdownDescription: "Configure HTTP bearer token authentication. **Note**: Services like GitHub use basic auth with your OAuth2 personal access token as the password.",
						Optional:            true,
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
							stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("basic")),
							stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("ssh_key")),
							stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("ssh_password")),
							stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("ssh_agent")),
							stringvalidator.ExactlyOneOf(
								path.MatchRelative().AtParent().AtName("basic"),
								path.MatchRelative().AtParent().AtName("bearer"),
								path.MatchRelative().AtParent().AtName("ssh_key"),
								path.MatchRelative().AtParent().AtName("ssh_password"),
								path.MatchRelative().AtParent().AtName("ssh_agent"),
							),
						},
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
					},
					"ssh_key": schema.SingleNestedAttribute{
						Description:         "Configure SSH public/private key authentication.",
						MarkdownDescription: "Configure SSH public/private key authentication.",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"username": schema.StringAttribute{
								Description:         "The SSH auth username.",
								MarkdownDescription: "The SSH auth username.",
								Optional:            true,
								Computed:            true,
								PlanModifiers: []planmodifier.String{
									modifiers.DefaultString("git"),
									stringplanmodifier.RequiresReplace(),
								},
							},
							"password": schema.StringAttribute{
								Description:         "The SSH key password.",
								MarkdownDescription: "The SSH key password.",
								Optional:            true,
								Computed:            true,
								PlanModifiers: []planmodifier.String{
									modifiers.DefaultString(""),
									stringplanmodifier.RequiresReplace(),
								},
							},
							"private_key_path": schema.StringAttribute{
								Description:         "The absolute path to the private SSH key.",
								MarkdownDescription: "The absolute path to the private SSH key.",
								Optional:            true,
								Validators: []validator.String{
									stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("private_key_pem")),
									stringvalidator.AtLeastOneOf(path.MatchRelative().AtParent().AtName("private_key_pem")),
								},
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.RequiresReplace(),
								},
							},
							"private_key_pem": schema.StringAttribute{
								Description:         "The private SSH key in PEM format.",
								MarkdownDescription: "The private SSH key in PEM format.",
								Optional:            true,
								Validators: []validator.String{
									stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("private_key_path")),
									stringvalidator.AtLeastOneOf(path.MatchRelative().AtParent().AtName("private_key_path")),
								},
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.RequiresReplace(),
								},
							},
							"known_hosts": schema.SetAttribute{
								Description:         "The list of known hosts files to accept. If none are specified, system defaults will be used.",
								MarkdownDescription: "The list of known hosts files to accept. If none are specified, system defaults will be used.",
								ElementType:         types.StringType,
								Optional:            true,
								PlanModifiers: []planmodifier.Set{
									setplanmodifier.RequiresReplace(),
								},
							},
						},
						Validators: []validator.Object{
							objectvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("basic")),
							objectvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("bearer")),
							objectvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("ssh_password")),
							objectvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("ssh_agent")),
							objectvalidator.ExactlyOneOf(
								path.MatchRelative().AtParent().AtName("basic"),
								path.MatchRelative().AtParent().AtName("bearer"),
								path.MatchRelative().AtParent().AtName("ssh_key"),
								path.MatchRelative().AtParent().AtName("ssh_password"),
								path.MatchRelative().AtParent().AtName("ssh_agent"),
							),
						},
					},
					"ssh_password": schema.SingleNestedAttribute{
						Description:         "Configure password based SSH authentication.",
						MarkdownDescription: "Configure password based SSH authentication.",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"username": schema.StringAttribute{
								Description:         "The SSH username.",
								MarkdownDescription: "The SSH username.",
								Required:            true,
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.RequiresReplace(),
								},
							},
							"password": schema.StringAttribute{
								Description:         "The SSH password.",
								MarkdownDescription: "The SSH password.",
								Required:            true,
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.RequiresReplace(),
								},
							},
							"known_hosts": schema.SetAttribute{
								Description:         "The list of known hosts files to accept. If none are specified, system defaults will be used.",
								MarkdownDescription: "The list of known hosts files to accept. If none are specified, system defaults will be used.",
								ElementType:         types.StringType,
								Optional:            true,
								PlanModifiers: []planmodifier.Set{
									setplanmodifier.RequiresReplace(),
								},
							},
						},
						Validators: []validator.Object{
							objectvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("basic")),
							objectvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("bearer")),
							objectvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("ssh_key")),
							objectvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("ssh_agent")),
							objectvalidator.ExactlyOneOf(
								path.MatchRelative().AtParent().AtName("basic"),
								path.MatchRelative().AtParent().AtName("bearer"),
								path.MatchRelative().AtParent().AtName("ssh_key"),
								path.MatchRelative().AtParent().AtName("ssh_password"),
								path.MatchRelative().AtParent().AtName("ssh_agent"),
							),
						},
						PlanModifiers: []planmodifier.Object{
							objectplanmodifier.RequiresReplace(),
						},
					},
					"ssh_agent": schema.SingleNestedAttribute{
						Description:         "Configure SSH agent based authentication.",
						MarkdownDescription: "Configure SSH agent based authentication.",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"username": schema.StringAttribute{
								Description:         "The system username of the user talking to the SSH agent. Use an empty string in order to automatically fetch this.",
								MarkdownDescription: "The system username of the user talking to the SSH agent. Use an empty string in order to automatically fetch this.",
								Optional:            true,
								Computed:            true,
								PlanModifiers: []planmodifier.String{
									modifiers.DefaultString(""),
									stringplanmodifier.RequiresReplace(),
								},
							},
							"known_hosts": schema.SetAttribute{
								Description:         "The list of known hosts files to accept. If none are specified, system defaults will be used.",
								MarkdownDescription: "The list of known hosts files to accept. If none are specified, system defaults will be used.",
								ElementType:         types.StringType,
								Optional:            true,
								PlanModifiers: []planmodifier.Set{
									setplanmodifier.RequiresReplace(),
								},
							},
						},
						Validators: []validator.Object{
							objectvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("basic")),
							objectvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("bearer")),
							objectvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("ssh_key")),
							objectvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("ssh_password")),
							objectvalidator.ExactlyOneOf(
								path.MatchRelative().AtParent().AtName("basic"),
								path.MatchRelative().AtParent().AtName("bearer"),
								path.MatchRelative().AtParent().AtName("ssh_key"),
								path.MatchRelative().AtParent().AtName("ssh_password"),
								path.MatchRelative().AtParent().AtName("ssh_agent"),
							),
						},
						PlanModifiers: []planmodifier.Object{
							objectplanmodifier.RequiresReplace(),
						},
					},
				},
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.RequiresReplace(),
				},
			},
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
	if err != git.NoErrAlreadyUpToDate && err != nil {
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
