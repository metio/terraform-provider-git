data "git_config" "config" {
  directory = "/path/to/git/repository"
  scope     = "local"
}
