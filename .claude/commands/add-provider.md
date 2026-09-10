---
description: Scaffold a new engage-ops provider (or a new resource type on an existing one)
---

Read CONTRIBUTING.md's "Adding a new provider" section and
`internal/provider/customerio/` (the reference implementation) before doing
anything else.

The user's request: $ARGUMENTS

Steps:
1. If this is a brand-new provider: create `internal/provider/<name>/provider.go`
   following the shape of `internal/provider/customerio/provider.go` —
   `Provider` struct embedding `provider.Base`, `New()` registering resource
   types (real ones if the API is known, `provider.NewStubResource`
   otherwise), `Configure` validating required credentials, and an `init()`
   calling `provider.Register`.
2. For each resource type being implemented for real (not stubbed), create
   `<resource>.go` implementing `core.Resource` — `Diff` delegating to
   `diff.Compute` unless a field requires replacement, `Apply` switching on
   `change.Action`, `Read` for drift detection. Persist any platform-assigned
   ID under an `_`-prefixed key (e.g. `_id`).
3. Add the provider's blank import to `internal/cli/root.go` if it's new.
4. Update the provider table in README.md's architecture section and the
   Status section if a previously-stubbed provider is now real.
5. Write unit tests for `Diff` and, if HTTP calls were added, for
   `Apply`/`Read` against a fake `httpClient` — not live credentials.
6. Run `gofmt -l .`, `go vet ./...`, `go build ./...`, `go test ./... -race`
   and fix anything they flag before considering the task done.

Ask before proceeding if $ARGUMENTS doesn't make clear which provider/
resource type is wanted, or whether real API integration vs. a stub is
intended.
