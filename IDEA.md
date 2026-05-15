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
- GraphQL at `/graphql`

**Non-goals:**
- No user accounts, registration, or login of any kind
- No admin web panel (server configured via `server.yml` only)
- No user-submitted words (EFF Diceware dataset only, updated with EFF releases)
- No definitions, parts of speech, or etymology
- No paid tiers, no API keys, no rate-limited access tiers
- No server-side passphrase storage (generated passphrases are never logged or retained)

### Roles & permissions

There are no user roles. All endpoints are public and require no authentication.

| Actor | Access |
|-------|--------|
| **Anonymous visitor (browser)** | Full read access to all web pages and API endpoints; passphrase generation |
| **Anonymous API client (curl/CLI)** | Full read access to all API endpoints; passphrase generation |
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

No PII stored or served. Generated passphrases are computed in memory, returned to the caller, and never logged or stored server-side.

### Trust boundaries & external services

| Boundary | Trust level | Notes |
|----------|-------------|-------|
| EFF Diceware word dataset (embedded at build) | Fully trusted | Static, compiled into binary |
| Incoming HTTP requests | **Untrusted** | Word count parameter validated and bounded |
| System CSPRNG (`crypto/rand`) | Fully trusted | Used for all passphrase generation; no `math/rand` |

No external services called at runtime.

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

- **No authentication on any endpoint**: intentional. Public read-only reference and generation API.
- **`crypto/rand` for passphrase generation**: mandatory. Using `math/rand` or any seeded PRNG for passphrase generation is a security defect; `crypto/rand` is the only acceptable source.
- **Passphrases never logged**: intentional. Logging a generated passphrase could expose it if a user uses it immediately after generation.
- **All responses include `Access-Control-Allow-Origin: *`**: intentional. Public API designed for cross-origin browser use.
