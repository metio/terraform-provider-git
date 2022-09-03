# SPDX-FileCopyrightText: The terraform-provider-git Authors
# SPDX-License-Identifier: 0BSD

output "directory" {
  value = data.git_remotes.all_remotes.directory
}

output "id" {
  value = data.git_remotes.all_remotes.id
}

output "remotes" {
  value = data.git_remotes.all_remotes.remotes
}
