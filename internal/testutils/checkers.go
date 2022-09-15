//go:build testing

/*
 * SPDX-FileCopyrightText: The terraform-provider-git Authors
 * SPDX-License-Identifier: 0BSD
 */

package testutils

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	"testing"
)

func CheckExactLength(expectedLength int) func(input string) error {
	return func(input string) error {
		if len(input) != expectedLength {
			return fmt.Errorf("expected length %d, actual length %d", expectedLength, len(input))
		}

		return nil
	}
}

func CheckMinLength(minimumLength int) func(input string) error {
	return func(input string) error {
		if len(input) < minimumLength {
			return fmt.Errorf("minimum length %d, actual length %d", minimumLength, len(input))
		}

		return nil
	}
}

func ComposeImportStateCheck(fs ...resource.ImportStateCheckFunc) resource.ImportStateCheckFunc {
	return func(s []*terraform.InstanceState) error {
		for i, f := range fs {
			if err := f(s); err != nil {
				return fmt.Errorf("check %d/%d error: %s", i+1, len(fs), err)
			}
		}

		return nil
	}
}

func CheckResourceAttrInstanceState(attributeName, attributeValue string) resource.ImportStateCheckFunc {
	return func(is []*terraform.InstanceState) error {
		if len(is) != 1 {
			return fmt.Errorf("unexpected number of instance states: %d", len(is))
		}

		s := is[0]

		attrVal, ok := s.Attributes[attributeName]
		if !ok {
			return fmt.Errorf("attribute '%s' not found in instance state", attributeName)
		}

		if attrVal != attributeValue {
			return fmt.Errorf("attribute '%s' expected: '%s', got: '%s'", attributeName, attributeValue, attrVal)
		}

		return nil
	}
}

func VerifySchemaDescriptions(t *testing.T, schema tfsdk.Schema) {
	assert.NotEmpty(t, schema.Description, "schema.Description")
	assert.NotEmpty(t, schema.MarkdownDescription, "schema.MarkdownDescription")

	for name, attribute := range schema.Attributes {
		assert.NotEmpty(t, attribute.Description, fmt.Sprintf("%s.Description", name))
		assert.NotEmpty(t, attribute.MarkdownDescription, fmt.Sprintf("%s.MarkdownDescription", name))

		if attribute.Attributes != nil {
			for nn, nv := range attribute.Attributes.GetAttributes() {
				assert.NotEmpty(t, nv.GetDescription(), fmt.Sprintf("%s.Description", nn))
				assert.NotEmpty(t, nv.GetMarkdownDescription(), fmt.Sprintf("%s.MarkdownDescription", nn))
			}
		}
	}
}
