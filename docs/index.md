<p style="text-align: center; margin-left: 1.6rem;">
  <img alt="Logotype depicting a lighthouse" src="./images/logo-450px.png" width="450" />
</p>
<h1 align="center">
  Watchtower
</h1>

<p align="center">
  A container-based solution for automating Docker container base image updates.
  <br/><br/>
  <a href="https://goreportcard.com/report/github.com/apivzero/watchtower">
    <img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/apivzero/watchtower" />
  </a>
  <a href="https://github.com/apivzero/watchtower/releases">
    <img alt="latest version" src="https://img.shields.io/github/tag/apivzero/watchtower.svg" />
  </a>
  <a href="https://www.apache.org/licenses/LICENSE-2.0">
    <img alt="Apache-2.0 License" src="https://img.shields.io/github/license/apivzero/watchtower.svg" />
  </a>
  <a href="https://github.com/apivzero/watchtower/pkgs/container/watchtower">
    <img alt="GHCR" src="https://img.shields.io/badge/GHCR-apivzero%2Fwatchtower-blue" />
  </a>
</p>

## Quick Start

With watchtower you can update the running version of your containerized app simply by pushing a new image to your
image registry. Watchtower will pull down your new image, gracefully shut down your existing container
and restart it with the same options that were used when it was deployed initially. Run the watchtower container with
the following command:

=== "docker run"

    ```bash
    $ docker run -d \
    --name watchtower \
    -v /var/run/docker.sock:/var/run/docker.sock \
    ghcr.io/apivzero/watchtower
    ```

=== "docker-compose.yml"

    ```yaml
    version: "3"
    services:
      watchtower:
        image: ghcr.io/apivzero/watchtower
        volumes:
          - /var/run/docker.sock:/var/run/docker.sock
    ```
