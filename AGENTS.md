# Repository Guidelines

## Project Structure & Module Organization
Go sources live at the repo root: `main.go` wires the CLI, `auth.go` owns credential handling, and `listfolder.go` / `uploadfolder.go` wrap the respective pCloud API calls. The Makefile drives compilation, release packaging, and version bumping; create new commands as additional top-level `.go` files until a dedicated package split becomes necessary.

## Build, Test, and Development Commands
Use `make build` to produce a `pcloud-uploader` binary with revision metadata, and `make install` to drop the tool into `$GOBIN`. `go run . ls /Backup` is the quickest feedback loop when iterating on a command. `make test` runs `go test -v ./...` after compiling, and `make clean` deletes build outputs plus cached artifacts.

## Coding Style & Naming Conventions
Follow standard Go conventions enforced by `gofmt` (tabs for indentation, exported names start with caps). Keep files focused on a single command and name new commands with short verbs mirroring the CLI aliases (`listfolder`, `uploadfile`, etc.). Config keys and environment variables are uppercase snake case; run `gofmt -w *.go` before committing.

## Testing Guidelines
Author tests with Go’s `testing` package and place them in `_test.go` siblings (e.g., `uploadfolder_test.go`). Mock HTTP interactions by extracting API calls into interfaces so unit tests can replace them with in-memory fakes. Target `go test ./...` as the baseline gate, and document any long-running integration tests separately.

## Commit & Pull Request Guidelines
Git history uses short, imperative summaries such as `change binary name` or `fix installation command`; follow that style and keep bodies for context only when needed. Every pull request should describe the user-visible change, include CLI examples if flags were modified, and link related GitHub issues. Confirm `make test` passes and mention how credentials should be supplied for reviewers (e.g., dummy env vars).

## Security & Configuration Tips
The client reads `PCLOUD_USERNAME` and `PCLOUD_PASSWORD` at runtime—load them via a password manager or `.envrc`, but never commit secrets. Prefer temporary application passwords over full account credentials, and sanitize logs before sharing because folder listings echo remote paths.
