# SPDX-FileCopyrightText: The terraform-provider-git Authors
# SPDX-License-Identifier: 0BSD

output "directory" {
  value = data.git_tags.all_tags.directory
}

output "id" {
  value = data.git_tags.all_tags.id
}

output "tags" {
  value = data.git_tags.all_tags.tags
}

output "annotated" {
  value = data.git_tags.all_tags.annotated
}

output "lightweight" {
  value = data.git_tags.all_tags.lightweight
}
