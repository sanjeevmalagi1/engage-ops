# CLAUDE.md

Guidance for Claude Code (and other agents) working in this repository.

## What this project is

`engage-ops` is a Terraform-style CLI, written in Go, that manages resources
(segments, campaigns, audiences, ...) across customer engagement platforms
— Customer.io, WebEngage, MoEngage, Twilio, Mailchimp — from declarative
YAML. The core workflow is `engageops plan` → `engageops apply`, backed by a
local JSON state file, mirroring Terraform's provider/resource/plan/apply
model exactly (see [README.md](README.md) for the full architecture map).

Read [README.md](README.md) for user-facing docs and
[CONTRIBUTING.md](CONTRIBUTING.md) before adding or modifying a provider —
it documents the required shape in detail.

## Architecture, in one paragraph

`internal/core` defines two interfaces — `Provider` and `Resource` — that
every engagement platform integration implements once
(`internal/core/resource.go`). `internal/core/plan` drives them generically:
it loads config (`internal/config`), loads state (`internal/core/state`),
configures each referenced provider, calls `Resource.Diff` for every
declared resource plus every resource in state no longer declared (→
deletion), and `Plan.Apply` calls `Resource.Apply` for each non-noop change.
The CLI (`internal/cli`, cobra-based) is a thin wrapper over `config` +
`plan` + `state`. `internal/provider/customerio` is the one fully-wired
reference provider; the other four are registered with
`provider.NewStubResource` placeholders — see CONTRIBUTING.md before
building one out.

## Working in this repo

- Before committing, run: `gofmt -l .` (must print nothing), `go vet ./...`,
  `go build ./...`, `go test ./... -race -cover`. CI
  (`.github/workflows/ci.yml`) runs the same checks and will fail the PR
  otherwise.
- Follow existing package boundaries. Don't add a dependency from
  `internal/core` on any `internal/provider/*` package — the direction of
  dependency is provider → core, never the reverse, or the plan/apply engine
  stops being generic.
- New provider work follows CONTRIBUTING.md's "Adding a new provider"
  section exactly — same file layout, same `Diff`/`Apply`/`Read` shape as
  `internal/provider/customerio`.
- State-shaped bugs (a resource re-created every plan, a delete never
  detected) are almost always in `diff.Compute`/`diff.Equal`
  (`internal/core/diff/diff.go`) or in a provider's `Apply` not persisting
  an `_`-prefixed field it should. Check there first.
- Never commit a real `*.tfstate.json` file or real API credentials, even in
  example/test fixtures — use `REPLACE_ME` or obviously-fake values, as the
  existing examples do.
- Keep provider `Apply`/`Read` implementations testable against a fake
  `httpClient` (see `internal/provider/customerio/http.go`) rather than
  hitting real APIs in tests.

## Conventions this codebase already follows — match them

- No doc comments beyond a one-line package comment and short comments on
  exported types/functions explaining *why*, not what. No comment blocks
  restating the code.
- Errors are wrapped with `fmt.Errorf("...: %w", err)` and include the
  resource address or file path they concern.
- CLI commands return errors from `RunE` rather than calling `os.Exit`
  directly; `cmd/engageops/main.go` is the only place that exits the
  process.
