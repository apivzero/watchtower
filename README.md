<div align="center">

  <img src="./logo.png" width="450" />

  # Watchtower

  A process for automating Docker container base image updates.
  <br/><br/>

  [![Go Report Card](https://goreportcard.com/badge/github.com/apivzero/watchtower)](https://goreportcard.com/report/github.com/apivzero/watchtower)
  [![latest version](https://img.shields.io/github/tag/apivzero/watchtower.svg)](https://github.com/apivzero/watchtower/releases)
  [![Apache-2.0 License](https://img.shields.io/github/license/apivzero/watchtower.svg)](https://www.apache.org/licenses/LICENSE-2.0)
  [![GHCR](https://img.shields.io/badge/GHCR-apivzero%2Fwatchtower-blue)](https://github.com/apivzero/watchtower/pkgs/container/watchtower)

</div>

## About

This is a maintained fork of [containrrr/watchtower](https://github.com/containrrr/watchtower), which was archived in December 2025. I run this for my own infrastructure and keep it working with modern Docker. If you stumble across it and find it useful, you're welcome to use it too.

I'm not offering heavy support, but I do welcome PRs.

### Key changes from upstream

- Docker SDK upgraded to v27.x with API version negotiation, fixing compatibility with Docker 28+ (which requires API v1.44+)
- Go version bumped to 1.22+
- Module path updated to `github.com/apivzero/watchtower`
- Docker labels (`com.centurylinklabs.watchtower.*`) are unchanged for backward compatibility

### Acknowledgments

This project was originally built by the [containrrr](https://github.com/containrrr) community. Huge thanks to all of the [original contributors](https://github.com/containrrr/watchtower/graphs/contributors) who designed, built, and maintained watchtower over the years. This fork wouldn't exist without their work.

## Quick Start

Watchtower updates the running version of your containerized app when a new image is pushed to Docker Hub or your own image registry. It pulls down the new image, gracefully shuts down the existing container, and restarts it with the same options that were used when it was deployed initially.

```
$ docker run --detach \
    --name watchtower \
    --volume /var/run/docker.sock:/var/run/docker.sock \
    ghcr.io/apivzero/watchtower
```

> **Note:** Watchtower is intended for homelabs, media centers, local dev environments, and similar. It is not recommended for commercial or production use. For that, look into Kubernetes or lighter-weight alternatives like [MicroK8s](https://microk8s.io/) and [k3s](https://k3s.io/).

## Releases

Container images are published to GHCR at [`ghcr.io/apivzero/watchtower`](https://github.com/apivzero/watchtower/pkgs/container/watchtower).

### How it works

There are two CI pipelines that produce artifacts:

- **Push to main** — every merge to `main` builds and pushes `ghcr.io/apivzero/watchtower:latest-dev`. This is a rolling dev image, always built from the latest commit on `main`.
- **Tagged release** — pushing a semver tag (e.g. `v2.0.0`) triggers a full production release: goreleaser builds multi-arch binaries and Docker images, creates a GitHub Release with archives, and pushes versioned + `latest` tags to GHCR.

### Creating a release

When you're ready to cut a release:

```bash
git tag v2.0.0
git push origin v2.0.0
```

This triggers the release workflow, which will:
1. Lint and test
2. Build binaries for linux/amd64, linux/386, linux/arm, linux/arm64
3. Build and push multi-arch Docker images to GHCR
4. Create a multi-arch manifest (`ghcr.io/apivzero/watchtower:2.0.0` and `:latest`)
5. Create a GitHub Release with downloadable archives
6. Notify pkg.go.dev of the new module version

There's no fixed release cadence. Tag a release whenever there are meaningful changes worth shipping.

### Manual workflow dispatch

The release workflow also supports manual triggering via the GitHub Actions UI ("Run workflow" button). When run on an untagged commit, it performs a snapshot build (lint, test, compile) without publishing anything. This is useful for verifying the release pipeline works end-to-end.

## Documentation

For general usage documentation, refer to the [upstream watchtower docs](https://containrrr.dev/watchtower). Note that image references should use `ghcr.io/apivzero/watchtower` instead of `containrrr/watchtower`.
