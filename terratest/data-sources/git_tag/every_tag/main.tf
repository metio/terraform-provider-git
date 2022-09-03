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

data "git_tags" "all_tags" {
  directory = var.directory
}

data "git_tag" "every_tag" {
  for_each  = data.git_tags.all_tags.tags
  directory = var.directory
  name      = each.key
}
