//go:build testing

/*
 * SPDX-FileCopyrightText: The terraform-provider-git Authors
 * SPDX-License-Identifier: 0BSD
 */

package testutils

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
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
