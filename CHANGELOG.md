# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [v2.1.0] - 2026-02-22

### Security

- **TLS certificate verification is now enabled by default for registry connections.**
  Previously, watchtower silently disabled TLS verification (`InsecureSkipVerify`) for
  all registry digest checks, exposing registry credentials to potential man-in-the-middle
  attacks on the network path between watchtower and the registry.

### Added

- `--no-tls-verify` flag (`WATCHTOWER_NO_TLS_VERIFY=true`) to opt out of TLS certificate
  verification for registry connections. This replaces the previous unconditional insecure
  behaviour and is intended only for private/self-hosted registries using self-signed
  certificates.

### Migration guide for self-hosted registry users

If you run watchtower against a private registry with a **self-signed or internally-signed
TLS certificate**, you must now explicitly pass `--no-tls-verify` (or set the environment
variable `WATCHTOWER_NO_TLS_VERIFY=true`). Without it, watchtower will fail to verify the
registry's certificate and fall back to a full image pull for every check cycle.

**Docker CLI / Docker Compose:**

```yaml
services:
  watchtower:
    image: ghcr.io/apivzero/watchtower
    environment:
      - WATCHTOWER_NO_TLS_VERIFY=true
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
```

**Command line:**

```sh
docker run ghcr.io/apivzero/watchtower --no-tls-verify
```

If you are unsure whether your registry uses a self-signed certificate, watchtower will now
log a clear warning when a certificate error is detected:

```
WARN Could not do a head request for "registry.internal/myimage:latest", falling back to regular pull.
WARN Reason: tls: failed to verify certificate: x509: certificate signed by unknown authority
WARN If "registry.internal/myimage:latest" uses a self-signed certificate, set --no-tls-verify (env: WATCHTOWER_NO_TLS_VERIFY=true) to disable certificate verification.
```

Users connecting to public registries (Docker Hub, GHCR, ECR, GCR, etc.) are unaffected —
those registries use trusted certificates and no configuration change is needed.

---

## [v2.0.1] - 2025-01-17

### Fixed

- Fix arm64 platform string in goreleaser `dockers_v2` config.
- Fix Dockerfile for goreleaser `dockers_v2` build context.

---

## [v2.0.0] - 2025-01-17

Initial release of the `apivzero/watchtower` fork from the archived
`containrrr/watchtower`. See the [migration guide](docs/migration.md) for
details on upgrading from `containrrr/watchtower`.

### Changed

- Module path changed to `github.com/apivzero/watchtower`.
- Container images published to GHCR only (`ghcr.io/apivzero/watchtower`).
- Upgraded Docker SDK to v27.5.1 with API version negotiation for Docker 28+ compatibility.
