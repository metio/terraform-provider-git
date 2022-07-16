data "git_tags" "all_tags" {
  directory = "/path/to/git/repository"
}

data "git_tags" "annotated_tags" {
  directory   = "/path/to/git/repository"
  annotated   = true
  lightweight = false
}

data "git_tags" "lightweight_tags" {
  directory   = "/path/to/git/repository"
  annotated   = false
  lightweight = true
}
