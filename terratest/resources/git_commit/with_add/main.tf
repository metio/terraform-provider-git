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

resource "git_add" "add" {
  directory = var.directory
  add_paths = var.add_paths
}

resource "git_commit" "commit" {
  directory = git_add.add.directory
  message   = "committed with terraform"

  author = {
    name  = "terraform"
    email = "automation@example.com"
  }

  lifecycle {
    replace_triggered_by = [git_add.add.id]
  }
}

data "git_commit" "commit" {
  directory = git_commit.commit.directory
  revision  = git_commit.commit.sha1
}

data "git_status" "single_file" {
  directory = git_commit.commit.directory
  file      = var.file
}
