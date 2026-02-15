# AGENTS.md

Instructions for AI coding agents working on this codebase. See https://agents.md for spec.

## Project Overview

Watchtower is a Go application that automatically updates running Docker containers when new images are pushed. It monitors containers, pulls updated images, gracefully stops containers, and restarts them with the same configuration.

This is a maintained fork of the archived `containrrr/watchtower`. The module path is `github.com/apivzero/watchtower`.

## Build and Test

```bash
go build ./...          # compile everything
go test ./...           # run all tests
go vet ./...            # check for issues
go fmt ./...            # format code
```

Tests must pass before any change is considered complete.

## Architecture

### Package Layout

```
cmd/                      CLI commands (cobra). Entry point is cmd/root.go
internal/
  actions/                Core update logic: staleness check, stop/start, cleanup
  actions/mocks/          Mock client/container for action tests
  flags/                  CLI flag registration, env var binding (viper), secret loading
  meta/                   Version metadata (set via ldflags)
  util/                   Small helpers (slice ops, random SHA generation)
pkg/
  api/                    HTTP API server (update trigger + Prometheus metrics)
  container/              Docker container abstraction and Docker SDK client wrapper
  container/mocks/        Mock HTTP handlers simulating Docker API responses
  filters/                Container selection predicates (name, label, scope, image)
  lifecycle/              Pre/post check/update hook execution
  metrics/                Prometheus metric definitions and async collection
  notifications/          Notification dispatch (Shoutrrr, legacy email/slack/teams/gotify)
  registry/               Docker registry auth, image digest comparison, manifest fetching
  session/                Update session progress tracking and reporting
  sorter/                 Dependency-aware container sorting (topological)
  types/                  Core interfaces: Container, Client, Filter, Notifier, Report
```

### Key Interfaces (all in `pkg/types/`)

- `Container` — abstraction over Docker container metadata and labels
- `Filter` — `func(FilterableContainer) bool` predicate for container selection
- `Notifier` — notification lifecycle (start batch, send, close)
- `Report` — update session results (scanned, updated, failed, skipped)

The concrete Docker client is in `pkg/container/client.go` implementing the `Client` interface from the same package.

### Update Flow

1. `cmd/root.go` sets up scheduler (cron) or runs once
2. `internal/actions/update.go` → `Update()` is the core loop
3. List containers → filter → check staleness → sort by dependencies
4. For each stale container: run pre-update hook → stop → start with new image → run post-update hook
5. Clean up old images if `--cleanup` is set
6. Send notification batch with results

### Docker Labels

All labels use the prefix `com.centurylinklabs.watchtower.*`. This is a user-facing API inherited from upstream — do NOT change the prefix. Existing containers in the wild have these labels applied.

Key labels (defined in `pkg/container/metadata.go`):
- `.enable` — opt-in/out with `--label-enable`
- `.monitor-only` — check for updates but don't restart
- `.no-pull` — skip image pull
- `.stop-signal` — custom stop signal (default SIGTERM)
- `.depends-on` — comma-separated dependency list for restart ordering
- `.scope` — isolate multiple watchtower instances
- `.lifecycle.{pre,post}-{check,update}` — hook commands
- `.lifecycle.{pre,post}-update-timeout` — hook timeout in minutes

### Docker SDK

Uses `github.com/docker/docker` v27.5.1 with API version negotiation (`WithAPIVersionNegotiation()`). This is critical for Docker 28+ compatibility. The minimum API version is 1.44 (set in `internal/flags/flags.go`).

When modifying Docker SDK interactions, use the new type locations:
- `container.ListOptions`, `container.StartOptions`, `container.RemoveOptions` (not `types.*`)
- `container.ExecOptions`, `container.ExecStartOptions` (not `types.ExecConfig`/`ExecStartCheck`)
- `image.PullOptions`, `image.RemoveOptions` (not `types.*`)
- `image.DeleteResponse` (not `types.ImageDeleteResponseItem`)

Types that still have aliases and are fine to use from `types`: `ContainerJSON`, `ContainerJSONBase`, `ContainerState`, `ImageInspect`, `Container`.

## Testing Patterns

### Frameworks

- **Ginkgo + Gomega** — used for most tests (`Describe`/`When`/`It` structure)
- **testify** — used in `internal/flags/` tests
- Both coexist; follow the pattern of the package you're modifying

### Test Structure

Each tested package has:
- `*_suite_test.go` — Ginkgo suite bootstrap
- `*_test.go` — test specs

### Mock Patterns

**Container mocks** (`pkg/container/container_mock_test.go`):
```go
container := MockContainer(WithImageName("docker.io/library/nginx:latest"))
container := MockContainer(WithLabels(map[string]string{"key": "val"}))
container := MockContainer(WithContainerState(types.ContainerState{Running: true}))
```

**Docker API mocks** (`pkg/container/mocks/`):
- Uses `ghttp.Server` (Gomega HTTP) to simulate Docker daemon
- JSON fixtures in `pkg/container/mocks/data/`
- Handlers: `ListContainersHandler()`, `GetContainerHandler()`, `RemoveImageHandler()`, etc.

**Action mocks** (`internal/actions/mocks/`):
- `MockClient` implementing the `container.Client` interface
- Simple stubs returning configured errors/data

### Test Data

- Mock JSON files: `pkg/container/mocks/data/*.json`
- These represent Docker API responses (container inspect, image inspect, container list)
- When adding new test scenarios, add fixtures here if needed

## Code Style

- Standard Go conventions, `gofmt` formatted
- No unnecessary comments or docstrings on obvious code
- Package-level doc comments on exported types/functions
- Error handling: wrap with context, don't swallow errors silently
- Use `log` (aliased from `logrus`) for logging, not `fmt.Print`

## CI

- GitHub Actions workflows in `.github/workflows/`
- All actions pinned to commit hashes (not floating tags)
- Go version: 1.24.x across all workflows
- Staticcheck for linting
- goreleaser v2 for builds and releases
- Container images published to GHCR only (`ghcr.io/apivzero/watchtower`), no Docker Hub
- All workflows authenticate to GHCR via `GITHUB_TOKEN` (no custom secrets needed)

## Things to Watch Out For

- `containrrr/shoutrrr` is an **external dependency** — don't rename its import path
- Container recreation removes image-default config to avoid overwriting the new image's defaults (`GetCreateConfig()`)
- The update lock is channel-based (`make(chan bool, 1)`) — shared between scheduler and HTTP API
- `EX_TEMPFAIL` (exit code 75) from lifecycle hooks means "skip this update" — not an error
- Circular dependency detection in `pkg/sorter/` returns an error, doesn't silently skip
- `_FILE` suffix on env vars triggers secret-from-file loading (e.g., `WATCHTOWER_HTTP_API_TOKEN_FILE`)
- The `tplprev/` directory contains a WASM build of the template previewer — has build tags
