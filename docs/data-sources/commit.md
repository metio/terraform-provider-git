---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "git_commit Data Source - terraform-provider-git"
subcategory: ""
description: |-
  Fetches information about a single commit.
---

# git_commit (Data Source)

Fetches information about a single commit.

## Example Usage

```terraform
# get commit by full sha1
data "git_commit" "full_sha1" {
  directory = "/path/to/git/repository"
  revision  = "dae86e1950b1277e545cee180551750029cfe735"
}

# get commit by short sha1
data "git_commit" "short_sha1" {
  directory = "/path/to/git/repository"
  revision  = "dae86e"
}

# get commit by refname
data "git_commit" "refname" {
  directory = "/path/to/git/repository"
  revision  = "main"
}

# get commit by head shortcut
data "git_commit" "head_shortcut" {
  directory = "/path/to/git/repository"
  revision  = "@"
}

# get commit by parent
data "git_commit" "head_parent" {
  directory = "/path/to/git/repository"
  revision  = "HEAD~1"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `directory` (String) The path to the local Git repository.
- `revision` (String) The [revision](https://www.git-scm.com/docs/gitrevisions) of the commit to fetch. Note that `go-git` does not [support](https://pkg.go.dev/github.com/go-git/go-git/v5#Repository.ResolveRevision) every revision type at the moment.

### Read-Only

- `author` (Attributes) The original author of the commit. (see [below for nested schema](#nestedatt--author))
- `committer` (Attributes) The person performing the commit. (see [below for nested schema](#nestedatt--committer))
- `files` (List of String) The files updated by the commit.
- `id` (String) The same value as the `revision` attribute.
- `message` (String) The message of the commit.
- `sha1` (String) The SHA1 hash of the resolved revision.
- `signature` (String) The signature of the commit.
- `tree_sha1` (String) The SHA1 checksum of the root tree of the commit.

<a id="nestedatt--author"></a>
### Nested Schema for `author`

Read-Only:

- `email` (String) The email address of the author.
- `name` (String) The name of the author.
- `timestamp` (String) The timestamp of the signature.


<a id="nestedatt--committer"></a>
### Nested Schema for `committer`

Read-Only:

- `email` (String) The email address of the committer.
- `name` (String) The name of the committer.
- `timestamp` (String) The timestamp of the signature.
