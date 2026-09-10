# Contributing to engage-ops

## Adding a new provider

Every provider lives in `internal/provider/<name>/` and implements
[`core.Provider`](internal/core/resource.go). `internal/provider/customerio` is
the reference implementation — copy its shape for a new platform:

1. Create `internal/provider/<name>/provider.go`:
   - Define a `Provider` struct embedding `provider.Base`.
   - `New() core.Provider` builds it and registers each resource type via
     `provider.NewBase(name, map[string]core.Resource{...})`.
   - `Configure(ctx, config core.Attributes) error` validates required
     credentials from the YAML `providers:` block and stores an API client.
   - Register the provider in an `init()` with `provider.Register(name, New)`.
2. For each resource type (e.g. `campaign`, `segment`), create a
   `<resource>.go` implementing `core.Resource`:
   - `Diff` — for a straightforward create/update/delete resource, just
     delegate to `diff.Compute`. Only write custom diff logic if some field
     changes require replacement (delete + recreate) instead of an in-place
     update.
   - `Apply` — switch on `change.Action` and call your platform's API,
     returning the attributes to persist in state (including any
     platform-assigned ID under an `_`-prefixed key, e.g. `_id`, so
     `diff.Equal` ignores it when comparing against desired config).
   - `Read` — fetch current remote state for drift detection; return
     `nil, nil` if the resource no longer exists remotely.
3. Wire the new package into `internal/cli/root.go`'s blank imports so its
   `init()` runs.
4. Add it to the provider list in `README.md` and note its resource types.

If you're scaffolding a provider you don't plan to fully implement yet,
register its resource types with `provider.NewStubResource(name, type)` —
this lets `engageops providers`/`validate`/`plan` work against it while
`apply` fails with a clear "not yet implemented" error instead of silently
no-opping.

## Testing

```
go build ./...
go vet ./...
go test ./... -race -cover
gofmt -l .   # should print nothing
```

Add unit tests next to the code you touch (`*_test.go`). Provider `Apply`/
`Read` methods that call real HTTP APIs should be tested against the
`httpClient` interface (see `internal/provider/customerio/http.go`) with a
fake transport, not live credentials.

## Commit style

Keep commits focused and explain *why* a change was made in the body when
it's not obvious from the diff alone.
