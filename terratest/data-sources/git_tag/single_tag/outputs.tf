# SPDX-FileCopyrightText: The terraform-provider-git Authors
# SPDX-License-Identifier: 0BSD

output "directory" {
  value = data.git_tag.single_tag.directory
}

output "id" {
  value = data.git_tag.single_tag.id
}

output "name" {
  value = data.git_tag.single_tag.name
}
