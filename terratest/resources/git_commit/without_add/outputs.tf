# SPDX-FileCopyrightText: The terraform-provider-git Authors
# SPDX-License-Identifier: 0BSD

output "sha1" {
  value = data.git_commit.commit.sha1
}

output "files" {
  value = data.git_commit.commit.files
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
