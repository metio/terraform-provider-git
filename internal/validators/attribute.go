/*
 * This file is part of terraform-provider-git. It is subject to the license terms in the LICENSE file found in the top-level
 * directory of this distribution and at https://creativecommons.org/publicdomain/zero/1.0/. No part of terraform-provider-git,
 * including this file, may be copied, modified, propagated, or distributed except according to the terms contained
 * in the LICENSE file.
 */

package validators

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// oneOfAttributeValidator checks that values held in the attribute
// are one of the given `acceptableValues`.
//
// This validator can be used with all primitive `types.*`, as well as
// collections (`types.List`, `types.Set`, `types.Map` and `types.Object`):
// for key/value collections, the validator will compare the values.
//
// Instances should be created via OneOf function.
type oneOfAttributeValidator struct {
	acceptableValues []attr.Value
}

// OneOf is a helper to instantiate a oneOfAttributeValidator.
func OneOf(acceptableValues ...attr.Value) tfsdk.AttributeValidator {
	return &oneOfAttributeValidator{acceptableValues}
}

var _ tfsdk.AttributeValidator = (*oneOfAttributeValidator)(nil)

func (av *oneOfAttributeValidator) Description(ctx context.Context) string {
	return av.MarkdownDescription(ctx)
}

func (av *oneOfAttributeValidator) MarkdownDescription(_ context.Context) string {
	return fmt.Sprintf("Ensures that the attribute is one of: %q", av.acceptableValues)
}

func (av *oneOfAttributeValidator) Validate(_ context.Context, req tfsdk.ValidateAttributeRequest, res *tfsdk.ValidateAttributeResponse) {
	if req.AttributeConfig.IsNull() || req.AttributeConfig.IsUnknown() {
		return
	}

	var values []attr.Value

	switch typedAttributeConfig := req.AttributeConfig.(type) {
	case types.List:
		values = typedAttributeConfig.Elems
	case types.Map:
		values = make([]attr.Value, 0, len(typedAttributeConfig.Elems))
		for _, v := range typedAttributeConfig.Elems {
			values = append(values, v)
		}
	case types.Set:
		values = typedAttributeConfig.Elems
	case types.Object:
		values = make([]attr.Value, 0, len(typedAttributeConfig.Attrs))
		for _, v := range typedAttributeConfig.Attrs {
			values = append(values, v)
		}
	default:
		values = []attr.Value{typedAttributeConfig}
	}

	for _, v := range values {
		if !av.isValid(v) {
			res.Diagnostics.AddAttributeError(
				req.AttributePath,
				"Invalid Attribute",
				fmt.Sprintf("Value %q must be one of %q", v, av.acceptableValues),
			)
		}
	}
}

func (av *oneOfAttributeValidator) isValid(v attr.Value) bool {
	for _, acceptableV := range av.acceptableValues {
		if v.Equal(acceptableV) {
			return true
		}
	}

	return false
}
