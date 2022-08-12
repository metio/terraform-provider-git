/*
 * SPDX-FileCopyrightText: Sebastian Hoß <seb@hoß.de>
 * SPDX-License-Identifier: 0BSD
 */

package provider

import "github.com/go-git/go-git/v5/config"

func mapConfigScope(userInput string) config.Scope {
	switch userInput {
	case "local":
		return config.LocalScope
	case "system":
		return config.SystemScope
	default:
		return config.GlobalScope
	}
}
