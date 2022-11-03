/*
 * SPDX-FileCopyrightText: The terraform-provider-git Authors
 * SPDX-License-Identifier: 0BSD
 */

package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func updatedUsingPlan(ctx context.Context, req *resource.UpdateRequest, res *resource.UpdateResponse, model interface{}) {
	res.Diagnostics.Append(req.Plan.Get(ctx, model)...)
	if res.Diagnostics.HasError() {
		return
	}
	res.Diagnostics.Append(res.State.Set(ctx, model)...)
}
