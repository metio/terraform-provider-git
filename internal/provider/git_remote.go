/*
 * SPDX-FileCopyrightText: Sebastian Hoß <seb@hoß.de>
 * SPDX-License-Identifier: BSD0
 */

package provider

import (
	"context"
	"github.com/go-git/go-git/v5"
	"github.com/hashicorp/terraform-plugin-framework/diag"
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
