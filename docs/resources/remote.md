---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "git_remote Resource - terraform-provider-git"
subcategory: ""
description: |-
  Manages remotes in a Git repository similar to git remote.
---

# git_remote (Resource)

Manages remotes in a Git repository similar to `git remote`.

## Example Usage

```terraform
resource "git_remote" "remote" {
  directory = "/path/to/git/repository"
  name      = "some-remote"
  urls      = ["https://github.com/some-org/some-repo.git"]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `directory` (String) The path to the local Git repository.
- `name` (String) The name of the Git remote to manage.
- `urls` (List of String) The URLs of the Git remote to manage. The first URL will be a fetch/pull URL. All other URLs will be push only.

### Read-Only

- `id` (String) The import ID to import this resource which has the form `'directory|name'`

## Import

Import is supported using the following syntax:

```shell
# git_remote resources can be imported by specifying the directory of the
# Git repository and the name of the remote to import. Both values are
# separated by a single '|'.
terraform import git_remote.remote 'path/to/your/git/repository|name-of-your-remote'
```
