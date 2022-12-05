/*
 * SPDX-FileCopyrightText: The terraform-provider-git Authors
 * SPDX-License-Identifier: 0BSD
 */

package modifiers

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// DefaultBool accepts a bool and uses the supplied value to set a default if the config for
// the attribute is null.
func DefaultBool(val bool) planmodifier.Bool {
	return &defaultBoolPlanModifier{types.BoolValue(val)}
}

type defaultBoolPlanModifier struct {
	defaultValue types.Bool
}

func (d *defaultBoolPlanModifier) Description(_ context.Context) string {
	return "If the config does not contain a value, a default will be set using defaultValue."
}

func (d *defaultBoolPlanModifier) MarkdownDescription(ctx context.Context) string {
	return d.Description(ctx)
}

// PlanModifyBool checks that the value of the attribute in the configuration and assigns the default value if
// the value in the config is null. This is a destructive operation in that it will overwrite any value
// present in the plan.
func (d *defaultBoolPlanModifier) PlanModifyBool(_ context.Context, req planmodifier.BoolRequest, resp *planmodifier.BoolResponse) {
	if req.ConfigValue.IsNull() {
		resp.PlanValue = d.defaultValue
	}
}
