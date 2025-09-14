resource "git_clone" "clone" {
  directory      = "/path/to/git/repository"
  url            = "https://github.com/orga/owner.git"
  reference_name = "some-branch"
}

resource "git_clone" "bare" {
  directory = "/path/to/git/repository"
  url            = "https://github.com/orga/owner.git"
  reference_name = "some-branch"
  bare           = true
}
