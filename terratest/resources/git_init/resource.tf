resource "git_init" "init" {
  directory = "/path/to/git/repository"
}

resource "git_init" "bare" {
  directory = "/path/to/git/repository"
  bare      = true
}
