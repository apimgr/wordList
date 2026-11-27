# Development

## Prerequisites

Docker is required for Go builds and tests. The project uses `casjaysdev/go:latest` and does not require a Go toolchain on the host.

## Build

```sh
make dev
make local
make build
```

`make dev` places both binaries in a temporary output directory. `make local` writes binaries under `binaries/`; `make build` cross-compiles the platform matrix.

`make docker` builds and pushes the multi-architecture image, then checks its required runtime tools and Tor health contract. To check a separately built or published image, run `./tests/image-contract.sh IMAGE`; Docker, curl, and jq must be available on the host. The same check runs in each CI provider after image builds.

## Run and test

Run the server binary produced by `make dev`. Run `make format` to apply `gofmt` to Go source and tests. `make test` runs the CI quality checks (Go vet, staticcheck, govulncheck, ShellCheck, and shfmt), all Go package tests with the required 60% coverage gate, and the basic compiled-binary/API checks in `tests/run_tests.sh`. Run `make lint` for vet, staticcheck, and golangci-lint. All Go tooling runs in `casjaysdev/go:latest`; no host Go installation is required.

Shared presentation helpers in `src/common/format` render localized durations, binary sizes, and grouped counts, plus user-facing timestamps with a zero-padded day and timezone. Human-readable pages use these helpers; machine-readable API responses preserve raw values.

Server-rendered pages detect Lynx, w3m, Links, ELinks, Browsh, Carbonyl, and NetSurf user agents and serve a separate no-JavaScript layout with ordinary navigation and form submissions. Responses vary on `User-Agent` so shared caches do not mix the browser layouts.

The basic binary and HTTP API checks run in Docker:

```sh
./tests/run_tests.sh
```

This invokes `tests/docker.sh`, which runs `make dev`, copies the server and client into a disposable Alpine container, checks their help/version output, starts the server, and checks health, text, and paginated word-list endpoints. The script removes its container and temporary files when it exits. Run `tests/run_tests.sh e2e` (or `make test-e2e`) for the Chromium browser tier and `tests/run_tests.sh incus` (or `make test-incus`) for systemd/OpenRC service checks. Pass `INCUS_MODE=systemd` or `INCUS_MODE=openrc` to select one init system.

## Project layout

- `src/server/` contains the HTTP server, templates, static assets, and OpenAPI document.
- `src/words/` contains word lookup and secure passphrase generation.
- `src/config/` and `src/paths/` manage configuration and OS-specific paths.
- `src/client/` contains the companion CLI.
- `docker/` contains Docker and Compose definitions.

## Code changes

Keep API documentation, rendered pages, README, and this documentation aligned with the routes and behavior actually implemented. Do not log generated passphrases or use a non-cryptographic random source for them.
