# SPDX-FileCopyrightText: The terraform-provider-git Authors
# SPDX-License-Identifier: 0BSD

terraform {
  required_providers {
    git = {
      source  = "localhost/metio/git"
      version = "9999.99.99"
    }
  }
}

provider "git" {
  # Configuration options
}

data "git_branches" "all" {
  directory = var.directory
}

data "git_branch" "every" {
  for_each  = data.git_branches.all.branches
  directory = var.directory
  name      = each.key
}
