# SPDX-FileCopyrightText: The terraform-provider-git Authors
# SPDX-License-Identifier: 0BSD

output "directory" {
  value = data.git_commit.current_head.directory
}

output "id" {
  value = data.git_commit.current_head.id
}

output "revision" {
  value = data.git_commit.current_head.revision
}

output "message" {
  value = data.git_commit.current_head.message
}

output "sha1" {
  value = data.git_commit.current_head.sha1
}

output "tree_sha1" {
  value = data.git_commit.current_head.tree_sha1
}

output "author" {
  value = data.git_commit.current_head.author
}

output "committer" {
  value = data.git_commit.current_head.committer
}
