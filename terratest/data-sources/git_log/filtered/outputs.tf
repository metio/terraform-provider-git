# SPDX-FileCopyrightText: The terraform-provider-git Authors
# SPDX-License-Identifier: 0BSD

output "directory" {
  value = data.git_log.filtered.directory
}

output "id" {
  value = data.git_log.filtered.id
}

output "commits" {
  value = data.git_log.filtered.commits
}
