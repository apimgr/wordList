# Tor hidden service

When a Tor binary is available, the server starts and owns a v3 onion service. If Tor is unavailable, the clearnet service continues to run. `server.tor.use_network` controls whether server-originated HTTP requests use Tor; it is separate from serving the onion site. Onion addresses are published on public surfaces only when `server.tor.expose` is enabled. Set `server.tor.contact_email` to use an onion-specific contact address.

The server stores Tor configuration and service keys in its configured directories and enforces private permissions on them. The onion hostname is persisted under `DATA_DIR/tor/site/hostname`. Private keys are never returned by status commands or control endpoints.

## Local maintenance commands

The server binary provides Tor status and maintenance commands:

```sh
wordList --maintenance tor status
wordList --maintenance tor validate
wordList --maintenance tor restart
wordList --maintenance tor regenerate [--yes]
wordList --maintenance tor vanity start --prefix abc [--workers 4]
wordList --maintenance tor vanity stop
wordList --maintenance tor vanity apply [ADDRESS] [--yes]
wordList --maintenance tor import-keys /absolute/path/to/key-directory [--yes]
```

Status and validation can run without a live server. The other commands require the running server because it owns the Tor process. Destructive operations prompt before changing keys; `--yes` confirms the action without a prompt. `status` and `validate` print JSON so scripts can consume the result.

Vanity prefixes use lowercase base32 characters `a-z` and `2-7`, up to six characters. By default the search uses one fewer worker than the logical CPU count, with a minimum of one. Search results are written as Tor-compatible key directories under `DATA_DIR/tor/vanity/`; a running search does not survive a server restart, while completed candidates remain available to apply.

An imported directory must contain a matching Tor v3 secret key, public key, and hostname. The server validates the key pair and address before replacing the active key. The old address stops working after regeneration or key replacement.

The CLI communicates with the running server over a loopback-only control channel. The channel is an internal implementation detail, not a public API. When no server is running, status and validation fall back to local state; commands that change Tor state require the server to be running.

See the [configuration guide](configuration.md) for Tor settings and the [CLI reference](cli.md) for other maintenance commands.
