# WordList

WordList serves the EFF Large Wordlist for Diceware passphrase generation. The Go server embeds the word data and provides a public REST API, GraphQL endpoint, server-rendered web pages, and a companion `wordList-cli` binary.

## Features

- Generate passphrases with Go's operating-system cryptographic random source (six words by default, one to 20 words).
- Look up dice codes, retrieve random words, search and browse the word list, and read word-list statistics.
- Use the responsive web interface without an account or API token.
- Call JSON or plain-text REST endpoints, or the GraphQL endpoint.
- Generate passphrases from a terminal with `wordList-cli`.
- Run without runtime data downloads; the EFF word list is embedded in the server binary.

Generated passphrases are sensitive. The server does not store them, and the API marks generation responses `Cache-Control: no-store` and `Referrer-Policy: no-referrer`.

## Build and run

Go builds and tests run in the `casjaysdev/go:latest` container through the Makefile. `make dev` builds the server and client into a temporary directory and prints its location. The default server port is `64851`.

Configuration is created on first start in the platform-specific configuration directory. In Docker, `CONFIG_DIR`, `DATA_DIR`, and `PORT` select the runtime paths and port.

## Web pages

| Path | Purpose |
|---|---|
| `/` | Word-list overview and links to search and generation |
| `/generate` | No-JavaScript passphrase form |
| `/search` | Search and browse the list |
| `/server/about` | Service information |
| `/server/help` | Usage help |
| `/server/help/mcp` | MCP connection guidance |
| `/server/privacy` | Privacy information |
| `/server/terms` | Terms of use |
| `/server/contact` | Abuse, DMCA, and privacy report forms (submission currently unavailable) |
| `/server/dpo` | Configured DPO contact and privacy link |
| `/server/security` | Security reporting and policy links |
| `/server/security/policy` | Disclosure policy |
| `/server/security/thanks` | Researcher acknowledgments |
| `/server/preferences` | Browser preferences |
| `/server/preferences/sync` | Export and import browser preferences |
| `/server/healthz` | Human-readable service status |
| `/server/docs/swagger` | REST API reference |
| `/server/docs/graphql` | GraphQL endpoint information |

## REST API

The versioned API root is `/api/v1`. Endpoints are public and do not require authentication.

| Method and path | Purpose |
|---|---|
| `GET /api/v1/` | API overview |
| `GET /api/v1/server/healthz` | JSON health status |
| `GET /api/v1/server/about`, `/help`, `/privacy`, `/terms` | Server page content as JSON |
| `GET /api/v1/server/dpo` | Configured DPO/general contact fallback and privacy URL |
| `GET /api/v1/server/contact` | Contact form fields and anti-spam requirements |
| `GET, POST /api/v1/server/preferences` | Read or save browser preferences |
| `GET /api/v1/server/preferences/sync` | Export preferences |
| `GET /api/v1/server/preferences/sync/import` | Import preferences |
| `GET /api/v1/server/swagger` | OpenAPI document |
| `POST /api/v1/server/graphql` | GraphQL endpoint |
| `GET /api/v1/random[/{count}]` | Random words (default 6, 1–100) |
| `GET /api/v1/random[/{count}].txt` | Random words as plain text |
| `GET /api/v1/passphrase[/{count}]` | Secure passphrase (default 6, 1–20) |
| `GET /api/v1/passphrase[/{count}].txt` | Passphrase as plain text |
| `GET /api/v1/search/{query}` | Search by prefix or `?mode=contains`; optional `limit` (1–1000) |
| `GET /api/v1/words` | Paginated list (`page`, `per_page`, default 100) |
| `GET /api/v1/dice/{code}` | Word for a five-digit dice code |
| `GET /api/v1/letters/{letter}` | Words beginning with a letter |
| `GET /api/v1/lengths/{length}` | Words of a specified length |
| `GET /api/v1/stats`, `/letters`, `/count` | Dataset statistics, letters, and count |

Append `.txt` where listed for plain-text responses. Passphrase endpoints accept `separator` (up to eight characters) and `capitalize` (`false` disables capitalization). JSON errors use the documented `ok`, `error`, and `message` fields.

```sh
curl -q -LSsf http://localhost:64851/api/v1/passphrase/8
curl -q -LSsf http://localhost:64851/api/v1/dice/11111
curl -q -LSsf 'http://localhost:64851/api/v1/words?page=1&per_page=100'
```

## GraphQL

The GraphQL endpoint is `/api/v1/server/graphql`; the implementation currently supports word count, word-list statistics, and the full word list. Submit a JSON object with a `query` field. The API documentation page describes the endpoint.

## CLI

```sh
wordList-cli passphrase --words 8
wordList-cli passphrase --words 6 --separator ' '
```

The CLI accepts `--words` from 1 to 20, `--separator` up to eight characters, and `--lowercase` to disable capitalization. `--help` and `--version` are also available.

## Docker

The Compose deployment is defined in `docker/docker-compose.yml`. The service publishes the configured container port and stores configuration and data in mounted directories. Its health check calls `/api/v1/server/healthz`. For automated binary and API checks, run `./tests/run_tests.sh`; it requires Docker and builds the server and CLI before running them in a disposable Alpine container. `make test` runs the Go test, vet, and coverage suite.

## Source and attribution

The word list is the [EFF Large Wordlist](https://www.eff.org/files/2016/07/18/eff_large_wordlist.txt). This project is distributed under the license in [LICENSE.md](LICENSE.md).
