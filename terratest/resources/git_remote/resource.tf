resource "git_remote" "remote" {
  directory = "/path/to/git/repository"
  name      = "some-remote"
  urls      = ["https://github.com/some-org/some-repo.git"]
}
