#
# SPDX-FileCopyrightText: Sebastian Hoß <seb@hoß.de>
# SPDX-License-Identifier: BSD0
#

# git_tag resources can be imported by specifying the directory of the Git repository and the name of the tag to
# import. Both values are separated by a single '|'.
terraform import git_remote.remote 'path/to/your/git/repository|name-of-your-tag'
