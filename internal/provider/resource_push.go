/*
 * SPDX-FileCopyrightText: The terraform-provider-git Authors
 * SPDX-License-Identifier: 0BSD
 */

package provider

import (
	"context"
	"github.com/go-git/go-git/v5"
	"github.com/hashicorp/terraform-plugin-framework-validators/schemavalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
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
					modifiers.DefaultValue(types.StringValue("origin")),
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
					modifiers.DefaultValue(types.BoolValue(false)),
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
					modifiers.DefaultValue(types.BoolValue(false)),
					resource.RequiresReplace(),
				},
			},
			"auth": {
				Description:         "The authentication credentials, if required, to use with the remote repository.",
				MarkdownDescription: "The authentication credentials, if required, to use with the remote repository.",
				Optional:            true,
				Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
					"basic": {
						Description:         "Configure basic auth authentication.",
						MarkdownDescription: "Configure basic auth authentication.",
						Optional:            true,
						Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
							"username": {
								Description:         "The basic auth username.",
								MarkdownDescription: "The basic auth username.",
								Type:                types.StringType,
								Required:            true,
								PlanModifiers: []tfsdk.AttributePlanModifier{
									resource.RequiresReplace(),
								},
							},
							"password": {
								Description:         "The basic auth password.",
								MarkdownDescription: "The basic auth password.",
								Type:                types.StringType,
								Required:            true,
								PlanModifiers: []tfsdk.AttributePlanModifier{
									resource.RequiresReplace(),
								},
							},
						}),
						Validators: []tfsdk.AttributeValidator{
							schemavalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("bearer")),
							schemavalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("ssh_key")),
							schemavalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("ssh_password")),
							schemavalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("ssh_agent")),
							schemavalidator.ExactlyOneOf(
								path.MatchRelative().AtParent().AtName("basic"),
								path.MatchRelative().AtParent().AtName("bearer"),
								path.MatchRelative().AtParent().AtName("ssh_key"),
								path.MatchRelative().AtParent().AtName("ssh_password"),
								path.MatchRelative().AtParent().AtName("ssh_agent"),
							),
						},
						PlanModifiers: []tfsdk.AttributePlanModifier{
							resource.RequiresReplace(),
						},
					},
					"bearer": {
						Description:         "Configure HTTP bearer token authentication. Note that services like GitHub use basic auth with your OAuth2 personal access token as the password.",
						MarkdownDescription: "Configure HTTP bearer token authentication. **Note**: Services like GitHub use basic auth with your OAuth2 personal access token as the password.",
						Type:                types.StringType,
						Optional:            true,
						Validators: []tfsdk.AttributeValidator{
							stringvalidator.LengthAtLeast(1),
							schemavalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("basic")),
							schemavalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("ssh_key")),
							schemavalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("ssh_password")),
							schemavalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("ssh_agent")),
							schemavalidator.ExactlyOneOf(
								path.MatchRelative().AtParent().AtName("basic"),
								path.MatchRelative().AtParent().AtName("bearer"),
								path.MatchRelative().AtParent().AtName("ssh_key"),
								path.MatchRelative().AtParent().AtName("ssh_password"),
								path.MatchRelative().AtParent().AtName("ssh_agent"),
							),
						},
						PlanModifiers: []tfsdk.AttributePlanModifier{
							resource.RequiresReplace(),
						},
					},
					"ssh_key": {
						Description:         "Configure SSH public/private key authentication.",
						MarkdownDescription: "Configure SSH public/private key authentication.",
						Optional:            true,
						Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
							"username": {
								Description:         "The SSH auth username.",
								MarkdownDescription: "The SSH auth username.",
								Type:                types.StringType,
								Optional:            true,
								Computed:            true,
								PlanModifiers: []tfsdk.AttributePlanModifier{
									modifiers.DefaultValue(types.StringValue("git")),
									resource.RequiresReplace(),
								},
							},
							"password": {
								Description:         "The SSH key password.",
								MarkdownDescription: "The SSH key password.",
								Type:                types.StringType,
								Optional:            true,
								Computed:            true,
								PlanModifiers: []tfsdk.AttributePlanModifier{
									modifiers.DefaultValue(types.StringValue("")),
									resource.RequiresReplace(),
								},
							},
							"private_key_path": {
								Description:         "The absolute path to the private SSH key.",
								MarkdownDescription: "The absolute path to the private SSH key.",
								Type:                types.StringType,
								Optional:            true,
								Validators: []tfsdk.AttributeValidator{
									schemavalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("private_key_pem")),
									schemavalidator.AtLeastOneOf(path.MatchRelative().AtParent().AtName("private_key_pem")),
								},
								PlanModifiers: []tfsdk.AttributePlanModifier{
									resource.RequiresReplace(),
								},
							},
							"private_key_pem": {
								Description:         "The private SSH key in PEM format.",
								MarkdownDescription: "The private SSH key in PEM format.",
								Type:                types.StringType,
								Optional:            true,
								Validators: []tfsdk.AttributeValidator{
									schemavalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("private_key_path")),
									schemavalidator.AtLeastOneOf(path.MatchRelative().AtParent().AtName("private_key_path")),
								},
								PlanModifiers: []tfsdk.AttributePlanModifier{
									resource.RequiresReplace(),
								},
							},
							"known_hosts": {
								Description:         "The list of known hosts to accept. If none are specified, system defaults will be used.",
								MarkdownDescription: "The list of known hosts to accept. If none are specified, system defaults will be used.",
								Type: types.SetType{
									ElemType: types.StringType,
								},
								Optional: true,
								PlanModifiers: []tfsdk.AttributePlanModifier{
									resource.RequiresReplace(),
								},
							},
						}),
						Validators: []tfsdk.AttributeValidator{
							schemavalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("basic")),
							schemavalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("bearer")),
							schemavalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("ssh_password")),
							schemavalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("ssh_agent")),
							schemavalidator.ExactlyOneOf(
								path.MatchRelative().AtParent().AtName("basic"),
								path.MatchRelative().AtParent().AtName("bearer"),
								path.MatchRelative().AtParent().AtName("ssh_key"),
								path.MatchRelative().AtParent().AtName("ssh_password"),
								path.MatchRelative().AtParent().AtName("ssh_agent"),
							),
						},
					},
					"ssh_password": {
						Description:         "Configure password based SSH authentication.",
						MarkdownDescription: "Configure password based SSH authentication.",
						Optional:            true,
						Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
							"username": {
								Description:         "The SSH username.",
								MarkdownDescription: "The SSH username.",
								Type:                types.StringType,
								Required:            true,
								PlanModifiers: []tfsdk.AttributePlanModifier{
									resource.RequiresReplace(),
								},
							},
							"password": {
								Description:         "The SSH password.",
								MarkdownDescription: "The SSH password.",
								Type:                types.StringType,
								Required:            true,
								PlanModifiers: []tfsdk.AttributePlanModifier{
									resource.RequiresReplace(),
								},
							},
							"known_hosts": {
								Description:         "The list of known hosts to accept. If none are specified, system defaults will be used.",
								MarkdownDescription: "The list of known hosts to accept. If none are specified, system defaults will be used.",
								Type: types.SetType{
									ElemType: types.StringType,
								},
								Optional: true,
								PlanModifiers: []tfsdk.AttributePlanModifier{
									resource.RequiresReplace(),
								},
							},
						}),
						Validators: []tfsdk.AttributeValidator{
							schemavalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("basic")),
							schemavalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("bearer")),
							schemavalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("ssh_key")),
							schemavalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("ssh_agent")),
							schemavalidator.ExactlyOneOf(
								path.MatchRelative().AtParent().AtName("basic"),
								path.MatchRelative().AtParent().AtName("bearer"),
								path.MatchRelative().AtParent().AtName("ssh_key"),
								path.MatchRelative().AtParent().AtName("ssh_password"),
								path.MatchRelative().AtParent().AtName("ssh_agent"),
							),
						},
						PlanModifiers: []tfsdk.AttributePlanModifier{
							resource.RequiresReplace(),
						},
					},
					"ssh_agent": {
						Description:         "Configure SSH agent based authentication.",
						MarkdownDescription: "Configure SSH agent based authentication.",
						Optional:            true,
						Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
							"username": {
								Description:         "The system username of the user talking to the SSH agent. Use an empty string in order to automatically fetch this.",
								MarkdownDescription: "The system username of the user talking to the SSH agent. Use an empty string in order to automatically fetch this.",
								Type:                types.StringType,
								Optional:            true,
								Computed:            true,
								PlanModifiers: []tfsdk.AttributePlanModifier{
									modifiers.DefaultValue(types.String{Value: ""}),
									resource.RequiresReplace(),
								},
							},
							"known_hosts": {
								Description:         "The list of known hosts to accept. If none are specified, system defaults will be used.",
								MarkdownDescription: "The list of known hosts to accept. If none are specified, system defaults will be used.",
								Type: types.SetType{
									ElemType: types.StringType,
								},
								Optional: true,
								PlanModifiers: []tfsdk.AttributePlanModifier{
									resource.RequiresReplace(),
								},
							},
						}),
						Validators: []tfsdk.AttributeValidator{
							schemavalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("basic")),
							schemavalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("bearer")),
							schemavalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("ssh_key")),
							schemavalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("ssh_password")),
							schemavalidator.ExactlyOneOf(
								path.MatchRelative().AtParent().AtName("basic"),
								path.MatchRelative().AtParent().AtName("bearer"),
								path.MatchRelative().AtParent().AtName("ssh_key"),
								path.MatchRelative().AtParent().AtName("ssh_password"),
								path.MatchRelative().AtParent().AtName("ssh_agent"),
							),
						},
						PlanModifiers: []tfsdk.AttributePlanModifier{
							resource.RequiresReplace(),
						},
					},
				}),
				PlanModifiers: []tfsdk.AttributePlanModifier{
					resource.RequiresReplace(),
				},
			},
		},
	}, nil
}

func (r *PushResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Create resource git_push")

	var inputs PushResourceModel
	diags := req.Config.Get(ctx, &inputs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// NOTE: It seems default values are not working?
	if inputs.Remote.IsNull() {
		inputs.Remote = types.StringValue("origin")
	}
	if inputs.Prune.IsNull() {
		inputs.Prune = types.BoolValue(false)
	}
	if inputs.Force.IsNull() {
		inputs.Force = types.BoolValue(false)
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
