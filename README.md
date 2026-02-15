<div align="center">

  <img src="./logo.png" width="450" />

  # Watchtower

  A process for automating Docker container base image updates.
  <br/><br/>

  [![Go Report Card](https://goreportcard.com/badge/github.com/apivzero/watchtower)](https://goreportcard.com/report/github.com/apivzero/watchtower)
  [![latest version](https://img.shields.io/github/tag/apivzero/watchtower.svg)](https://github.com/apivzero/watchtower/releases)
  [![Apache-2.0 License](https://img.shields.io/github/license/apivzero/watchtower.svg)](https://www.apache.org/licenses/LICENSE-2.0)
  [![Pulls from DockerHub](https://img.shields.io/docker/pulls/apivzero/watchtower.svg)](https://hub.docker.com/r/apivzero/watchtower)

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
    apivzero/watchtower
```

> **Note:** Watchtower is intended for homelabs, media centers, local dev environments, and similar. It is not recommended for commercial or production use. For that, look into Kubernetes or lighter-weight alternatives like [MicroK8s](https://microk8s.io/) and [k3s](https://k3s.io/).

## Documentation

For general usage documentation, refer to the [upstream watchtower docs](https://containrrr.dev/watchtower). Note that image references should use `apivzero/watchtower` instead of `containrrr/watchtower`.
