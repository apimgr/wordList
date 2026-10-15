# TODO.AI.md — WordList

## Verified source and product work

- [x] Keep the word list embedded and validate that the loaded EFF dataset contains 7,776 sorted words.
- [x] Use the OS CSPRNG for generated words and passphrases; expose secure generation errors to HTTP callers.
- [x] Map dice codes from `11111` through `66666` to the corresponding ordered EFF words.
- [x] Add bounded word counts, request body size limits, request timeouts, security headers, and configurable per-IP rate limiting.
- [x] Add server-rendered service information pages, versioned API routes, direct health/OpenAPI/GraphQL aliases, dice lookup, and an embedded OpenAPI document.
- [x] Make passphrase generation and word search usable without JavaScript; remove passphrase history from browser storage and mark generated API responses no-store.
- [x] Add an installable PWA app shell with offline word-list caching and WebCrypto passphrase fallback.
- [x] Add `wordList-cli passphrase --words N` for local secure passphrase generation.
- [x] Replace the unsafe/release-mutating Makefile behavior and align builds to the `casjaysdev/go:latest` container and platform names.
- [x] Remove Tor from the runtime image, use a minimal non-root application process, simplify the entrypoint, fix Compose mounts, and add the development image expected by the Docker workflow.
- [x] Correct server health checks and web links to the canonical routes; update the existing route tests.
- [x] Add test coverage for secure generation, entropy, and the EFF dice mapping.
- [x] Document the EFF dataset source in the README.

## Remaining work

- [ ] Replace the current limited GraphQL field detector with a standards-compliant GraphQL parser/schema/resolver implementation. The endpoint currently accepts only a shallow subset and returns errors for unsupported fields; do not describe it as a complete GraphQL API until then.
- [ ] Expand `src/server/openapi.json` and the Swagger UI to document every implemented REST route, response shape, error, and query parameter. The embedded document currently covers the primary routes only.
- [ ] Complete the CLI contract from AI.md PART 32: remote `--server` operation, user-scoped config, JSON output, interactive TUI/setup wizard, and safe version/update workflow. The current CLI covers the IDEA's local passphrase command only.
- [ ] Add end-to-end offline PWA fallback coverage and finish no-JavaScript browsing by letter and word length. Search and server-side passphrase generation already have HTML fallbacks.
- [ ] Add source/license attribution for the EFF word list in `LICENSE.md` and `/server/terms` after confirming the applicable EFF dataset terms. The published source URL is documented; the source page's legal terms have not been resolved.
- [ ] Review remaining AI.md platform requirements outside IDEA's core scope (trusted-proxy resolution, full server CLI/config behavior, localization, and release/update behavior) and implement or document each applicable requirement.
- [ ] Primary session must run the repository's Docker-based test and lint/build gates after integration. These checks were intentionally not run in this subagent pass.
