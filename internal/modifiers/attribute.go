/*
 * SPDX-FileCopyrightText: Sebastian Hoß <seb@hoß.de>
 * SPDX-License-Identifier: 0BSD
 */

package modifiers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// DefaultValue accepts an attr.Value and uses the supplied value to set a default if the config for
// the attribute is null.
func DefaultValue(val attr.Value) tfsdk.AttributePlanModifier {
	return &defaultValueAttributePlanModifier{val}
}

type defaultValueAttributePlanModifier struct {
	defaultValue attr.Value
}

func (d *defaultValueAttributePlanModifier) Description(ctx context.Context) string {
	return "If the config does not contain a value, a default will be set using defaultValue."
}

func (d *defaultValueAttributePlanModifier) MarkdownDescription(ctx context.Context) string {
	return d.Description(ctx)
}

// Modify checks that the value of the attribute in the configuration and assigns the default value if
// the value in the config is null. This is a destructive operation in that it will overwrite any value
// present in the plan.
func (d *defaultValueAttributePlanModifier) Modify(ctx context.Context, req tfsdk.ModifyAttributePlanRequest, resp *tfsdk.ModifyAttributePlanResponse) {
	// If the attribute configuration is not null, we are done here
	if !req.AttributeConfig.IsNull() {
		return
	}

	resp.AttributePlan = d.defaultValue
}
