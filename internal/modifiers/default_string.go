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

// DefaultString accepts a string and uses the supplied value to set a default if the config for
// the attribute is null.
func DefaultString(val string) planmodifier.String {
	return &defaultStringPlanModifier{types.StringValue(val)}
}

type defaultStringPlanModifier struct {
	defaultValue types.String
}

func (d *defaultStringPlanModifier) Description(_ context.Context) string {
	return "If the config does not contain a value, a default will be set using defaultValue."
}

func (d *defaultStringPlanModifier) MarkdownDescription(ctx context.Context) string {
	return d.Description(ctx)
}

// PlanModifyString checks that the value of the attribute in the configuration and assigns the default value if
// the value in the config is null. This is a destructive operation in that it will overwrite any value
// present in the plan.
func (d *defaultStringPlanModifier) PlanModifyString(_ context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if req.ConfigValue.IsNull() {
		resp.PlanValue = d.defaultValue
	}
}
