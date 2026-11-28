# Integrations and discovery

WordList exposes a public REST API, GraphQL, and MCP endpoints. Outbound integrations include SMTP notifications, release checks, GeoIP/blocklist data downloads, Let's Encrypt ACME, and optional public-key publication to configured OpenPGP keyservers. Configured trusted-proxy hostnames are resolved periodically through the system DNS resolver.

## Server and CLI discovery

`GET /api/autodiscover` returns JSON containing the service name and version, current API version, usable service URLs, published CLI platform versions and SHA-256 digests, and the configured minimum CLI version. Overlay addresses are included only when the corresponding `server.tor.expose` or `server.i2p.expose` setting is enabled. Requests received over an overlay do not disclose clearnet service URLs.

The CLI updater can use operator-provisioned artifacts in `cli.binary_download.dir` (default `{data_dir}/cli`). Publish a `version.txt`, `sha256.txt`, and platform-named CLI binaries there; only binaries with matching checksums are advertised. Downloads use `/api/{api_version}/cli/binaries/{platform-binary-name}`. Set `cli.binary_download.require_auth: true` to require a bearer token for downloads. Do not put untrusted binaries in this directory.

The service publishes its API catalog at `/.well-known/llms.txt` and `/llms.txt`, OpenAPI JSON at `/api/swagger`, and GraphQL at `/api/graphql`. The versioned API forms are `/api/v1/server/swagger` and `/api/v1/server/graphql`. MCP is stateless Streamable HTTP at `/api/v1/server/mcp`, with direct aliases `/api/mcp` and `/mcp`; each request requires a bearer token and an allowed Origin when an Origin header is sent. Use a linked owner token for an individual MCP client. The operator token exposes operator-only tools only when an operation is explicitly marked for operators.

WordList does not provide OAuth/OIDC, webhooks, federated identity, or external identity-provider login.

## Tor hidden service

When a Tor executable is available, the server starts a managed v3 onion hidden service. Tor bootstrap runs in the background so clearnet service remains available. The onion hostname is persisted in the application data directory; `/api/autodiscover` and health output omit it unless `server.tor.expose` is enabled. Onion requests resolve absolute URLs to `http://{onion}`, keep Secure cookies, and use the configured `server.tor.contact_email` on the security and DPO/contact surfaces. If that address is empty, clearnet contact addresses are not substituted. Rendered onion pages remove clearnet and third-party link targets while retaining same-origin onion links. Clearnet HTML pages advertise the onion URL with `Onion-Location`. Tor-routed outbound HTTP clients fail when the Tor SOCKS service is unavailable rather than falling back to a direct connection.

Set `server.tor.use_network: true` to enable the managed SOCKS listener for outbound requests. The server does not act as a Tor relay or exit node. The hidden service uses a separate loopback backend listener and does not share the clearnet listener.

## I2P eepsite

I2P support is included in the server. Activation defaults to off on bare metal and in virtual machines, and on in an application container. When enabled, the server prefers to start and manage its own dedicated `i2pd` process. If no `i2pd` binary is available, it tries the configured external SAM bridge (`127.0.0.1:7656` by default). With neither provider available, health reports the eepsite as `disabled` while clearnet service continues; the manager retries provider discovery. The destination key is persisted under the application data directory so the `.b32.i2p` address survives restarts.

The address is omitted from autodiscovery and its hostname field is empty in health output unless `server.i2p.expose` is enabled. I2P requests use plain HTTP and the configured virtual port; they do not set Secure cookies. I2P does not provide outbound request proxying, and the I2P backend uses a separate loopback listener without Tor's PROXY-protocol header. Settings can be changed through `server.i2p`, `I2P_ENABLED`, or `--i2p[=BOOL]`; supported I2P settings reload dynamically.

## Discovery and documentation

- `/api/v1` returns the service summary and key endpoint paths; `/api/v1/` redirects to the canonical root.
- `/api/v1/server/swagger` returns the OpenAPI document; `/api/swagger` is an unversioned alias.
- `/api/v1/server/graphql` accepts GraphQL requests; `/api/graphql` is its unversioned alias.
- `/api/v1/server/mcp` exposes the API operations as stateless Streamable HTTP tools; `/api/mcp` and `/mcp` are direct aliases. Every MCP request requires a Bearer token, and linked tokens are recommended for individual clients. Use `read` scope when the client only needs lookup tools; write tools are omitted for that scope. Linked tokens default to `read-write`.
- `/api/v1/server/help/mcp` returns request-resolved MCP connection details and the public tool registry; `/server/help/mcp` provides the same guidance as a rendered page.
- `/.well-known/llms.txt` and `/llms.txt` publish the OpenAPI endpoint catalog for AI clients.
- `/robots.txt` provides crawler guidance. `/.well-known/security.txt` is the canonical RFC 9116 security contact file; `/security.txt` returns `404`.
- `/sitemap.xml` lists the currently allowed public HTML pages; it excludes the contact submission form and API routes.
- `/manifest.json` and `/sw.js` provide the web app manifest and service worker.
- `/.well-known/pgp-key.asc` publishes the active security-report encryption key while one exists. OpenPGP Web Key Directory is available at `/.well-known/openpgpkey/hu/{hash}` and its policy endpoint only when the key exists, a valid security contact is configured, and the request host matches the contact domain.

The service's well-known namespace is intentionally limited to the paths documented here. Unknown `/.well-known/*` paths return `404`; WordList does not enable WebFinger, OpenID provider metadata, Android App Links, Apple association, or MTA-STS.

The service worker activates automatically on first install. When an updated worker is ready, the page offers **Update now** or **Later**; choosing **Update now** activates the waiting worker and reloads the page, with a five-second fallback if the browser does not report control handover. Update prompts follow the page's selected language.

Static assets and the manifest use cache-first delivery after their first successful fetch. API requests bypass the worker. Navigations try the network for up to eight seconds, then use a cached copy only when the response was marked public. Public read-only pages emit the marker only when they contain no per-visitor state; the worker keys their copies by theme, language and consent. Other offline navigation uses a synthesized `503` page with translated text stored from the last visited page (English defaults are used before any page has stored them). There is no static offline route. Failed subresource requests return `504` responses.

The project currently exposes MCP over HTTP. OAuth/OIDC, webhooks, federation, and external identity integrations are not implemented.
