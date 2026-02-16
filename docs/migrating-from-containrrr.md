# Migrating from containrrr/watchtower

This fork is a drop-in replacement for [containrrr/watchtower](https://github.com/containrrr/watchtower), which was archived in December 2025. The only change required is swapping the image name. All Docker labels (`com.centurylinklabs.watchtower.*`) are unchanged, so existing per-container configuration continues to work.

## docker run

If you started watchtower with `docker run`:

1. Find the running container:

    ```bash
    docker ps --filter "name=watchtower" --no-trunc
    ```

2. Inspect it and save your configuration — restart policy, environment variables, volume mounts, port bindings, network mode, and any extra flags:

    ```bash
    docker inspect watchtower
    ```

    Look at `HostConfig.RestartPolicy`, `HostConfig.Binds`, `Config.Env`, `HostConfig.PortBindings`, and `HostConfig.NetworkMode` in the output. You'll need these to recreate the container with the same options.

3. Stop and remove the old container:

    ```bash
    docker stop watchtower && docker rm watchtower
    ```

4. Pull the new image and start with the same options:

    ```bash
    docker pull ghcr.io/apivzero/watchtower
    docker run -d \
      --name watchtower \
      --restart unless-stopped \
      -v /var/run/docker.sock:/var/run/docker.sock \
      ghcr.io/apivzero/watchtower
    ```

    Add back any environment variables (`-e`), volume mounts (`-v`), port bindings (`-p`), or extra flags from the inspect output in step 2.

5. Verify:

    ```bash
    docker ps --filter name=watchtower
    docker logs watchtower --tail 20
    ```

6. Clean up old images and any leftover networks:

    ```bash
    docker image rm containrrr/watchtower
    docker network prune -f
    ```

## docker-compose

Update the `image:` field in your compose file:

```yaml
services:
  watchtower:
    image: ghcr.io/apivzero/watchtower  # was containrrr/watchtower
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
```

Then:

```bash
docker compose pull watchtower
docker compose up -d watchtower
```

After verifying everything works, clean up:

```bash
docker image prune -f
docker network prune -f
```

## What changed

- **Image location**: `containrrr/watchtower` (Docker Hub) → `ghcr.io/apivzero/watchtower` (GHCR)
- **Docker SDK**: Upgraded to v27.x with API version negotiation, fixing compatibility with Docker 28+
- **Labels**: Unchanged. `com.centurylinklabs.watchtower.enable`, `com.centurylinklabs.watchtower.monitor-only`, and all other label-based configuration works as before.
- **CLI flags and environment variables**: Unchanged.
