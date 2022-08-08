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
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func openRepository(ctx context.Context, directory string, diag *diag.Diagnostics) *git.Repository {
	repository, err := git.PlainOpen(directory)
	if err != nil {
		diag.AddError(
			"Cannot open repository",
			"Could not open git repository ["+directory+"] because of: "+err.Error(),
		)
		return nil
	}
	tflog.Trace(ctx, "opened repository", map[string]interface{}{
		"directory": directory,
	})
	return repository
}
