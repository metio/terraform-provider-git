# SPDX-FileCopyrightText: The terraform-provider-git Authors
# SPDX-License-Identifier: 0BSD

output "directory" {
  value = data.git_repository.repository.directory
}

output "id" {
  value = data.git_repository.repository.id
}

output "sha1" {
  value = data.git_repository.repository.sha1
}
