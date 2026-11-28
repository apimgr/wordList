# WordList documentation

WordList is a public Diceware word-list and passphrase service. This guide documents the browser interface, REST and GraphQL APIs, CLI, runtime configuration, security behavior, and development workflow.

## Quick start

Follow the [installation guide](installation.md) to start WordList with Docker Compose, a release binary, or a system service. To build from a source checkout, use `make dev`; it builds the server and companion CLI in the Go toolchain container. The default listener uses a random port from `64000` through `64999` outside an application container and port `80` inside one.

## Features

- The EFF Large Wordlist is embedded in the server binary.
- Passphrases use the operating system cryptographic random source.
- The REST API returns JSON by default and plain text on the documented `.txt` routes.
- The server-rendered browser pages work without JavaScript and include accessible preferences, contact, privacy, security, and health flows.
- MCP clients can call the API through the stateless Streamable HTTP endpoint with a linked bearer token.
- No account or API token is needed for public endpoints.
- Operators can enable local GeoIP country policies, IP/domain blocklists, Tor and I2P services, metrics, and automated backups.

## Guides

- [Installation](installation.md)
- [Configuration](configuration.md)
- [REST and GraphQL API](api.md)
- [CLI](cli.md)
- [Tor hidden service](tor.md)
- [Database](database.md)
- [Security](security.md)
- [Integrations and discovery](integrations.md)
- [IP and domain blocklists](blocklists.md)
- [GeoIP](geoip.md)
- [Metrics](metrics.md)
- [PGP key lifecycle](pgp.md)
- [TLS certificates](ssl.md)
- [Interactive API documentation assets](api-doc-assets.md)
- [Development](development.md)
- [Translating](translating.md)

## Links

- [Source repository](https://github.com/apimgr/wordList)
- [EFF Large Wordlist source](https://www.eff.org/files/2016/07/18/eff_large_wordlist.txt)

## License

WordList is distributed under the terms in [LICENSE.md](https://github.com/apimgr/wordList/blob/main/LICENSE.md).
