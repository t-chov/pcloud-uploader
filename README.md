# pcloud-uploader

[![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/t-chov/pcloud-uploader/main/LICENSE)

A command-line client for [pCloud](https://my.pcloud.com/).

## Installation

```sh
go install github.com/t-chov/pcloud-uploader@latest
```

## Usage

```text
NAME:
   pcloud-uploader - A new cli application

USAGE:
   pcloud-uploader [global options] command [command options] [arguments...]

VERSION:
   0.1.0

COMMANDS:
   ls, listfolder  Receive data for a folder.
   up, uploadfile  Upload a file.
   help, h         Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```

## Requirements

You must have a registered pCloud account.

You can provide your credentials via environment variables to avoid interactive login prompts every time:

- `PCLOUD_USERNAME`: The email address associated with your pCloud account.
- `PCLOUD_PASSWORD`: Your pCloud login password.

```sh
export PCLOUD_USERNAME=john.smith@example.com
export PCLOUD_PASSWORD=IamJohnSmith
```

If these environment variables are not set, the application will prompt you to enter your username and password interactively upon running a command.

## Commands

### `listfolder` / `ls`

Receive metadata and structure for a folder on pCloud. It prints a tree-like outline of the files and directories.

```text
$ pcloud-uploader ls --help
NAME:
   pcloud-uploader ls - Receive data for a folder.

USAGE:
   pcloud-uploader ls [command options] <PATH>

OPTIONS:
   --recursive, -r  full directory tree will be returned (default: false)
   --showdeleted    deleted files and folders that can be undeleted will be displayed (default: false)
   --nofiles        only the folder (sub)structure will be returned (default: false)
   --noshares       only user's own folders and files will be displayed (default: false)
   --help, -h       show help
```

**Example:**
```sh
# List contents of the root directory
$ pcloud-uploader ls /

# List contents recursively
$ pcloud-uploader ls -r /Photos
```

### `uploadfile` / `up`

Upload a local file to a specified remote destination path. Prints the updated file's hash upon successful upload.

```text
$ pcloud-uploader up --help
NAME:
   pcloud-uploader up - Upload a file.

USAGE:
   pcloud-uploader up [command options] <SOURCE_FILE> <DEST_PATH>

OPTIONS:
   --help, -h  show help
```

**Example:**
```sh
# Upload a local file.txt to /Documents/file.txt on pCloud
$ pcloud-uploader up ./file.txt /Documents/file.txt
```
