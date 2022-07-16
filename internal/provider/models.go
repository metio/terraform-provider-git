/*
 * This file is part of terraform-gitProvider-git. It is subject to the license terms in the LICENSE file found in the top-level
 * directory of this distribution and at https://creativecommons.org/publicdomain/zero/1.0/. No part of terraform-gitProvider-git,
 * including this file, may be copied, modified, propagated, or distributed except according to the terms contained
 * in the LICENSE file.
 */

package provider

import "github.com/hashicorp/terraform-plugin-framework/types"

type GitTag struct {
	Annotated   types.Bool   `tfsdk:"annotated"`
	Lightweight types.Bool   `tfsdk:"lightweight"`
	SHA1        types.String `tfsdk:"sha1"`
}

type GitBranch struct {
	SHA1   types.String `tfsdk:"sha1"`
	Remote types.String `tfsdk:"remote"`
	Rebase types.String `tfsdk:"rebase"`
}

type GitRemote struct {
	URLs []types.String `tfsdk:"urls"`
}
