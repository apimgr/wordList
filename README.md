# WordList

WordList serves the EFF Large Wordlist for cryptographically secure Diceware passphrase generation. It embeds the 7,776-word list, maps words to five-die codes, and provides a public versioned REST API, GraphQL endpoint, server-rendered web app, and companion CLI.

## Features

- Secure passphrase generation using the operating system CSPRNG (six words by default, up to 20)
- Dice code lookup, random word selection, search, browsing, and word list statistics
- Server rendered, mobile friendly web interface with installable offline app shell
- OpenAPI JSON, GraphQL endpoint, and `wordList-cli passphrase --words 8`
- No accounts, external runtime dependencies, or server-side passphrase retention

## Build

Builds use the `casjaysdev/go:latest` container through the Makefile. The default listen port is `64851`.

```sh
make dev
```

The configuration file is created at the platform-specific WordList config directory on first start. In containers, configure `CONFIG_DIR`, `DATA_DIR`, and `PORT`.

## Web pages

- `/` — overview
- `/generate` — passphrase generator
- `/search` — word search
- `/server/about`, `/server/help`, `/server/privacy`, `/server/terms`
- `/server/healthz`
- `/server/docs/swagger` and `/server/docs/graphql`

## API

API v1 is rooted at `/api/v1`; the OpenAPI document is `/api/v1/server/swagger`, health status is `/api/v1/server/healthz`, and GraphQL accepts POST requests at `/api/v1/server/graphql` (`/api/graphql` is an alias).

```sh
curl 'http://localhost:64851/api/v1/passphrase/8'
curl 'http://localhost:64851/api/v1/dice/11111'
curl 'http://localhost:64851/api/v1/words?page=1&per_page=100'
```

## CLI

```sh
wordList-cli passphrase --words 8
```

Generated passphrases are not logged, cached, or saved by the server or browser. The EFF Large Wordlist is sourced from [EFF's published list](https://www.eff.org/files/2016/07/18/eff_large_wordlist.txt).
