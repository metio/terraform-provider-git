# SPDX-FileCopyrightText: The terraform-provider-git Authors
# SPDX-License-Identifier: 0BSD

output "directory" {
  value = data.git_remote.single_remote.directory
}

output "id" {
  value = data.git_remote.single_remote.id
}

output "name" {
  value = data.git_remote.single_remote.name
}

output "urls" {
  value = data.git_remote.single_remote.urls
}
