# pcloud-uploader

A command line client for [pCloud](https://my.pcloud.com/)

## Installation

```
go install github.com/t-chov/pcloud@latest
```

## Usage

```
NAME:
   pcloud - A new cli application

USAGE:
   pcloud [global options] command [command options] [arguments...]

VERSION:
   0.0.1

COMMANDS:
   ls, listfolder  Receive data for a folder.
   up, uploadfile  Upload a file.
   help, h         Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```

## Requirements

You have to register pCloud.

And set two environment variables

- `PCLOUD_USERNAME` : email address you have registered.
- `PCLOUD_PASSWORD` : login password for pCloud.

```
$ export PCLOUD_USERNAME=john.smith@example.com
$ export PCLOUD_PASSWORD=IamJohnSmith
```

## Commands

```
$ pcloud listfolder --help
NAME:
   pcloud ls - Receive data for a folder.

USAGE:
   pcloud ls [command options] <PATH>

OPTIONS:
   --recursive, -r  full directory tree will be returned (default: false)
   --showdeleted    deleted files and folders that can be undeleted will be displayed (default: false)
   --nofiles        only the folder (sub)structure will be returned (default: false)
   --noshares       only user's own folders and files will be displayed (default: false)
   --help, -h       show help
```

```
$ pcloud uploadfile --help
NAME:
   pcloud up - Upload a file.

USAGE:
   pcloud up [command options] <SOURCE_FILE> <DEST_PATH>

OPTIONS:
   --help, -h  show help
```
