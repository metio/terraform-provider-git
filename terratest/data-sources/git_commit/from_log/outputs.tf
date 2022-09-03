# SPDX-FileCopyrightText: The terraform-provider-git Authors
# SPDX-License-Identifier: 0BSD

output "from_log" {
  value = data.git_commit.from_log
}
