# SPDX-FileCopyrightText: The terraform-provider-git Authors
# SPDX-License-Identifier: 0BSD

output "directory" {
  value = git_add.add.directory
}

output "id" {
  value = git_add.add.id
}

output "add_paths" {
  value = git_add.add.add_paths
}

output "files" {
  value = data.git_statuses.multiple_files.files
}
