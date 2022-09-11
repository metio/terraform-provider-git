# add single file
resource "git_add" "single_file" {
  directory = "/path/to/git/repository"
  add_paths = ["path/to/file/in/repository"]
}

# add all files in directory and its subdirectory recursively
resource "git_add" "single_directory" {
  directory = "/path/to/git/repository"
  add_paths = ["path/to/directory/in/repository"]
}

# add files matching pattern
resource "git_add" "glob_pattern" {
  directory = "/path/to/git/repository"
  add_paths = ["path/*/in/repo*"]
}

# mix exact paths and glob patterns
resource "git_add" "glob_pattern" {
  directory = "/path/to/git/repository"
  add_paths = [
    "path/*/in/repo*",
    "another/path/to/file/here",
    "this/could/be/a/directory",
  ]
}
