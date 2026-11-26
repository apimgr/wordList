# WordList

[![CI](https://github.com/apimgr/wordList/actions/workflows/ci.yml/badge.svg)](https://github.com/apimgr/wordList/actions/workflows/ci.yml)
[![Release](https://img.shields.io/github/v/release/apimgr/wordList)](https://github.com/apimgr/wordList/releases)
[![License](https://img.shields.io/github/license/apimgr/wordList)](LICENSE.md)

## About

WordList serves the EFF Large Wordlist for Diceware passphrase generation. The Go server embeds all 7,776 words and provides a versioned REST API, GraphQL endpoint, server-rendered web interface, and companion `wordList-cli` binary.

## Features

- Generate cryptographically secure passphrases with `crypto/rand` (six words by default, one to 20 words).
- Look up dice codes, retrieve random words, search and browse the list, and read word-list statistics.
- Use the responsive web interface without an account.
- Call JSON or plain-text REST endpoints, or the GraphQL endpoint.
- Explore REST endpoints in the locally hosted Swagger UI and run GraphQL queries in GraphiQL; both follow the site theme.
- Connect MCP clients to the stateless Streamable HTTP endpoint using an owner bearer token.
- Generate passphrases locally or from a configured server with `wordList-cli`.
- Run without runtime word-list downloads; the EFF data is embedded in the server binary.

## Production

### Docker Compose

Review the port and persistent volume locations in `docker/docker-compose.yml`, then start the service:

```sh
docker compose -f docker/docker-compose.yml up -d
```

The default host mapping is `172.17.0.1:64580`; change it to suit the host network. The Compose service mounts configuration and application data and checks `/api/v1/server/healthz`.

### Release binaries

Download the matching `wordList-{os}-{arch}` server and `wordList-cli-{os}-{arch}` client from [GitHub Releases](https://github.com/apimgr/wordList/releases). Windows artifacts use `.exe`. Verify the published checksums and place the binaries in a directory on `PATH`.

The server includes service management commands for supported operating systems. Use `wordList --service --install` (or `wordList --service install`) to install, enable, and start the service. Use `wordList --service start`, `stop`, `restart`, `reload`, or `status` to manage it. Use `wordList --service --disable` to stop and disable automatic startup. `--service --uninstall` prompts before removing the service registration.

## Client

The companion client can generate locally or call a configured server:

```sh
wordList-cli passphrase --words 8
wordList-cli passphrase --words 6 --separator ' ' --capitalize=false
```

Use `wordList-cli --help` for available commands. Remote calls accept `--server` and `--token`; configure the primary server with `WORDLIST_SERVER_PRIMARY` or the client configuration file. Pass `--capitalize=false` to keep generated words lowercase. Do not place an operator token in a shared or automated client configuration.

## Configuration

On first start, the server creates `server.yml` and a restricted `secrets.yml` in the platform-specific configuration directory. The generated operator token is printed once. There is no web-based server configuration interface.

| Setting | Default | Purpose |
|---|---|---|
| `server.address` | `[::]` | Dual-stack listen address |
| `server.port` | random `64000`–`64999` (`80` in an application container) | Listen port; `PORT` overrides it |
| `server.baseurl` | `/` | Optional URL path prefix |
| `server.mode` | `production` | Application mode |
| `server.rate_limit.enabled` | `true` | Enable request rate limiting |
| `server.tokens.max_linked` | `25` | Active linked-token limit per owner |
| `server.tokens.cookie_max_age` | `1y` | Sliding owner-cookie lifetime |
| `server.tokens.idle_purge_days` | `0` | Idle owner purge delay; zero disables it |
| `server.i18n.default_language` | `en` | Default locale |

The server reads `server.yml` and optional file-based area settings such as `yaml.d/pages.yml` from its configuration directory. In containers, `CONFIG_DIR` and `DATA_DIR` select configuration and data locations. See the [configuration guide](docs/configuration.md) for the current settings and proxy URL behavior.

## API

The canonical REST API base is `/api/v1`. Public word-list and passphrase routes do not require a token. Owner-token management routes require the specified bearer token; browser cookies are ignored by API routes.

| Endpoint | Purpose |
|---|---|
| `GET /api/v1/passphrase[/{count}]` | Generate a secure passphrase |
| `GET /api/v1/random[/{count}]` | Return random words |
| `GET /api/v1/search/{query}` | Search by prefix or substring |
| `GET /api/v1/words` | Browse the word list with pagination |
| `GET /api/v1/dice/{code}` | Resolve a dice code |
| `GET /api/v1/letters[/{letter}]` | List starting letters or matching words |
| `GET /api/v1/lengths/{length}` | List words of a selected length |
| `GET /api/v1/stats`, `/count` | Read word-list statistics |
| `GET /api/v1/server/healthz` | Read service health; SQLite or word-list failures return `503` with public-safe check values |
| `GET /api/v1/server/auth/check` | Validate a bearer token and return its kind without echoing the token |
| `GET /api/autodiscover` | Discover the API version, public origins, and available CLI release metadata |
| `GET /api/v1/server/swagger` | Retrieve the OpenAPI document |
| `GET /api/v1/server/help/mcp` | Read request-resolved MCP connection details and the live tool list |
| `GET, POST /api/v1/server/graphql` | Run a read-only GraphQL query |
| `POST /api/v1/server/mcp` | Use the API operations as authenticated MCP tools |
| `POST /api/v1/server/preferences/tokens/primary` | Issue an owner token |

Append `.txt` to the supported word-list and passphrase routes for plain-text responses. Passphrase endpoints accept `separator` (up to eight characters) and `capitalize=false`. JSON errors use the documented `ok`, `error`, and `message` fields.

```sh
export SERVER_URL='https://your-public-host.example'
curl -q -LSsf "$SERVER_URL/api/v1/passphrase/8"
curl -q -LSsf "$SERVER_URL/api/v1/dice/11111"
curl -q -LSsf "$SERVER_URL/api/v1/words?page=1&per_page=100"
```

The GraphQL endpoint is `/api/v1/server/graphql`; it provides read-only word-count, statistics, pagination, search, dice-code lookup, letter, and length queries. Its schema supports introspection and applies query-size, field-count, and nesting-depth limits. Consult the [API guide](docs/api.md) for details and examples.

The MCP endpoint is `/api/v1/server/mcp` (also `/api/mcp` and `/mcp`). It requires a bearer token and MCP protocol version `2026-07-28`. Create a linked token for each MCP client and see the [API guide](docs/api.md) for connection details.

## Other

The browser interface includes passphrase generation, search, About and Help pages, the configured privacy policy and DPO contact, terms, contact and security pages, preferences, and owner-token management. The first-visit consent banner works without JavaScript; when data sales are enabled, the privacy page offers a CSRF-protected CCPA opt-out. Privacy-policy deployment facts and consent copy are configured under `server.privacy`; see the [configuration guide](docs/configuration.md). Confirmed abuse, DMCA, and privacy submissions are emailed when SMTP is ready. Confirmed security reports are encrypted and stored for operator review; they are not emailed.

The embedded locale set includes English, Spanish, French, German, Italian, and Arabic. Arabic uses right-to-left layout. See the [translation guide](docs/translating.md).

The word data is the [EFF Large Wordlist](https://www.eff.org/files/2016/07/18/eff_large_wordlist.txt).

## Development

Development instructions are for contributors.

Go builds and tests run in the `casjaysdev/go:latest` container through the Makefile. Docker is required.

```sh
git clone https://github.com/apimgr/wordList.git
cd wordList
make dev
make format
make test
make lint
```

`make dev` builds the server and client into a temporary directory. `make build` cross-compiles the configured platform matrix. `make test` runs Go tests and vet with coverage enforcement; `make lint` runs vet, staticcheck, and golangci-lint. `./tests/run_tests.sh` runs the basic Docker binary and API checks. Read [Development](docs/development.md) before contributing.

## Disclaimer

This software is provided "as is" without warranty of any kind. Use it at your own risk.

- **No warranty:** The authors are not responsible for damages, data loss, or other issues arising from use of this software.
- **No professional advice:** This software does not provide legal, financial, medical, or other professional advice.
- **Third-party services:** The terms of connected external services apply separately.
- **Security:** No software can be guaranteed free of vulnerabilities.
- **Production use:** Evaluate the software and its configuration before deploying it in production.

By using this software, you acknowledge that you have read and understood this disclaimer.

## License

See [LICENSE.md](LICENSE.md) for the project license. The EFF word-list source has its own attribution and licensing terms.
