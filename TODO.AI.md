# Remaining AI.md alignment work

This list records verified implementation gaps that remain after synchronizing the project documentation and OpenAPI inventory. Completion entries must remain until the change that resolves them is committed, per AI.md PART 0.

## Priority order

1. Implement the mandatory server/API foundations and reconcile API route semantics, response envelopes, content negotiation, health output, security middleware, and project route inventory.
2. Implement the complete standard server page family and shared page shell, including preferences/help hubs, MCP, contact, DPO, and security pages with their required APIs and workflows.
3. Bring all rendered pages into PART 16 conformance: no-JS-first forms, accessible shared chrome, responsive theme tokens, PWA behavior, SEO, and removal of inline styles/scripts/handlers.
4. Finish i18n alignment: implement plural-message runtime/validation and translate the remaining user-facing strings across server, CLI, API, and email.
5. Complete the required client behavior and align its command/config/update lifecycle with AI.md PART 32.
6. Complete required root/project infrastructure and supporting docs/tests: provider workflows, Docker variants, service templates, integration and E2E scripts, page conformance, route coverage, and API/docs parity.
7. Verify the complete implementation using the project-required container commands and report actual results.

## Verified outstanding implementation gaps

These items were checked against AI.md and the current repository after the previous changes. They are implementation work, not audit conclusions.

### Server and data lifecycle

- Complete bearer-token validation and the owner-token lifecycle: primary and linked-token APIs and web forms use the database, and idle-owner deletion is wired to the cleanup task. Strict attempt limiting, audit events, scheduler persistence, and some response/auth edge cases remain missing.
- Implement the required MCP streamable HTTP endpoint and aliases, protocol/header/origin validation, OpenAPI-derived tools, and the live MCP help page/API. The current help page advertises an endpoint that is not registered.
- Finish request-aware URL resolution: current working tree validates hosts, gates forwarded host/proto/port/prefix headers on trusted peers, honors the `DOMAIN` allow-list and `server.baseurl`, resolves configured trusted-proxy DNS names every five minutes, and uses the resolver for absolute URLs, page links, redirects, manifest and service-worker scope. Still required: Tor/I2P runtime URL overrides, live domain learning, and parity audits for every URL-bearing route and client asset.
- Replace contact POST's 503 response with the required two-step validated submission, persistent encrypted reports, SMTP delivery, replay/rate controls, and security-report flow. Add the installation-secret security ID and live security acknowledgments instead of empty/static pages.
- Complete DPO, privacy, terms, security, preferences, and help page/API content from the specified config and product sources. The current DPO and legal pages are generic translations and the preferences/help registries omit required pages.
- Wire SQLite to application state and implement required app-secret initialization, audit, token, report, rate-limit, and scheduler persistence. Add the required libsql driver/configuration path and database-mode tests.
- Persist scheduler task state; implement the required registered jobs, catch-up and timezone behavior, and add panic/recovery coverage. The server now starts and stops the scheduler and registers the 15-minute token/captcha cleanup task, but scheduler state remains in memory only.
- Complete SMTP auto-detection, TLS modes, customizable embedded email templates, environment precedence, notification delivery, and email tests.
- Complete configured Let's Encrypt challenge support, renewal integration, certificate storage/security, and startup wiring. The current SSL package is not connected to the server and does not implement the required DNS provider behavior.
- Implement backups/restore and encryption/retention, GeoIP and blocklist loading, internal authenticated metrics, update checking/install lifecycle, privilege drop, and service manager templates/installation behavior.
- Implement Tor and opt-in/outside-container I2P setup, health monitoring, advertised addresses, and related help/footer behavior.

### Client, pages, docs, and release infrastructure

- Complete the client specification: CLI command/flag/config parity, localization and completions, token-aware API behavior, update flow, context detection, first-run wizard, and required TUI/GUI behavior.
- Finish PART 16/PART 30 page requirements: shared component conformance, complete responsive/accessibility checks, user-consent/preferences behavior, language parity for every rendered string, and offline/PWA install/update coverage.
- Add the required Incus, end-to-end, VM, and route/page conformance test runners; ensure generated OpenAPI, route registrations, MCP tools, llms.txt, IDEA inventory, and documentation stay in parity.
- Complete all required provider CI/release/security workflows and service/container artifacts from their applicable PARTs, then run the required containerized build, tests, linters, documentation checks, and CI verification.

## Completed in current working tree (not committed)

- Replaced README with a route- and behavior-focused guide.
- Added the required MkDocs/ReadTheDocs configuration and documentation pages for installation, configuration, REST/GraphQL, CLI, security, integrations, development, and translation status.
- Added the required root Claude loader, its AGENTS.md copy, Claude memory index, and PART-grouped Claude rule summaries.
- Removed the unsupported unversioned `/graphql` route and changed the GraphQL API routes to POST-only.
- Expanded the embedded OpenAPI document to cover the implemented JSON routes.
- Added the required IDEA.md project route inventory and removed its unsupported `/graphql` product claim.
- Added request-aware host/proto/port resolution, Host-header allow-listing, trusted-proxy configuration with DNS refresh, and path-prefix routing/link handling for server-rendered pages and PWA metadata; the startup banner uses the same host resolver. Tests cover spoof rejection, trusted proxy behavior, DNS cache expiry, custom/default ports, and prefixed routes.
- Added translation registry and embedded-locale consistency checks for file-set parity, locale metadata, nonempty values, message keys, and interpolation placeholders; added coverage for invalid `Accept-Language` quality values and corrected Italian CLDR plural metadata.
- Added `server.i18n.default_language` and `server.i18n.available_languages` configuration with required-locale validation; replaced the obsolete English-only translator placeholder with a current contribution guide and documented the locale settings.
- Added database primitives to issue an anonymous owner with a SHA-256-indexed primary token, store an AES-256-GCM encrypted token copy, authenticate unexpired tokens, touch use timestamps, and reveal a token for an authorized owner action. Token rows now use UUID identifiers, with a preserving migration from the previous integer-ID schema. Primary and linked-token HTTP/API integration is implemented and tested; CLI token integration remains outstanding.
- Generated a 32-byte AES encryption key into owner-only `secrets.yml` with version metadata and verified that restarts preserve it. The key is used by the primary and linked token routes; key rotation, migration of operator-set secrets, and audit events remain outstanding.
- Added the primary owner-token page with new-token, restore, one-use reveal, rotate, revoke, and confirmed forget flows; API issue/reveal/rotate/delete endpoints use bearer headers only. Primary tokens are encrypted, cookie values are HttpOnly/SameSite=Strict with the configured sliding expiry (one year by default), and tests cover issuance, one-use display, rotation, revocation, and API cookie isolation. Strict attempt limiting, audit, and remaining lifecycle rules are still outstanding.
- Added `server.tokens.max_linked`, `cookie_max_age`, and `idle_purge_days` configuration with validation; the configured sliding cookie expiry, linked-token cap, and periodic idle-owner deletion are applied.
- Added masked token listing and linked-token create/reveal/rotate/rename/expiry/revoke API routes and complete linked-token SSR controls, primary-token-only bearer authorization, active linked-token limits, date/relative expiry semantics, UUID identifiers, localized status and action notices, responsive layouts, OpenAPI/README/IDEA inventory, and API/database/web lifecycle tests.
- Started and stopped the built-in scheduler with the server and registered periodic cleanup for expired contact challenges and optionally idle owners; persistent state, recovery/catch-up, timezone behavior, and the other required tasks remain outstanding.
- Fixed the token screen route being replaced by the generic hub handler, aligned GraphQL response/error shapes and root-field selection, preserved valid saved CLI server configuration, updated the README structure and capitalization option, and moved Docker binary-test output into a private temporary directory.
