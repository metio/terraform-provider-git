# SPDX-FileCopyrightText: The terraform-provider-git Authors
# SPDX-License-Identifier: 0BSD

output "directory" {
  value = data.git_log.current_log.directory
}

output "id" {
  value = data.git_log.current_log.id
}

output "commits" {
  value = data.git_log.current_log.commits
}
