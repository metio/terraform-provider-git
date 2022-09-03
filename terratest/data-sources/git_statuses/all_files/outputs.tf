# SPDX-FileCopyrightText: The terraform-provider-git Authors
# SPDX-License-Identifier: 0BSD

output "is_clean" {
  value = data.git_statuses.all_files.is_clean
}

output "files" {
  value = data.git_statuses.all_files.files
}
