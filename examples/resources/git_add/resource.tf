# add single file
resource "git_add" "file" {
  directory  = "/path/to/git/repository"
  exact_path = "path/to/file/in/repository"
}

# add all files in directory and its subdirectory recursively
resource "git_add" "directory" {
  directory  = "/path/to/git/repository"
  exact_path = "path/to/directory/in/repository"
}

# add files matching pattern
resource "git_add" "glob" {
  directory = "/path/to/git/repository"
  glob_path = "path/*/in/repo*"
}

# add all modified files
resource "git_add" "all" {
  directory = "/path/to/git/repository"
  all       = true
}
