resource "git_tag" "tag" {
  directory = "/path/to/git/repository"
  name      = "v1.2.3"
}

resource "git_tag" "annotated_tag" {
  directory = "/path/to/git/repository"
  name      = "v1.2.3"
  message   = "some message for the new tag"
}

resource "git_tag" "specific_commit" {
  directory = "/path/to/git/repository"
  name      = "v1.2.3"
  revision  = "b1af8d13f5131c9b4de9ddd06e311c2e79fdb285"
}

resource "git_tag" "head" {
  directory = "/path/to/git/repository"
  name      = "v1.2.3"
  revision  = "HEAD"
}

resource "git_tag" "branch" {
  directory = "/path/to/git/repository"
  name      = "v1.2.3"
  revision  = "main"
}
