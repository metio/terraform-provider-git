---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "git_tag Data Source - terraform-provider-git"
subcategory: ""
description: |-
  Reads information about a specific tag of a Git repository.
---

# git_tag (Data Source)

Reads information about a specific tag of a Git repository.

## Example Usage

```terraform
data "git_tag" "tag" {
  directory = "/path/to/git/repository"
  name      = "v1.2.3"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `directory` (String) The path to the local Git repository.
- `name` (String) The name of the tag to gather information about.

### Read-Only

- `annotated` (Boolean) Whether the given tag is an annotated tag.
- `id` (String) The same value as the `name` attribute.
- `lightweight` (Boolean) Whether the given tag is a lightweight tag.
- `message` (String) The associated message of an annotated tag.
- `sha1` (String) The SHA1 checksum of the commit the given tag is pointing at.
