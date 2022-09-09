# SPDX-FileCopyrightText: The terraform-provider-git Authors
# SPDX-License-Identifier: 0BSD

variable "directory" {
  type = string
}

variable "add_paths" {
  type = list(string)
}
