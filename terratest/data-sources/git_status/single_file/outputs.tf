# SPDX-FileCopyrightText: The terraform-provider-git Authors
# SPDX-License-Identifier: 0BSD

output "directory" {
  value = data.git_status.single_file.directory
}

output "id" {
  value = data.git_status.single_file.id
}

output "file" {
  value = data.git_status.single_file.file
}

output "staging" {
  value = data.git_status.single_file.staging
}

output "worktree" {
  value = data.git_status.single_file.worktree
}
