# Migrating from containrrr/watchtower

This fork is a drop-in replacement for [containrrr/watchtower](https://github.com/containrrr/watchtower), which was archived in December 2025. The only change required is swapping the image name. All Docker labels (`com.centurylinklabs.watchtower.*`) are unchanged, so existing per-container configuration continues to work.

## docker run

If you started watchtower with `docker run`:

1. Find the running container:

    ```bash
    docker ps --filter "name=watchtower" --no-trunc
    ```

2. Inspect it and note your configuration (restart policy, environment variables, volume mounts, port bindings, network mode):

    ```bash
    docker inspect watchtower
    ```

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

    Add back any environment variables (`-e`), volume mounts (`-v`), port bindings (`-p`), or extra flags you had before.

5. Verify:

    ```bash
    docker ps --filter name=watchtower
    docker logs watchtower --tail 20
    ```

6. Clean up the old image:

    ```bash
    docker image rm containrrr/watchtower
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
docker image prune -f
```

## What changed

- **Image location**: `containrrr/watchtower` (Docker Hub) → `ghcr.io/apivzero/watchtower` (GHCR)
- **Docker SDK**: Upgraded to v27.x with API version negotiation, fixing compatibility with Docker 28+
- **Labels**: Unchanged. `com.centurylinklabs.watchtower.enable`, `com.centurylinklabs.watchtower.monitor-only`, and all other label-based configuration works as before.
- **CLI flags and environment variables**: Unchanged.
