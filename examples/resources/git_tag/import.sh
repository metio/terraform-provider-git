# git_tag resources can be imported by specifying the directory of the
# Git repository, the name of the tag to import, and the revision. All
# values are separated by a single '|'. The revision is optional and
# will default to 'HEAD' if not specified.
terraform import git_tag.tag 'path/to/your/git/repository|name-of-your-tag|revision'
