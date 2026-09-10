# engage-ops

Infrastructure-as-code for customer engagement platforms. `engage-ops` lets
you declare segments, campaigns, and other resources in Customer.io,
WebEngage, MoEngage, Twilio, and Mailchimp as YAML, then plan and apply
changes the way [Terraform](https://www.terraform.io) manages cloud
infrastructure.

```
engageops plan     # show what would change, touch nothing
engageops apply     # make it so
```

## Why

Engagement platform configuration (segments, campaigns, journeys, audiences)
tends to live only in each vendor's dashboard: no diff, no review, no audit
trail, no way to promote a change from staging to production. `engage-ops`
brings the same discipline Terraform brought to cloud infrastructure —
declarative config, a plan/apply workflow, and a state file — to this
category of tooling.

## Status

Early scaffold. The plan/apply engine, state file, and config loader are
functional. **Customer.io** (`segment` resource) is the one fully wired
provider today; WebEngage, MoEngage, Twilio, and Mailchimp are registered
with stubbed resources so the CLI and config validate against them, but
`apply` will error with "not yet implemented" until each is built out. See
[CONTRIBUTING.md](CONTRIBUTING.md) for how to add a provider.

## Install

Requires Go 1.26+.

```sh
go build -o bin/engageops ./cmd/engageops
```

## Quick start

```sh
mkdir my-engagement-config && cd my-engagement-config
engageops init                 # scaffolds engageops.yml
$EDITOR engageops.yml          # fill in provider credentials + resources
engageops validate             # parse config, resolve providers/resource types
engageops plan                 # preview changes
engageops apply                # apply, with interactive confirmation
```

See [`examples/basic`](examples/basic) for a runnable example config.

## Configuration

Config is one or more `*.yml`/`*.yaml` files in a directory (default: `.`,
override with `--dir`). Files are merged; every declared resource address
must be unique across the directory.

```yaml
providers:
  customerio:
    app_api_key: "REPLACE_ME"

resources:
  - provider: customerio
    type: segment
    name: active_users          # local name; forms the address customerio_segment.active_users
    attributes:
      name: "Active Users"      # fields passed straight to the provider
```

`engageops` does not yet expand `${VAR}` placeholders in YAML (see
[Roadmap](#roadmap)) — keep real credentials out of version control via a
gitignored local file or a config generated from a secrets manager.

## State

Applied resources are tracked in a local JSON state file (default:
`engageops.tfstate.json`, override with `--state`), analogous to
`terraform.tfstate`. It maps each resource address to the attributes last
written to the remote platform (including platform-assigned IDs), so future
plans diff against reality instead of assuming a clean slate.

Treat the state file like a secret: it may contain identifiers and
attributes from your engagement platforms. It's gitignored by default —
don't force-add it.

## Commands

| Command                | Purpose                                                              |
|-------------------------|-----------------------------------------------------------------------|
| `engageops init`        | Scaffold a starter config file                                       |
| `engageops validate`    | Parse config and resolve provider/resource references, no API calls  |
| `engageops plan`        | Diff config against state, print planned changes                     |
| `engageops apply`       | Apply planned changes (prompts for confirmation unless `--auto-approve`) |
| `engageops destroy`     | Delete every resource currently tracked in state                     |
| `engageops providers`   | List registered providers and their resource types                   |

## Architecture

```
cmd/engageops/            CLI entrypoint
internal/cli/              cobra commands (plan, apply, destroy, ...)
internal/config/           YAML config loading
internal/core/              Provider/Resource interfaces, Address, Change
internal/core/diff/         generic attribute-map diffing
internal/core/state/        state file read/write
internal/core/plan/         plan computation + apply orchestration
internal/provider/          provider registry + built-in providers
  customerio/                reference provider (segment resource, real HTTP client)
  webengage/ moengage/       stubbed providers — see CONTRIBUTING.md
  twilio/ mailchimp/
```

The core abstraction is `core.Provider` / `core.Resource`
([internal/core/resource.go](internal/core/resource.go)) — a small interface
each engagement platform implements once. The plan/apply engine
([internal/core/plan](internal/core/plan)) knows nothing about any specific
platform; it just diffs config against state and calls `Diff`/`Apply` on
whatever `Resource` a provider hands back.

## Roadmap

- Environment variable / secrets-manager interpolation in YAML config
- Remote state backends (S3, GCS) with locking, for team use
- WebEngage, MoEngage, Twilio, Mailchimp resource implementations
- `engageops import` to adopt existing remote resources into state
- `engageops state list/show/rm` subcommands

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

## License

[MIT](LICENSE)
