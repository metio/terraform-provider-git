# SPDX-FileCopyrightText: The terraform-provider-git Authors
# SPDX-License-Identifier: 0BSD

output "directory" {
  value = data.git_branch.current.directory
}

output "name" {
  value = data.git_branch.current.name
}

output "id" {
  value = data.git_branch.current.id
}

output "sha1" {
  value = data.git_branch.current.sha1
}
