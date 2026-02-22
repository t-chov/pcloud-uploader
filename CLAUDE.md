# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

A Go CLI client for the [pCloud](https://my.pcloud.com/) cloud storage API. Built with `urfave/cli/v2` for command routing and `fatih/color` for terminal output.

## Build & Development Commands

- `make build` — compile binary (`pcloud-uploader`) with git revision metadata
- `make install` — install to `$GOBIN`
- `make test` — build then run `go test -v ./...`
- `make clean` — remove build outputs and cached artifacts
- `go run . ls /Backup` — quickest feedback loop when iterating on a command
- `gofmt -w .` — format before committing

## Architecture

The project uses two packages: the `main` package at the repo root for CLI wiring, and the `pcloud` subpackage for API client logic.

### `pcloud/` — API client package

- **`pcloud/client.go`** — `API` interface (`Authenticate`, `ListFolder`, `UploadFile`), `Client` struct with `BaseURL`/`HTTPClient` fields, `NewClient` constructor
- **`pcloud/auth.go`** — `Client.Authenticate` method: obtains a session token from `api.pcloud.com/userinfo`
- **`pcloud/listfolder.go`** — `Client.ListFolder` method + `ListFolderResult`, `FolderMetadata` types
- **`pcloud/uploadfile.go`** — `Client.UploadFile` method + `UploadFileResult`, `FileMetadata`, `Checksum` types

### Root (`main` package) — CLI layer

- **`main.go`** — minimal entry point: `main()` + constants (`APP_NAME`, `VERSION`, env var names)
- **`commands.go`** — CLI wiring: `run`, `runWithClient`, `buildCommands`, `login` hook, `loadFromInput`, `printListfolder`, `printError`

To test commands, inject a mock `pcloud.API` into `runWithClient`. New API operations should be added as methods on `Client` and exposed via the `API` interface. New CLI commands should be added in `buildCommands` in `commands.go`.

## Coding Conventions

- Standard Go conventions enforced by `gofmt` (tabs, capitalized exports)
- Keep each file focused on a single command
- Environment variables use uppercase snake case
- Tests go in `_test.go` siblings; mock HTTP by extracting API calls into interfaces
- Commit messages: short imperative summaries (e.g., "change binary name", "fix installation command")

## Security

Credentials are read at runtime from env vars — never commit secrets. Sanitize logs before sharing as folder listings echo remote paths.
