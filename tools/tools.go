//go:build tools
// +build tools

/*
 * SPDX-FileCopyrightText: Sebastian Hoß <seb@hoß.de>
 * SPDX-License-Identifier: BSD0
 */

package tools

import (
	// document generation
	_ "github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs"
)
