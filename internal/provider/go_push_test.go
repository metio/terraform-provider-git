/*
 * SPDX-FileCopyrightText: The terraform-provider-git Authors
 * SPDX-License-Identifier: 0BSD
 */

package provider_test

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"net"
	"os"
	"path/filepath"
	"testing"

	gitssh "github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/metio/terraform-provider-git/internal/provider"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

func TestCreatePushOptions_EmptyModel(t *testing.T) {
	ctx := context.TODO()
	model := &provider.PushResourceModel{}
	diagnostics := &diag.Diagnostics{}

	options := provider.CreatePushOptions(ctx, model, "", diagnostics)

	assert.Nil(t, options)
	assert.False(t, diagnostics.HasError())
}

func TestCreatePushOptions_RefSpecs(t *testing.T) {
	ctx := context.TODO()
	model := &provider.PushResourceModel{}
	var diagnostics diag.Diagnostics

	model.RefSpecs, diagnostics = types.ListValueFrom(ctx, types.StringType, []string{"refs/heads/main"})
	model.Auth = types.ObjectNull(map[string]attr.Type{})
	options := provider.CreatePushOptions(ctx, model, "", &diagnostics)

	assert.NotNil(t, options)
	assert.False(t, diagnostics.HasError())
	assert.Equal(t, 1, len(options.RefSpecs))
	assert.Equal(t, "refs/heads/main", options.RefSpecs[0].String())
	assert.Empty(t, options.RemoteName)
	assert.False(t, options.Prune)
	assert.False(t, options.Force)
	assert.Nil(t, options.Auth)
}

func TestCreatePushOptions_Remote(t *testing.T) {
	ctx := context.TODO()
	model := &provider.PushResourceModel{}
	diagnostics := &diag.Diagnostics{}

	model.RefSpecs, _ = types.ListValueFrom(ctx, types.StringType, []string{"refs/heads/main"})
	model.Auth = types.ObjectNull(map[string]attr.Type{})
	model.Remote = types.StringValue("origin")
	options := provider.CreatePushOptions(ctx, model, "", diagnostics)

	assert.NotNil(t, options)
	assert.False(t, diagnostics.HasError())
	assert.Equal(t, 1, len(options.RefSpecs))
	assert.Equal(t, "refs/heads/main", options.RefSpecs[0].String())
	assert.Equal(t, "origin", options.RemoteName)
	assert.False(t, options.Prune)
	assert.False(t, options.Force)
	assert.Nil(t, options.Auth)
}

func TestCreatePushOptions_Prune(t *testing.T) {
	ctx := context.TODO()
	model := &provider.PushResourceModel{}
	diagnostics := &diag.Diagnostics{}

	model.RefSpecs, _ = types.ListValueFrom(ctx, types.StringType, []string{"refs/heads/main"})
	model.Auth = types.ObjectNull(map[string]attr.Type{})
	model.Prune = types.BoolValue(true)
	options := provider.CreatePushOptions(ctx, model, "", diagnostics)

	assert.NotNil(t, options)
	assert.False(t, diagnostics.HasError())
	assert.Equal(t, 1, len(options.RefSpecs))
	assert.Equal(t, "refs/heads/main", options.RefSpecs[0].String())
	assert.Empty(t, options.RemoteName)
	assert.True(t, options.Prune)
	assert.False(t, options.Force)
	assert.Nil(t, options.Auth)
}

func TestCreatePushOptions_Force(t *testing.T) {
	ctx := context.TODO()
	model := &provider.PushResourceModel{}
	diagnostics := &diag.Diagnostics{}

	model.RefSpecs, _ = types.ListValueFrom(ctx, types.StringType, []string{"refs/heads/main"})
	model.Auth = types.ObjectNull(map[string]attr.Type{})
	model.Force = types.BoolValue(true)
	options := provider.CreatePushOptions(ctx, model, "", diagnostics)

	assert.NotNil(t, options)
	assert.False(t, diagnostics.HasError())
	assert.Equal(t, 1, len(options.RefSpecs))
	assert.Equal(t, "refs/heads/main", options.RefSpecs[0].String())
	assert.Empty(t, options.RemoteName)
	assert.False(t, options.Prune)
	assert.True(t, options.Force)
	assert.Nil(t, options.Auth)
}

func TestCreatePushOptions_Auth_Empty(t *testing.T) {
	ctx := context.TODO()
	model := &provider.PushResourceModel{}
	diagnostics := &diag.Diagnostics{}

	model.RefSpecs, _ = types.ListValueFrom(ctx, types.StringType, []string{"refs/heads/main"})
	model.Auth = types.ObjectNull(map[string]attr.Type{})
	options := provider.CreatePushOptions(ctx, model, "", diagnostics)

	assert.NotNil(t, options)
	assert.False(t, diagnostics.HasError())
	assert.Equal(t, 1, len(options.RefSpecs))
	assert.Equal(t, "refs/heads/main", options.RefSpecs[0].String())
	assert.Empty(t, options.RemoteName)
	assert.False(t, options.Prune)
	assert.False(t, options.Force)
	assert.Nil(t, options.Auth)
}

func TestCreatePushOptions_Auth_BearerAuth_Null(t *testing.T) {
	ctx := context.TODO()
	model := &provider.PushResourceModel{}
	diagnostics := &diag.Diagnostics{}

	model.RefSpecs, _ = types.ListValueFrom(ctx, types.StringType, []string{"refs/heads/main"})
	model.Auth = types.ObjectValueMust(
		map[string]attr.Type{
			"bearer": types.StringType,
		},
		map[string]attr.Value{
			"bearer": types.StringNull(),
		},
	)
	options := provider.CreatePushOptions(ctx, model, "", diagnostics)

	assert.NotNil(t, options)
	assert.False(t, diagnostics.HasError())
	assert.Equal(t, 1, len(options.RefSpecs))
	assert.Equal(t, "refs/heads/main", options.RefSpecs[0].String())
	assert.Empty(t, options.RemoteName)
	assert.False(t, options.Prune)
	assert.False(t, options.Force)
	assert.Nil(t, options.Auth)
}

func TestCreatePushOptions_Auth_BearerAuth_Token(t *testing.T) {
	ctx := context.TODO()
	model := &provider.PushResourceModel{}
	diagnostics := &diag.Diagnostics{}

	model.RefSpecs, _ = types.ListValueFrom(ctx, types.StringType, []string{"refs/heads/main"})
	model.Auth = types.ObjectValueMust(
		map[string]attr.Type{
			"bearer": types.StringType,
		},
		map[string]attr.Value{
			"bearer": types.StringValue("secret-token"),
		},
	)
	options := provider.CreatePushOptions(ctx, model, "", diagnostics)

	assert.NotNil(t, options)
	assert.False(t, diagnostics.HasError())
	assert.Equal(t, 1, len(options.RefSpecs))
	assert.Equal(t, "refs/heads/main", options.RefSpecs[0].String())
	assert.Empty(t, options.RemoteName)
	assert.False(t, options.Prune)
	assert.False(t, options.Force)
	assert.NotNil(t, options.Auth)
}

func TestCreatePushOptions_Auth_BasicAuth_Null(t *testing.T) {
	ctx := context.TODO()
	model := &provider.PushResourceModel{}
	diagnostics := &diag.Diagnostics{}

	model.RefSpecs, _ = types.ListValueFrom(ctx, types.StringType, []string{"refs/heads/main"})
	model.Auth = types.ObjectValueMust(
		map[string]attr.Type{
			"basic": types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"username": types.StringType,
					"password": types.StringType,
				},
			},
		},
		map[string]attr.Value{
			"basic": types.ObjectNull(map[string]attr.Type{
				"username": types.StringType,
				"password": types.StringType,
			}),
		},
	)
	options := provider.CreatePushOptions(ctx, model, "", diagnostics)

	assert.NotNil(t, options)
	assert.False(t, diagnostics.HasError())
	assert.Equal(t, 1, len(options.RefSpecs))
	assert.Equal(t, "refs/heads/main", options.RefSpecs[0].String())
	assert.Empty(t, options.RemoteName)
	assert.False(t, options.Prune)
	assert.False(t, options.Force)
	assert.Nil(t, options.Auth)
}

func TestCreatePushOptions_Auth_BasicAuth_Valid(t *testing.T) {
	ctx := context.TODO()
	model := &provider.PushResourceModel{}
	diagnostics := &diag.Diagnostics{}

	model.RefSpecs, _ = types.ListValueFrom(ctx, types.StringType, []string{"refs/heads/main"})
	model.Auth = types.ObjectValueMust(
		map[string]attr.Type{
			"basic": types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"username": types.StringType,
					"password": types.StringType,
				},
			},
		},
		map[string]attr.Value{
			"basic": types.ObjectValueMust(
				map[string]attr.Type{
					"username": types.StringType,
					"password": types.StringType,
				},
				map[string]attr.Value{
					"username": types.StringValue("user"),
					"password": types.StringValue("secret"),
				},
			),
		},
	)
	options := provider.CreatePushOptions(ctx, model, "", diagnostics)

	assert.NotNil(t, options)
	assert.False(t, diagnostics.HasError())
	assert.Equal(t, 1, len(options.RefSpecs))
	assert.Equal(t, "refs/heads/main", options.RefSpecs[0].String())
	assert.Empty(t, options.RemoteName)
	assert.False(t, options.Prune)
	assert.False(t, options.Force)
	assert.NotNil(t, options.Auth)
}

func TestCreatePushOptions_Auth_SshKeyAuth_Null(t *testing.T) {
	ctx := context.TODO()
	model := &provider.PushResourceModel{}
	diagnostics := &diag.Diagnostics{}

	model.RefSpecs, _ = types.ListValueFrom(ctx, types.StringType, []string{"refs/heads/main"})
	model.Auth = types.ObjectValueMust(
		map[string]attr.Type{
			"ssh_key": types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"username":         types.StringType,
					"password":         types.StringType,
					"private_key_path": types.StringType,
					"private_key_pem":  types.StringType,
					"known_hosts":      types.ListType{ElemType: types.StringType},
				},
			},
		},
		map[string]attr.Value{
			"ssh_key": types.ObjectNull(map[string]attr.Type{
				"username":         types.StringType,
				"password":         types.StringType,
				"private_key_path": types.StringType,
				"private_key_pem":  types.StringType,
				"known_hosts":      types.ListType{ElemType: types.StringType},
			}),
		},
	)
	options := provider.CreatePushOptions(ctx, model, "", diagnostics)

	assert.NotNil(t, options)
	assert.False(t, diagnostics.HasError())
	assert.Equal(t, 1, len(options.RefSpecs))
	assert.Equal(t, "refs/heads/main", options.RefSpecs[0].String())
	assert.Empty(t, options.RemoteName)
	assert.False(t, options.Prune)
	assert.False(t, options.Force)
	assert.Nil(t, options.Auth)
}

func TestCreatePushOptions_Auth_SshKeyAuth_NoPrivateKey(t *testing.T) {
	ctx := context.TODO()
	model := &provider.PushResourceModel{}
	diagnostics := &diag.Diagnostics{}

	model.RefSpecs, _ = types.ListValueFrom(ctx, types.StringType, []string{"refs/heads/main"})
	model.Auth = types.ObjectValueMust(
		map[string]attr.Type{
			"ssh_key": types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"username":         types.StringType,
					"password":         types.StringType,
					"private_key_path": types.StringType,
					"private_key_pem":  types.StringType,
					"known_hosts":      types.ListType{ElemType: types.StringType},
				},
			},
		},
		map[string]attr.Value{
			"ssh_key": types.ObjectValueMust(
				map[string]attr.Type{
					"username":         types.StringType,
					"password":         types.StringType,
					"private_key_path": types.StringType,
					"private_key_pem":  types.StringType,
					"known_hosts":      types.ListType{ElemType: types.StringType},
				},
				map[string]attr.Value{
					"username":         types.StringValue("git"),
					"password":         types.StringValue(""),
					"private_key_path": types.StringNull(),
					"private_key_pem":  types.StringNull(),
					"known_hosts":      types.ListNull(types.StringType),
				},
			),
		},
	)
	options := provider.CreatePushOptions(ctx, model, "", diagnostics)

	assert.Nil(t, options)
	assert.True(t, diagnostics.HasError())
}

func TestCreatePushOptions_Auth_SshKeyAuth_PrivateKeyPath_Invalid(t *testing.T) {
	ctx := context.TODO()
	model := &provider.PushResourceModel{}
	diagnostics := &diag.Diagnostics{}

	model.RefSpecs, _ = types.ListValueFrom(ctx, types.StringType, []string{"refs/heads/main"})
	model.Auth = types.ObjectValueMust(
		map[string]attr.Type{
			"ssh_key": types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"username":         types.StringType,
					"password":         types.StringType,
					"private_key_path": types.StringType,
					"private_key_pem":  types.StringType,
					"known_hosts":      types.ListType{ElemType: types.StringType},
				},
			},
		},
		map[string]attr.Value{
			"ssh_key": types.ObjectValueMust(
				map[string]attr.Type{
					"username":         types.StringType,
					"password":         types.StringType,
					"private_key_path": types.StringType,
					"private_key_pem":  types.StringType,
					"known_hosts":      types.ListType{ElemType: types.StringType},
				},
				map[string]attr.Value{
					"username":         types.StringValue("git"),
					"password":         types.StringValue(""),
					"private_key_path": types.StringValue("~/.ssh/unknown_key_here"),
					"private_key_pem":  types.StringNull(),
					"known_hosts":      types.ListNull(types.StringType),
				},
			),
		},
	)
	options := provider.CreatePushOptions(ctx, model, "", diagnostics)

	assert.Nil(t, options)
	assert.True(t, diagnostics.HasError())
}

func TestCreatePushOptions_Auth_SshAgentAuth_Null(t *testing.T) {
	ctx := context.TODO()
	model := &provider.PushResourceModel{}
	diagnostics := &diag.Diagnostics{}

	model.RefSpecs, _ = types.ListValueFrom(ctx, types.StringType, []string{"refs/heads/main"})
	model.Auth = types.ObjectValueMust(
		map[string]attr.Type{
			"ssh_agent": types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"username":    types.StringType,
					"known_hosts": types.ListType{ElemType: types.StringType},
				},
			},
		},
		map[string]attr.Value{
			"ssh_agent": types.ObjectNull(map[string]attr.Type{
				"username":    types.StringType,
				"known_hosts": types.ListType{ElemType: types.StringType},
			}),
		},
	)
	options := provider.CreatePushOptions(ctx, model, "", diagnostics)

	assert.NotNil(t, options)
	assert.False(t, diagnostics.HasError())
	assert.Equal(t, 1, len(options.RefSpecs))
	assert.Equal(t, "refs/heads/main", options.RefSpecs[0].String())
	assert.Empty(t, options.RemoteName)
	assert.False(t, options.Prune)
	assert.False(t, options.Force)
	assert.Nil(t, options.Auth)
}

func TestCreatePushOptions_Auth_SshPasswordAuth_Null(t *testing.T) {
	ctx := context.TODO()
	model := &provider.PushResourceModel{}
	diagnostics := &diag.Diagnostics{}

	model.RefSpecs, _ = types.ListValueFrom(ctx, types.StringType, []string{"refs/heads/main"})
	model.Auth = types.ObjectValueMust(
		map[string]attr.Type{
			"ssh_password": types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"username":    types.StringType,
					"password":    types.StringType,
					"known_hosts": types.ListType{ElemType: types.StringType},
				},
			},
		},
		map[string]attr.Value{
			"ssh_password": types.ObjectNull(map[string]attr.Type{
				"username":    types.StringType,
				"password":    types.StringType,
				"known_hosts": types.ListType{ElemType: types.StringType},
			}),
		},
	)
	options := provider.CreatePushOptions(ctx, model, "", diagnostics)

	assert.NotNil(t, options)
	assert.False(t, diagnostics.HasError())
	assert.Equal(t, 1, len(options.RefSpecs))
	assert.Equal(t, "refs/heads/main", options.RefSpecs[0].String())
	assert.Empty(t, options.RemoteName)
	assert.False(t, options.Prune)
	assert.False(t, options.Force)
	assert.Nil(t, options.Auth)
}

func TestCreatePushOptions_Auth_SshPasswordAuth_Valid(t *testing.T) {
	ctx := context.TODO()
	model := &provider.PushResourceModel{}
	diagnostics := &diag.Diagnostics{}

	model.RefSpecs, _ = types.ListValueFrom(ctx, types.StringType, []string{"refs/heads/main"})
	model.Auth = types.ObjectValueMust(
		map[string]attr.Type{
			"ssh_password": types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"username":    types.StringType,
					"password":    types.StringType,
					"known_hosts": types.ListType{ElemType: types.StringType},
				},
			},
		},
		map[string]attr.Value{
			"ssh_password": types.ObjectValueMust(
				map[string]attr.Type{
					"username":    types.StringType,
					"password":    types.StringType,
					"known_hosts": types.ListType{ElemType: types.StringType},
				},
				map[string]attr.Value{
					"username":    types.StringValue("user"),
					"password":    types.StringValue("secret"),
					"known_hosts": types.ListNull(types.StringType),
				},
			),
		},
	)
	options := provider.CreatePushOptions(ctx, model, "", diagnostics)

	assert.NotNil(t, options)
	assert.False(t, diagnostics.HasError())
	assert.Equal(t, 1, len(options.RefSpecs))
	assert.Equal(t, "refs/heads/main", options.RefSpecs[0].String())
	assert.Empty(t, options.RemoteName)
	assert.False(t, options.Prune)
	assert.False(t, options.Force)
	assert.NotNil(t, options.Auth)
}

func TestCreatePushOptions_Auth_SshPasswordAuth_KnownHosts_Valid(t *testing.T) {
	ctx := context.TODO()
	model := &provider.PushResourceModel{}
	diagnostics := &diag.Diagnostics{}

	model.RefSpecs, _ = types.ListValueFrom(ctx, types.StringType, []string{"refs/heads/main"})
	list, _ := types.ListValueFrom(ctx, types.StringType, []string{"github.com 123abc"})
	model.Auth = types.ObjectValueMust(
		map[string]attr.Type{
			"ssh_password": types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"username":    types.StringType,
					"password":    types.StringType,
					"known_hosts": types.ListType{ElemType: types.StringType},
				},
			},
		},
		map[string]attr.Value{
			"ssh_password": types.ObjectValueMust(
				map[string]attr.Type{
					"username":    types.StringType,
					"password":    types.StringType,
					"known_hosts": types.ListType{ElemType: types.StringType},
				},
				map[string]attr.Value{
					"username":    types.StringValue("user"),
					"password":    types.StringValue("secret"),
					"known_hosts": list,
				},
			),
		},
	)
	options := provider.CreatePushOptions(ctx, model, "", diagnostics)

	assert.NotNil(t, options)
	assert.False(t, diagnostics.HasError())
	assert.Equal(t, 1, len(options.RefSpecs))
	assert.Equal(t, "refs/heads/main", options.RefSpecs[0].String())
	assert.Empty(t, options.RemoteName)
	assert.False(t, options.Prune)
	assert.False(t, options.Force)
	assert.NotNil(t, options.Auth)
}

func TestCreatePushOptions_Auth_SshPasswordAuth_KnownHosts_Null(t *testing.T) {
	ctx := context.TODO()
	model := &provider.PushResourceModel{}
	diagnostics := &diag.Diagnostics{}

	model.RefSpecs, _ = types.ListValueFrom(ctx, types.StringType, []string{"refs/heads/main"})
	model.Auth = types.ObjectValueMust(
		map[string]attr.Type{
			"ssh_password": types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"username":    types.StringType,
					"password":    types.StringType,
					"known_hosts": types.ListType{ElemType: types.StringType},
				},
			},
		},
		map[string]attr.Value{
			"ssh_password": types.ObjectValueMust(
				map[string]attr.Type{
					"username":    types.StringType,
					"password":    types.StringType,
					"known_hosts": types.ListType{ElemType: types.StringType},
				},
				map[string]attr.Value{
					"username":    types.StringValue("user"),
					"password":    types.StringValue("secret"),
					"known_hosts": types.ListNull(types.StringType),
				},
			),
		},
	)
	options := provider.CreatePushOptions(ctx, model, "", diagnostics)

	assert.NotNil(t, options)
	assert.False(t, diagnostics.HasError())
	assert.Equal(t, 1, len(options.RefSpecs))
	assert.Equal(t, "refs/heads/main", options.RefSpecs[0].String())
	assert.Empty(t, options.RemoteName)
	assert.False(t, options.Prune)
	assert.False(t, options.Force)
	assert.NotNil(t, options.Auth)
}

func TestCreatePushOptions_Auth_SshPasswordAuth_KnownHosts_File(t *testing.T) {
	ctx := context.TODO()
	model := &provider.PushResourceModel{}
	diagnostics := &diag.Diagnostics{}

	hostKey := hostPublicKey(t)
	model.RefSpecs, _ = types.ListValueFrom(ctx, types.StringType, []string{"refs/heads/main"})
	model.Auth = sshPasswordAuthWithKnownHosts(t, knownHostsFile(t, knownhosts.Line([]string{"example.com"}, hostKey)))
	options := provider.CreatePushOptions(ctx, model, "git@example.com:metio/terraform-provider-git.git", diagnostics)

	assert.NotNil(t, options)
	assert.False(t, diagnostics.HasError())
	auth, ok := options.Auth.(*gitssh.Password)
	assert.True(t, ok)
	assert.Equal(t, []string{ssh.KeyAlgoED25519}, auth.HostKeyAlgorithms)
	remote := &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 22}
	assert.NoError(t, auth.HostKeyCallback("example.com:22", remote, hostKey))
	assert.Error(t, auth.HostKeyCallback("example.com:22", remote, hostPublicKey(t)))
}

func TestCreatePushOptions_Auth_SshPasswordAuth_KnownHosts_File_CustomPort(t *testing.T) {
	ctx := context.TODO()
	model := &provider.PushResourceModel{}
	diagnostics := &diag.Diagnostics{}

	hostKey := hostPublicKey(t)
	model.RefSpecs, _ = types.ListValueFrom(ctx, types.StringType, []string{"refs/heads/main"})
	model.Auth = sshPasswordAuthWithKnownHosts(t, knownHostsFile(t, knownhosts.Line([]string{"[example.com]:2222"}, hostKey)))
	options := provider.CreatePushOptions(ctx, model, "ssh://git@example.com:2222/metio/terraform-provider-git.git", diagnostics)

	assert.NotNil(t, options)
	assert.False(t, diagnostics.HasError())
	auth, ok := options.Auth.(*gitssh.Password)
	assert.True(t, ok)
	assert.Equal(t, []string{ssh.KeyAlgoED25519}, auth.HostKeyAlgorithms)
}

func TestCreatePushOptions_Auth_SshPasswordAuth_KnownHosts_File_Missing(t *testing.T) {
	ctx := context.TODO()
	model := &provider.PushResourceModel{}
	diagnostics := &diag.Diagnostics{}

	model.RefSpecs, _ = types.ListValueFrom(ctx, types.StringType, []string{"refs/heads/main"})
	model.Auth = sshPasswordAuthWithKnownHosts(t, filepath.Join(t.TempDir(), "does-not-exist"))
	options := provider.CreatePushOptions(ctx, model, "git@example.com:metio/terraform-provider-git.git", diagnostics)

	assert.Nil(t, options)
	assert.True(t, diagnostics.HasError())
}

func hostPublicKey(t *testing.T) ssh.PublicKey {
	publicKey, _, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	sshPublicKey, err := ssh.NewPublicKey(publicKey)
	if err != nil {
		t.Fatal(err)
	}
	return sshPublicKey
}

func knownHostsFile(t *testing.T, line string) string {
	file := filepath.Join(t.TempDir(), "known_hosts")
	if err := os.WriteFile(file, []byte(line+"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	return file
}

func sshPasswordAuthWithKnownHosts(t *testing.T, file string) types.Object {
	knownHosts, diags := types.SetValueFrom(context.TODO(), types.StringType, []string{file})
	if diags.HasError() {
		t.Fatal(diags)
	}
	attributeTypes := map[string]attr.Type{
		"username":    types.StringType,
		"password":    types.StringType,
		"known_hosts": types.SetType{ElemType: types.StringType},
	}
	return types.ObjectValueMust(
		map[string]attr.Type{
			"ssh_password": types.ObjectType{AttrTypes: attributeTypes},
		},
		map[string]attr.Value{
			"ssh_password": types.ObjectValueMust(
				attributeTypes,
				map[string]attr.Value{
					"username":    types.StringValue("user"),
					"password":    types.StringValue("secret"),
					"known_hosts": knownHosts,
				},
			),
		},
	)
}
