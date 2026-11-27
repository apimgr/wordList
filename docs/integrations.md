# Integrations and discovery

WordList exposes a public REST API, GraphQL, and MCP endpoints. SMTP can be configured for contact notifications; there are no other outbound API integrations enabled by default. Configured trusted-proxy hostnames are resolved periodically through the system DNS resolver.

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
- `/api/v1/server/mcp` exposes the API operations as stateless Streamable HTTP tools; `/api/mcp` and `/mcp` are direct aliases. Every MCP request requires a Bearer token, and linked tokens are recommended for individual clients.
- `/api/v1/server/help/mcp` returns request-resolved MCP connection details and the public tool registry; `/server/help/mcp` provides the same guidance as a rendered page.
- `/.well-known/llms.txt` and `/llms.txt` publish the OpenAPI endpoint catalog for AI clients.
- `/robots.txt` provides crawler guidance. `/.well-known/security.txt` is the canonical RFC 9116 security contact file; `/security.txt` returns `404`.
- `/sitemap.xml` lists the currently allowed public HTML pages; it excludes the contact submission form and API routes.
- `/manifest.json` and `/sw.js` provide the web app manifest and service worker.

The service worker activates automatically on first install. When an updated worker is ready, the page offers **Update now** or **Later**; choosing **Update now** activates the waiting worker and reloads the page, with a five-second fallback if the browser does not report control handover. Update prompts follow the page's selected language.

Static assets and the manifest use cache-first delivery after their first successful fetch. API requests bypass the worker. Navigations try the network for up to eight seconds, then use a cached copy only when the response was marked public. `/offline.html` is currently the only route that emits the public-cache marker, and only when the request has no authorization or user-preference cookies. Other offline navigation uses a synthesized `503` page with translated text stored from the last visited page (English defaults are used before any page has stored them). Failed subresource requests return `504` responses.

The project currently exposes MCP over HTTP. OAuth/OIDC, webhooks, federation, and external identity integrations are not implemented.
