/*
 * This file is part of terraform-provider-git. It is subject to the license terms in the LICENSE file found in the top-level
 * directory of this distribution and at https://creativecommons.org/publicdomain/zero/1.0/. No part of terraform-provider-git,
 * including this file, may be copied, modified, propagated, or distributed except according to the terms contained
 * in the LICENSE file.
 */

package provider

import (
	"context"
	"github.com/go-git/go-git/v5"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type GitRemote struct {
	URLs types.List `tfsdk:"urls"`
}

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
