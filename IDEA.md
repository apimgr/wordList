## Project description

WordList is a full-stack Go web application serving the EFF Large Wordlist for Diceware passphrases through a versioned REST API, GraphQL endpoint, and a server-side rendered web UI. It provides the complete 7,776-word EFF Diceware list, cryptographically secure passphrase generation, and word lookup by dice roll code. All word data is embedded in the binary at build time. A companion CLI tool enables passphrase generation directly from the terminal. Deployed as a single self-contained static binary.

## Project variables

project_name: wordList
project_org: apimgr
internal_name: wordList
internal_org: apimgr
app_name: WordList API
repo: https://github.com/apimgr/wordList
license: MIT
binary: wordList
client_binary: wordList-cli
owner_token: wordList_owner_token_qwhzkfta

## Business logic

### Product scope & non-goals

**In scope:**
- EFF Large Wordlist: 7,776 words, each mapped to a 5-digit dice roll code (11111–66666)
- Cryptographically secure passphrase generation using `crypto/rand` (default: 6 words, ~77.5 bits of entropy)
- Configurable word count for passphrase generation
- Word lookup by 5-digit dice roll code
- Random word retrieval
- Full word list retrieval as JSON
- Word count (total words in dataset)
- Full web frontend (server-side Go templates, dark/light/auto theme, PWA, mobile-first) with interactive passphrase generator
- Server pages: `/server/about`, `/server/help`, `/server/healthz`, `/server/privacy`, `/server/terms`
- CLI client (`wordList-cli`) for terminal passphrase generation: `wordList-cli passphrase --words 8`
- OpenAPI/Swagger docs at `/api/{api_version}/server/swagger`
- GraphQL at `/api/{api_version}/server/graphql` (with `/api/graphql` alias)

**Non-goals:**
- No user accounts, registration, or login of any kind
- No admin web panel (server configured via `server.yml` only)
- No user-submitted words (EFF Diceware dataset only, updated with EFF releases)
- No definitions, parts of speech, or etymology
- No paid tiers or API keys for the public word-list and passphrase features; required AI.md operator and owner token flows still apply to system functions
- No server-side passphrase storage (generated passphrases are never logged or retained)

### Roles & permissions

There are no end-user accounts or product roles. Public word-list and passphrase features require no authentication. Mandatory AI.md system functions use their specified operator or owner tokens; this does not add user registration or paid access tiers.

| Actor | Access |
|-------|--------|
| **Anonymous visitor (browser)** | Read public pages and use public word-list and passphrase features |
| **Anonymous API client (curl/CLI)** | Read public word-list data and generate passphrases without an API key |
| **Server operator** | Configures server via `server.yml` only; no web management interface |

### Data model & sensitivity

**Word record** (embedded at build time, no PII):

| Field | Type | Sensitivity |
|-------|------|-------------|
| `dice` | string — 5-digit dice roll code (e.g., `24146`) | Public |
| `word` | string — the corresponding Diceware word (e.g., `entail`) | Public |

**Passphrase response** (generated per-request, never stored):

| Field | Type | Sensitivity |
|-------|------|-------------|
| `passphrase` | string — space-separated words | **Sensitive** — never logged |
| `words` | string[] — individual words | **Sensitive** — never logged |
| `entropy` | float — entropy in bits | Public |

The passphrase and word-list features do not store PII. Generated passphrases are computed in memory, returned to the caller, and never logged or stored server-side. Mandatory AI.md system features may store operational records, owner/token metadata, and encrypted contact or security reports as specified there; that data is separate from passphrase generation and must follow the required retention, encryption, and privacy controls.

### Trust boundaries & external services

| Boundary | Trust level | Notes |
|----------|-------------|-------|
| EFF Diceware word dataset (embedded at build) | Fully trusted | Static, compiled into binary |
| Incoming HTTP requests | **Untrusted** | Word count parameter validated and bounded |
| System CSPRNG (`crypto/rand`) | Fully trusted | Used for all passphrase generation; no `math/rand` |

The passphrase and word-list features do not require external services. Mandatory AI.md capabilities may make outbound SMTP, Tor, I2P, GeoIP, update, and optional remote database connections when configured or enabled; each remains disabled or unavailable when its required configuration is absent.

### Threat model & abuse cases

**Primary assets:** service availability; passphrase entropy integrity (ensuring generated passphrases are truly random).

**Attacker/abuser goals:**
- DoS via high-rate passphrase generation requests
- Bulk dump of the full word list
- Undermining passphrase randomness (requires compromising the OS CSPRNG — outside scope)

**Defenses:**
- Rate limiting on all endpoints
- Passphrase word count capped at a configurable maximum (default: 20 words) to prevent excessive entropy computation
- Passphrase and individual words are never logged — only path and status are recorded
- No user accounts eliminates credential stuffing and privilege escalation entirely

### Security decisions & exceptions

- **No end-user accounts or login**: intentional. Public word-list and passphrase routes remain open; AI.md-required operator and anonymous owner tokens apply only to system features and their API routes.
- **`crypto/rand` for passphrase generation**: mandatory. Using `math/rand` or any seeded PRNG for passphrase generation is a security defect; `crypto/rand` is the only acceptable source.
- **Passphrases never logged**: intentional. Logging a generated passphrase could expose it if a user uses it immediately after generation.
- **CORS follows AI.md configuration and origin resolution**: the public API is designed for cross-origin use, while credentials, MCP origin checks, and configured allow-lists follow the mandatory server security rules.

### API routes

Public word-list and passphrase routes are unauthenticated; mandatory AI.md system routes use their specified token requirements. The native API uses the `/api/v1` prefix. JSON endpoints return the documented response object; `.txt` endpoints return plain text. Requests that generate a passphrase accept a word count, optional separator, and optional capitalization setting. The server-rendered pages are listed separately because their HTML responses are the browser interface for these features.

| Method | Path | Auth | Purpose and request/response summary |
|---|---|---|---|
| GET | `/api/v1/` | None | API summary, version, word count, and endpoint list |
| GET | `/api/v1/passphrase` | None | Generate a secure passphrase; `count` query defaults to 6 and must be 1–20; returns passphrase, words, entropy, and options |
| GET | `/api/v1/passphrase/{count}` | None | Generate a secure passphrase with 1–20 words; returns passphrase, words, entropy, and options |
| GET | `/api/v1/passphrase.txt` | None | Generate a six-word passphrase as plain text |
| GET | `/api/v1/passphrase/{count}.txt` | None | Generate a 1–20 word passphrase as plain text |
| GET | `/api/v1/random` | None | Return 6 random words by default; query count must be 1–100 |
| GET | `/api/v1/random/{count}` | None | Return the requested 1–100 random words |
| GET | `/api/v1/random.txt` | None | Return 6 random words as plain text by default |
| GET | `/api/v1/random/{count}.txt` | None | Return the requested 1–100 random words as plain text |
| GET | `/api/v1/search/{query}` | None | Search by prefix or `mode=contains`; optional `limit` from 1 to 1000; returns matching words and metadata |
| GET | `/api/v1/search/{query}.txt` | None | Return matching words as plain text |
| GET | `/api/v1/words` | None | Return a page of words; `page` defaults to 1, `per_page` defaults to 100 and is capped at 1000 |
| GET | `/api/v1/words.txt` | None | Return the complete word list as plain text |
| GET | `/api/v1/dice/{code}` | None | Resolve a five-digit dice code to its word |
| GET | `/api/v1/dice/{code}.txt` | None | Return the matching word as plain text |
| GET | `/api/v1/letters/{letter}` | None | Return words beginning with one ASCII letter |
| GET | `/api/v1/letters/{letter}.txt` | None | Return those words as plain text |
| GET | `/api/v1/lengths/{length}` | None | Return words with the requested length |
| GET | `/api/v1/lengths/{length}.txt` | None | Return those words as plain text |
| GET | `/api/v1/stats` | None | Return word-list statistics |
| GET | `/api/v1/letters` | None | Return distinct starting letters |
| GET | `/api/v1/count` | None | Return the word count as JSON |
| GET | `/api/v1/count.txt` | None | Return the word count as plain text |
| GET | `/api/v1/server/healthz` | Per AI.md | Return the public-safe service health response |
| GET | `/api/v1/server/about` | None | Return About page data as JSON |
| GET | `/api/v1/server/help` | None | Return Help page data as JSON |
| GET | `/api/v1/server/privacy` | None | Return Privacy page data as JSON |
| GET | `/api/v1/server/terms` | None | Return Terms page data as JSON |
| GET | `/api/v1/server/contact` | None | Return contact form types, field descriptors, and anti-spam requirements |
| POST | `/api/v1/server/contact` | Per AI.md | Submit or confirm a validated contact or security report |
| GET, POST | `/api/v1/server/preferences` | Per AI.md | Read or save browser preferences |
| POST | `/api/v1/server/preferences/tokens/primary` | None; valid bearer token returns conflict | Issue a new anonymous owner and primary token; returns `owner_token` once |
| GET | `/api/v1/server/preferences/sync` | None | Export supported preferences as a URL and code |
| GET | `/api/v1/server/preferences/sync/import` | None | Import supported preferences from a URL |
| GET | `/api/v1/server/dpo` | None | Return data-protection contact information |
| GET | `/api/v1/server/security` | None | Return public security policy and reporting information |
| POST | `/api/v1/server/mcp` | Bearer token | Streamable HTTP MCP endpoint; exposes the existing API operations as tools |
| GET | `/api/v1/server/swagger` | None | Return the OpenAPI document |
| POST | `/api/v1/server/graphql` | Per AI.md | Accept a JSON GraphQL query and return supported word count, statistics, or word-list fields |
| GET | `/api/healthz` | None | Unversioned health alias |
| GET | `/api/swagger` | None | Unversioned OpenAPI alias |
| POST | `/api/graphql` | None | Unversioned GraphQL alias |
| POST | `/api/mcp` | Bearer token | Unversioned alias for the MCP endpoint |
| POST | `/mcp` | Bearer token | Root alias for the MCP endpoint |

The browser surface includes `/`, `/generate`, `/search`, `/server/about`, `/server/help`, `/server/help/mcp`, `/server/privacy`, `/server/terms`, `/server/contact`, `/server/dpo`, `/server/security`, `/server/security/policy`, `/server/security/thanks`, `/server/preferences`, `/server/preferences/tokens`, `/server/preferences/sync`, `/server/healthz`, `/server/docs/swagger`, and `/server/docs/graphql`. `/graphql` is not a supported route.
