//go:build testing

/*
 * SPDX-FileCopyrightText: The terraform-provider-git Authors
 * SPDX-License-Identifier: 0BSD
 */

package testutils

import (
	"fmt"
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
