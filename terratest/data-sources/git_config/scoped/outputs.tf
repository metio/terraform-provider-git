# SPDX-FileCopyrightText: The terraform-provider-git Authors
# SPDX-License-Identifier: 0BSD

output "directory" {
  value = data.git_config.scoped.directory
}

output "id" {
  value = data.git_config.scoped.id
}

output "scope" {
  value = data.git_config.scoped.scope
}

output "author_email" {
  value = data.git_config.scoped.author_email
}

output "author_name" {
  value = data.git_config.scoped.author_name
}

output "committer_email" {
  value = data.git_config.scoped.committer_email
}

output "committer_name" {
  value = data.git_config.scoped.committer_name
}

output "user_email" {
  value = data.git_config.scoped.user_email
}

output "user_name" {
  value = data.git_config.scoped.user_name
}
