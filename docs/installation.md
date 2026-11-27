# Installation

## Docker Compose

The checked-in Compose definition is `docker/docker-compose.yml`. It configures the `ghcr.io/apimgr/wordlist:latest` image, publishes the service port, mounts persistent configuration and data directories, and checks `/api/v1/server/healthz`.

Review the published port and volume locations before deployment. The included mapping binds to `172.17.0.1:64580`; adjust it for the host network where the service will run.

## Build from source

Go is used only inside `casjaysdev/go:latest` through the Makefile:

```sh
make dev
```

This builds both `wordList` and `wordList-cli` in a temporary output directory. `make local` writes local binaries under the ignored `binaries/` directory; `make build` cross-compiles the configured platform matrix.

## Release binaries

Download the matching `wordList-{os}-{arch}` and `wordList-cli-{os}-{arch}` artifacts from the [GitHub Releases page](https://github.com/apimgr/wordList/releases). Windows artifacts use the `.exe` suffix. Verify the release checksums before installation and place the binaries in a directory on `PATH`.

## Running the server

Start `wordList` directly. On first start it creates `server.yml` in its platform-specific configuration directory. The default listen address is `[::]` with a random port from `64000` through `64999`; application containers default to port `80`. Set `PORT`, `CONFIG_DIR`, and `DATA_DIR` as needed.

The image entrypoint supplies defaults for `TZ` (`America/New_York`), `CONFIG_DIR` (`/config/wordList`), `DATA_DIR` (`/data/wordList`), `LISTEN` (`[::]`), and `PORT` (`80`). Values already supplied by the container environment take precedence. It then replaces the entrypoint shell with the server process and passes any additional arguments through.

## System service

The server can install a service for the detected system manager (systemd, OpenRC, runit, s6, SysVinit, launchd, Windows Service, or BSD rc.d):

```sh
wordList --service --install
wordList --service status
```

Use `--service start`, `stop`, `restart`, `reload`, or `disable` to manage it. `--service --uninstall` prompts before removing the service registration. System service installation requires administrator privileges; user-service installation is not currently available.
