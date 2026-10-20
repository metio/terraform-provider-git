/*
 * SPDX-FileCopyrightText: The terraform-provider-git Authors
 * SPDX-License-Identifier: 0BSD
 */

package provider

import (
	"context"
	"errors"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func getRemote(ctx context.Context, repository *git.Repository, remoteName string, diag *diag.Diagnostics) *git.Remote {
	remote, err := repository.Remote(remoteName)
	if err != nil {
		diag.AddError(
			"Cannot read remote",
			"Could not read remote ["+remoteName+"] because of: "+err.Error(),
		)
		return nil
	}
	tflog.Trace(ctx, "read remote", map[string]interface{}{
		"remote": remoteName,
	})
	return remote
}

// getRemoteHeadBranch reads the branch that the local refs/remotes/<name>/HEAD reference points to.
func getRemoteHeadBranch(ctx context.Context, repository *git.Repository, remoteName string, diag *diag.Diagnostics) types.String {
	reference, err := repository.Reference(plumbing.NewRemoteHEADReferenceName(remoteName), false)
	if errors.Is(err, plumbing.ErrReferenceNotFound) {
		return types.StringNull()
	}
	if err != nil {
		diag.AddError(
			"Cannot read remote HEAD",
			"Could not read HEAD of remote ["+remoteName+"] because of: "+err.Error(),
		)
		return types.StringNull()
	}
	tflog.Trace(ctx, "read remote HEAD", map[string]interface{}{
		"remote":    remoteName,
		"reference": reference.String(),
	})

	// branch names can contain slashes, so strip the prefix instead of using Short()
	prefix := "refs/remotes/" + remoteName + "/"
	target := reference.Target().String()
	if reference.Type() != plumbing.SymbolicReference || !strings.HasPrefix(target, prefix) {
		return types.StringNull()
	}
	return types.StringValue(strings.TrimPrefix(target, prefix))
}

// queryRemoteHeadBranch asks the remote which branch its HEAD points to, similar to 'git remote show'.
func queryRemoteHeadBranch(ctx context.Context, remote *git.Remote, auth types.Object, diag *diag.Diagnostics) types.String {
	remoteName := remote.Config().Name

	authMethod := authOptions(ctx, auth, diag)
	if diag.HasError() {
		return types.StringNull()
	}

	references, err := remote.List(&git.ListOptions{
		Auth: authMethod,
	})
	if errors.Is(err, transport.ErrEmptyRemoteRepository) {
		return types.StringNull()
	}
	if err != nil {
		diag.AddError(
			"Cannot list remote",
			"Could not list references of remote ["+remoteName+"] because of: "+err.Error(),
		)
		return types.StringNull()
	}
	tflog.Trace(ctx, "listed remote references", map[string]interface{}{
		"remote":     remoteName,
		"references": len(references),
	})

	for _, reference := range references {
		if reference.Name() == plumbing.HEAD && reference.Type() == plumbing.SymbolicReference && reference.Target().IsBranch() {
			return types.StringValue(strings.TrimPrefix(reference.Target().String(), "refs/heads/"))
		}
	}
	return types.StringNull()
}
