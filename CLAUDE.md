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
- `gofmt -w *.go` — format before committing

## Architecture

All Go sources live at the repo root in a single `main` package. No subdirectories or subpackages.

- **`main.go`** — CLI entry point: defines commands (`ls`, `up`), handles credential loading from `PCLOUD_USERNAME`/`PCLOUD_PASSWORD` env vars or stdin, runs `login()` as a pre-command hook to obtain an auth token
- **`auth.go`** — authenticates against `api.pcloud.com/userinfo` and returns a session token
- **`listfolder.go`** — wraps `api.pcloud.com/listfolder`, recursively prints folder tree with colored output
- **`uploadfolder.go`** — wraps `api.pcloud.com/uploadfile` via multipart POST, prints file hash on success

New commands should be added as additional top-level `.go` files named with short verbs matching the CLI aliases.

## Coding Conventions

- Standard Go conventions enforced by `gofmt` (tabs, capitalized exports)
- Keep each file focused on a single command
- Environment variables use uppercase snake case
- Tests go in `_test.go` siblings; mock HTTP by extracting API calls into interfaces
- Commit messages: short imperative summaries (e.g., "change binary name", "fix installation command")

## Security

Credentials are read at runtime from env vars — never commit secrets. Sanitize logs before sharing as folder listings echo remote paths.
