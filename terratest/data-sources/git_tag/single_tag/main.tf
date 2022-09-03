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

data "git_tag" "single_tag" {
  directory = var.directory
  name      = var.name
}
