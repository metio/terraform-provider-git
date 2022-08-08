/*
 * This file is part of terraform-provider-git. It is subject to the license terms in the LICENSE file found in the top-level
 * directory of this distribution and at https://creativecommons.org/publicdomain/zero/1.0/. No part of terraform-provider-git,
 * including this file, may be copied, modified, propagated, or distributed except according to the terms contained
 * in the LICENSE file.
 */

package provider

import "github.com/hashicorp/terraform-plugin-framework/types"

type GitStatus struct {
	Staging  types.String `tfsdk:"staging"`
	Worktree types.String `tfsdk:"worktree"`
}
