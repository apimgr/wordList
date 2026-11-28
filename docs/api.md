# API reference

The public word-list and passphrase REST API is rooted at `/api/v1` and requires no authentication. JSON success responses use `{"ok":true,"data":...}`. Errors use `{"ok":false,"error":"...","message":"..."}` with an optional `details` object when structured context is needed; the HTTP status carries the status code, and absent fields are omitted. A `.txt` route or `Accept: text/plain` returns plain text with a trailing newline. Owner-token management endpoints require a primary owner token in the `Authorization: Bearer` header. Cookie credentials are ignored by API routes.

## REST API

| Method | Path | Parameters and behavior |
|---|---|---|
| GET | `/api/v1` | Service overview and endpoint list; `/api/v1/` redirects here |
| GET | `/api/v1/random` or `/api/v1/random/{count}` | Return 6 random words by default; count 1–100 |
| GET | `/api/v1/random.txt` or `/api/v1/random/{count}.txt` | Same as above in plain text |
| GET | `/api/v1/passphrase` or `/api/v1/passphrase/{count}` | Secure passphrase; count 1–20, default 6; optional `separator` and `capitalize=false` |
| GET | `/api/v1/passphrase.txt` or `/api/v1/passphrase/{count}.txt` | Same passphrase as plain text |
| GET | `/api/v1/search/{query}` | Prefix search; `mode=contains` enables substring search; `limit` is 1–1000 |
| GET | `/api/v1/search/{query}.txt` | Search result words as plain text |
| GET | `/api/v1/words` | Paginated list; `page` defaults to 1 and `per_page` to 100 (maximum 1000) |
| GET | `/api/v1/words.txt` | Complete word list as plain text |
| GET | `/api/v1/dice/{code}` | Look up a five-die code |
| GET | `/api/v1/dice/{code}.txt` | Word for a dice code as plain text |
| GET | `/api/v1/letters/{letter}` | Words beginning with one ASCII letter |
| GET | `/api/v1/letters/{letter}.txt` | Matching words as plain text |
| GET | `/api/v1/lengths/{length}` | Words of the requested length (1–100) |
| GET | `/api/v1/lengths/{length}.txt` | Matching words as plain text |
| GET | `/api/v1/stats` | Dataset statistics |
| GET | `/api/v1/letters` | Available starting letters |
| GET | `/api/v1/count` | Word count as JSON |
| GET | `/api/v1/count.txt` | Word count as plain text |

### Owner token management

The primary token creates an anonymous owner and has `global` scope. A primary token can create and manage linked tokens; linked tokens cannot manage tokens. Linked tokens have immutable `read` or `read-write` scope, defaulting to `read-write`. A `read` token may read owner resources through REST and MCP but cannot mutate them; MCP also omits write tools. `global` is reserved for the primary token and is rejected when creating a linked token. Every token metadata object includes `scope`; raw token values are returned only by create, reveal, and rotate operations. Store each raw value securely when it is returned.

Create a linked token with an explicit scope when needed:

```json
{"name":"automation","scope":"read","expires":"30d"}
```

Omitting `scope` creates a `read-write` token. Scope cannot be changed later; revoke and recreate the linked token to change it.

| Method | Path | Purpose |
|---|---|---|
| POST | `/api/v1/server/preferences/tokens/primary` | Create an owner and primary token; returns `409` when the bearer token already identifies an owner |
| POST | `/api/v1/server/preferences/tokens/primary/reveal` | Reveal the caller's primary token |
| POST | `/api/v1/server/preferences/tokens/primary/rotate` | Replace the primary token without changing linked tokens |
| DELETE | `/api/v1/server/preferences/tokens/primary` | Delete the owner after last-four-character confirmation |
| GET | `/api/v1/server/preferences/tokens` | List masked token metadata |
| POST | `/api/v1/server/preferences/tokens` | Create a linked token |
| POST | `/api/v1/server/preferences/tokens/{id}/reveal` | Reveal a linked token |
| POST | `/api/v1/server/preferences/tokens/{id}/rotate` | Rotate a linked token |
| PATCH | `/api/v1/server/preferences/tokens/{id}` | Rename a linked token or change its expiry |
| DELETE | `/api/v1/server/preferences/tokens/{id}` | Revoke a linked token |

When `server.healthz.root.enabled` is `true`, `GET /healthz` mounts the same frontend health handler as `/server/healthz`. It negotiates HTML, text, or JSON without redirecting. The OpenAPI document, MCP tool list, and `llms.txt` endpoint catalog include the root alias only while this setting is enabled.

## Server routes

| Method | Path | Purpose |
|---|---|---|
| GET | `/api/v1/server/healthz` | Bare health payload with cached dependency checks, public-safe overlay state, restart/maintenance state, and persisted aggregate request statistics; unhealthy, maintenance, and shutdown states return `503` |
| GET | `/api/v1/server/auth/check` | Machine-only bearer-token validation; returns token kind without echoing the token |
| GET | `/api/v1/server/about` | About title and Markdown content derived from IDEA.md, or the language-resolved `server.branding.content.about` source; `format` identifies Markdown, HTML, or text |
| GET | `/api/v1/server/help` | Localized help content and available help pages |
| GET | `/api/v1/server/privacy` | Privacy summary, cookie descriptions, sharing, retention, third-party services, CCPA state, consent/usage text, and the language-resolved operator extension source/format |
| GET | `/api/v1/server/terms` | Localized built-in Markdown terms, or the language-resolved `server.branding.content.terms` source; `format` identifies Markdown, HTML, or text |
| GET | `/api/v1/server/help/mcp` | Resolved MCP URL, aliases, authentication scheme, live tool names, and WordList guidance |
| GET | `/api/v1/server/contact` | Field descriptors with conditional requirements, one-use form token, selected CAPTCHA provider and simple no-JavaScript fallback challenge, security URL, and public abuse email fallback. Tor and I2P requests always receive the simple provider |
| POST | `/api/v1/server/contact` | Two-step JSON submission: submit checks the single-use form token, one CAPTCHA proof, rate limits, and content rules, then returns a review token without delivery. Confirm repeats every business field and requires `confirmed: true`; rate limits count at confirm and duplicate reports are silently dropped. Abuse, DMCA, and privacy reports require working SMTP; security reports are stored encrypted and return a tracking ID. Security notifications and status email contain no plaintext report content |
| GET, POST | `/api/v1/server/preferences` | Read or partially update `theme`, `lang`, and `wordList_pref_{key}` values; POST reports `applied` and `ignored`, and requires accepted Preferences-cookie consent to write cookies |
| GET | `/api/v1/server/preferences/sync` | Export effective theme/language and accepted app preference cookies as a URL and short code; consent and credential cookies are excluded |
| GET | `/api/v1/server/preferences/sync/import` | Import from a `code` or direct query parameters; direct values win and missing Preferences-cookie consent blocks writes without an error |
| GET | `/api/v1/server/dpo` | DPO name and configured DPO/general email fallback; Tor requests use only `server.tor.contact_email` and an onion privacy URL |
| GET | `/api/v1/server/security` | Ordered public reporting contacts, expiry, PGP availability, request-resolved policy/thanks/contact URLs, and contact-free researcher acknowledgments in `thanks` |
| GET | `/api/v1/server/swagger` | OpenAPI JSON |
| GET, POST | `/api/v1/server/graphql` | Read-only GraphQL query endpoint; GET accepts `query`, optional `variables`, and optional `operationName` |
| POST | `/api/v1/server/mcp` | Stateless Streamable HTTP MCP endpoint (Bearer token required) |

The unversioned compatibility paths `/api/healthz`, `/api/swagger`, `/api/graphql`, `/api/mcp`, and `/mcp` are also available. `/api/autodiscover` returns the current API version, request origin and configured `DOMAIN` origins, plus the configured `cli_min_version`; `cli_versions` is empty until release binaries are provisioned in the server's distribution directory. It includes `onion` and `i2p` addresses only when their respective exposure settings are enabled. Requests to a Tor or I2P address return an empty `urls` list and include only the matching overlay address field, when that overlay is enabled and exposed.

`GET /cli/binaries/wordList-cli-{os}-{arch}[.exe]` streams an operator-provisioned client binary only when its SHA-256 matches the checksum manifest. Downloads are public unless `cli.binary_download.require_auth` is enabled, in which case a valid bearer token is required. See the [CLI guide](cli.md) for update behavior and the [configuration guide](configuration.md) for the distribution directory.

The privacy response uses the standard `{ "ok": true, "data": { ... } }` envelope. Its `ccpa.applicable` value follows `server.privacy.data.sold`, and `ccpa.user_opted_out` reflects the browser's `ccpa_opt_out` cookie. `content.operator_extension` and `content.operator_format` expose the language-resolved `server.branding.content.privacy` source; an empty value means the built-in policy is used without an operator extension. Generated passphrases are never stored.

All health routes render from one shared result. `/server/healthz` and the optional `/healthz` alias negotiate HTML, human-readable text for curl and other command-line clients, or JSON. The versioned and unversioned API routes return a bare JSON object by default; `.txt`, `Accept: text/plain`, and non-interactive API requests without an explicit format receive dot-notation machine text. API JSON/text use compact uptime and RFC 3339 UTC timestamps; frontend HTML/text localize labels and humanize durations, timestamps, and counts.

Health JSON fields follow a stable order: `project`, `status`, optional `pending_restart`, `restart_reason`, and `maintenance`, then version/build/runtime values, `features`, `checks`, and `stats`. Dataset checks are cached and refreshed by the scheduler at most every 60 seconds and after data updates; health requests do not scan datasets or query the database for statistics. The core database, cache, disk, and scheduler probes also run outside request handling. Health output contains only public-safe status values and aggregate counts. The status is `503` only for `unhealthy`, `maintenance`, or `shutting_down`; `degraded` and `restart_required` remain `200`.

## Model Context Protocol

MCP uses the same bearer credentials as the REST API. Create one linked token per MCP client on `/server/preferences/tokens`; token creation, reveal, rotation, and revocation are never MCP tools. The endpoint is stateless and accepts JSON-RPC over HTTP `POST` using protocol version `2026-07-28`.

For example, verify a connection by listing tools:

```sh
export SERVER_URL='https://your-public-host.example'
curl -q -LSsf "$SERVER_URL/api/v1/server/mcp" \
  -H "Authorization: Bearer $MCP_TOKEN" \
  -H 'Accept: application/json, text/event-stream' \
  -H 'Content-Type: application/json' \
  -H 'MCP-Protocol-Version: 2026-07-28' \
  -H 'Mcp-Method: tools/list' \
  -d '{"jsonrpc":"2.0","id":1,"method":"tools/list","_meta":{"io.modelcontextprotocol/protocolVersion":"2026-07-28"}}'
```

The root alias is `/mcp`; `/api/mcp` is the unversioned API alias. Browser `Origin` headers must match the request's resolved public origin or an explicit configured CORS origin. The wildcard CORS default does not authorize MCP origins.

## Examples

```sh
curl -q -LSsf "$SERVER_URL/api/v1/passphrase/8"
curl -q -LSsf "$SERVER_URL/api/v1/dice/11111"
curl -q -LSsf "$SERVER_URL/api/v1/search/pass?mode=prefix&limit=20"
curl -q -LSsf "$SERVER_URL/api/v1/words?page=2&per_page=100"
```

## GraphQL

The versioned endpoint and `/api/graphql` alias accept read-only queries over GET or POST. GET accepts `query`, optional JSON-object `variables`, and optional `operationName`; without `query`, it returns a usage message. POST accepts a JSON object containing `query` and optional `variables` and `operationName`.

The code-first schema supports `wordCount`, `stats`, paginated `words(page, perPage)`, `search(query, mode, limit)`, `wordByDiceCode(code)`, `wordsByLetter(letter)`, `wordsByLength(length)`, and `letters`. `stats` exposes total, length, letter, and length-frequency data. Queries are limited to 4,096 characters, 250 fields, and 12 nested field levels. The endpoint is read-only; GraphQL mutations are not configured. For example:

```sh
curl -q -LSsf --get --data-urlencode 'query={ wordCount stats { totalWords averageLength } words(page: 1, perPage: 5) { items total } }' \
  "$SERVER_URL/api/v1/server/graphql"
```

POST remains available for larger requests:

```sh
curl -q -LSsf -H 'Content-Type: application/json' \
  -d '{"query":"{ wordCount }"}' \
  "$SERVER_URL/api/v1/server/graphql"
```

Responses use the standard GraphQL shape (`{"data": {...}}`) without the REST `{ok, data}` envelope. Query parse and validation errors use `{"data": null, "errors": [...]}` with HTTP 200; malformed request bodies retain the API error shape and status.

The interactive GraphiQL page is `/server/docs/graphql`. It uses locally served, pinned browser assets and follows the site’s dark, light, and automatic theme. Without JavaScript, the page links to the plain GraphQL endpoint and retains a POST form.

The interactive Swagger UI is `/server/docs/swagger`; the OpenAPI JSON remains at `/api/v1/server/swagger` with `/api/swagger` as the current-version alias. Its browser assets are also served locally and follow the site theme.
