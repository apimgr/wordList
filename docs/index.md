# WordList documentation

WordList is a public Diceware word-list and passphrase service. This guide documents the browser interface, REST and GraphQL APIs, CLI, runtime configuration, security behavior, and development workflow.

## Quick start

Build the server and companion CLI in the Go toolchain container with `make dev`. The server listens on a random port from `64000` through `64999` by default (port `80` in an application container). Open `/` for the web interface or `/api/swagger` for the API reference.

## Features

- The EFF Large Wordlist is embedded in the server binary.
- Passphrases use the operating system cryptographic random source.
- The REST API returns JSON by default and plain text on the documented `.txt` routes.
- MCP clients can call the API through the stateless Streamable HTTP endpoint with a linked bearer token.
- The server-rendered generator works without JavaScript.
- No account or API token is needed for public endpoints.
- Optional country blocking uses local ip-location-db datasets and fails open when data is unavailable.

## Guides

- [Installation](installation.md)
- [Configuration](configuration.md)
- [REST and GraphQL API](api.md)
- [CLI](cli.md)
- [Security](security.md)
- [Integrations and discovery](integrations.md)
- [Development](development.md)
- [Translating](translating.md)

## Links

- [Source repository](https://github.com/apimgr/wordList)
- [EFF Large Wordlist source](https://www.eff.org/files/2016/07/18/eff_large_wordlist.txt)

## License

WordList is distributed under the terms in [LICENSE.md](https://github.com/apimgr/wordList/blob/main/LICENSE.md).
