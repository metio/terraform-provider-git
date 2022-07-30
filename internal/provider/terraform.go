/*
 * This file is part of terraform-provider-git. It is subject to the license terms in the LICENSE file found in the top-level
 * directory of this distribution and at https://creativecommons.org/publicdomain/zero/1.0/. No part of terraform-provider-git,
 * including this file, may be copied, modified, propagated, or distributed except according to the terms contained
 * in the LICENSE file.
 */

package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func updatedUsingPlan(ctx context.Context, req *tfsdk.UpdateResourceRequest, res *tfsdk.UpdateResourceResponse, model interface{}) {
	// Read the plan
	res.Diagnostics.Append(req.Plan.Get(ctx, model)...)
	if res.Diagnostics.HasError() {
		return
	}

	// Set it as the new state
	res.Diagnostics.Append(res.State.Set(ctx, model)...)
}

func stringsToList(strings []string) types.List {
	var values []attr.Value
	for _, url := range strings {
		values = append(values, types.String{Value: url})
	}
	return types.List{
		Elems:    values,
		ElemType: types.StringType,
	}
}
