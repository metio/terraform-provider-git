# SPDX-FileCopyrightText: The terraform-provider-git Authors
# SPDX-License-Identifier: 0BSD

output "every_tag" {
  value = data.git_tag.every_tag
}
