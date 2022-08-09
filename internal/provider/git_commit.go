/*
 * This file is part of terraform-provider-git. It is subject to the license terms in the LICENSE file found in the top-level
 * directory of this distribution and at https://creativecommons.org/publicdomain/zero/1.0/. No part of terraform-provider-git,
 * including this file, may be copied, modified, propagated, or distributed except according to the terms contained
 * in the LICENSE file.
 */

package provider

import (
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func signatureAsMap(signature *object.Signature) types.Object {
	data := make(map[string]attr.Value)
	data["name"] = types.String{Value: signature.Name}
	data["email"] = types.String{Value: signature.Email}
	data["timestamp"] = types.String{Value: signature.When.String()}

	return types.Object{
		AttrTypes: map[string]attr.Type{
			"name":      types.StringType,
			"email":     types.StringType,
			"timestamp": types.StringType,
		},
		Attrs: data,
	}
}
