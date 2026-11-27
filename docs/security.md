# Security

WordList has no end-user accounts. Public word-list and passphrase endpoints are unauthenticated. Anonymous owners can issue a primary token and manage linked tokens on `/server/preferences/tokens` or through the corresponding versioned API. The primary token is stored as a SHA-256 hash and an AES-256-GCM encrypted copy so it can be explicitly revealed; linked tokens have the same storage model. API token management accepts bearer credentials only. The browser token cookie is `HttpOnly`, `SameSite=Strict`, and scoped to token management forms. Passphrases are generated with Go's `crypto/rand`, returned to the caller, and not intentionally persisted by the application.

## Input and abuse controls

Passphrase count is bounded to 1–20 words, random-word count to 1–100, search query length to 100 characters, and request bodies to 1 MiB. The in-memory rate limiter defaults to 120 requests per 60 seconds and can be configured in `server.yml`. Tor requests use the exported circuit identity when available; I2P requests share an overlay-wide rate-limit bucket instead of treating the local tunnel address as a visitor IP.

## Response and browser protections

The server sets content-type sniffing, framing, referrer, permissions, and content security policy headers. CORS defaults to `*` without credentials. Passphrase responses disable caching and referrer transmission.

Health routes check the SQLite connection and loaded word list. A failed core check reports `unhealthy` with HTTP `503`; enabled Tor or I2P services that are unavailable produce `degraded` with HTTP `200`. Checks expose public-safe status values without database paths or error details. Overlay addresses are included only when the matching `expose` setting is enabled.

The consent banner uses a server-rendered, CSRF-protected POST and works without JavaScript. Consent state is stored in the `cookie_consent` cookie; the preferences form updates the same categories and updates the `ccpa_opt_out` cookie only when data sales are enabled.

## Public security endpoints

- `/.well-known/security.txt` is the canonical RFC 9116 contact file. On clearnet it publishes the configured private vulnerability reporting URL, an instance contact URL with a rolling identifier, the optional configured security email, the security policy URL, and a one-year expiry. On Tor it uses onion URLs and publishes only `server.tor.contact_email` when configured; it never falls back to clearnet reporting or contact addresses. The identifier is derived from a persistent 32-byte installation secret stored in SQLite. `/security.txt` returns `404`.
- `/server/healthz` and `/api/v1/server/healthz` expose service health and word count.
- `/server/security` and `/server/security/policy` publish security reporting information and policy. Tor requests receive onion URLs and the configured Tor contact email only.
- `/server/contact` renders abuse, DMCA, privacy, or security report forms with CSRF and arithmetic challenge controls. Security mode requires the current or previous 48-hour identifier and echoes a valid identifier only in the hidden revalidation field. Invalid identifiers fall back to ordinary report mode. Every form uses a signed review step. Field descriptors drive conditional, enum, and length validation; report text links must point to an allowed host, and name/subject fields reject links. A confirmed security report is encrypted at rest and returns a tracking ID. When SMTP is ready, the maintainer receives only the encrypted report body and the researcher receives a content-free status link. Missing SMTP does not prevent secure storage. Use GitHub private vulnerability reporting as the primary reporting channel.
- Security report notification templates are `security_report.txt`, `security_ack.txt`, and `security_status.txt` in the email template directory. Their subjects and body text come from the `email.*` locale keys; custom templates may rearrange translated text and permitted values, but researcher acknowledgment and status messages must not include report content.
- `/server/security/report/{tracking_id}?token=...` displays only the report's triage status, researcher-facing comment, and expected disclosure date. The SHA-256 token hash is stored; a valid token can read status once per UTC day and stops working 30 days after closure. The response is `no-store`, `no-referrer`, and `noindex`.
- `/server/security/thanks` lists credited researchers after disclosure, including the disclosure year and severity. `/api/v1/server/security` returns the same public entries in its `thanks` field. Hall-of-fame rows contain no contact information and outlive report retention.
- Closed reports are retained for `server.security.reports.retention_days` days (default 365). The non-disableable daily `security_cleanup` task deletes expired closed reports; open reports are retained.
- `/server/docs/swagger` publishes the OpenAPI documentation.

The configured `server.token` is accepted as the operator Bearer credential for MCP. OpenAPI operations marked `x-operatorOnly` are hidden from other MCP callers and rejected when called with a non-operator token; no operations currently carry that marker. Repository vulnerabilities should be reported through [GitHub's private vulnerability reporting](https://github.com/apimgr/wordList/security/advisories/new). The instance form is the secondary disclosure channel for deployed-instance findings.
