# SPDX-FileCopyrightText: The terraform-provider-git Authors
# SPDX-License-Identifier: 0BSD

output "directory" {
  value = data.git_config.defaults.directory
}

output "id" {
  value = data.git_config.defaults.id
}

output "scope" {
  value = data.git_config.defaults.scope
}

output "author_email" {
  value = data.git_config.defaults.author_email
}

output "author_name" {
  value = data.git_config.defaults.author_name
}

output "committer_email" {
  value = data.git_config.defaults.committer_email
}

output "committer_name" {
  value = data.git_config.defaults.committer_name
}

output "user_email" {
  value = data.git_config.defaults.user_email
}

output "user_name" {
  value = data.git_config.defaults.user_name
}
