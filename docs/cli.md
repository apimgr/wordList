# CLI reference

## Server binary

The server starts in the foreground by default. Its current top-level options include:

```sh
wordList [--config DIR] [--data DIR] [--cache DIR] [--log DIR] [--backup DIR]
         [--pid FILE] [--address ADDR] [--port PORT] [--baseurl PATH]
         [--mode MODE] [--daemon] [--debug] [--color auto|yes|no]
wordList --status
wordList --service --install
wordList --service start|stop|restart|reload|disable|status
```

The service installer supports the detected system service manager and currently requires administrator privileges. `--service --uninstall` asks for confirmation before removing the service registration. Application containers stay in the foreground and do not create PID files.

Maintenance commands include backups, restore, update operations, and persistent server-mode changes:

```sh
wordList --maintenance backup [FILENAME] [--include-ssl] [--include-data]
wordList --maintenance backup test [FILE]
wordList --maintenance email test [--to ADDRESS] [--dry-run]
wordList --maintenance restore FILE
wordList --maintenance mode {production|development|debug}
wordList --maintenance security list [--status STATUS]
wordList --maintenance security show TRACKING_ID
wordList --maintenance security export TRACKING_ID [--output FILE]
wordList --maintenance security update TRACKING_ID [--status STATUS] [--comment TEXT] [--disclosure-date DATE] [--severity LEVEL]
wordList --update [check|yes|branch {stable,beta,devel}]
```

Restore requires root, a first-run empty database, or the service user's operator token, and always requires typed confirmation. Encrypted restore and verification read the password without echoing it; passwords are never accepted as command-line arguments.

Changing the server mode requires root or the service user's operator token. The command updates only `server.mode` in `yaml.d/maintenance.yml`, preserves other drop-in values, and requires a service restart to take effect. Root changes print a warning because development diagnostics can expose additional information.

Security report maintenance requires root or the `wordList` service user with the operator token. `list` and `show` print metadata only. `export` writes the encrypted report payload to a new mode-`0600` file; use `--output -` to write it to standard output.

The email test command requires root or operator-token authorization. Use `--dry-run` to validate and preview all embedded or custom templates without sending; invalid overrides are reported and make the command exit with status 65. A real test requires a ready SMTP relay and uses `server.notifications.email.reply_to` as its default recipient; pass `--to` to override it. Test sends are recorded in the audit log.

`security update` changes report status, researcher-safe comment, expected disclosure date, or severity. Statuses advance through `received`, `triaged`, `confirmed`, `patching`, and `disclosed`, or can move to `wont_fix`; a closed status cannot be reopened.

`wordList-cli` can generate a passphrase locally or use a configured WordList server for supported API commands.

```sh
wordList-cli passphrase [--words 6] [--separator -] [--capitalize=false]
wordList-cli --config staging --server https://example.test random --count 6
wordList-cli --server https://example.test words
wordList-cli --server https://example.test count
wordList-cli --server https://example.test letter a
wordList-cli --server https://example.test length 5
wordList-cli --server https://example.test list
wordList-cli --server https://example.test search QUERY
wordList-cli --server https://example.test --limit 20 green words
wordList-cli --server https://example.test dice CODE
wordList-cli --update check
wordList-cli --update yes
```

For the WordList search API, bare arguments are joined into a query. Search
options may appear before or after the terms: `--limit` accepts values from 1
to 1000, and `--mode` accepts `prefix` or `contains`.

| Option | Default | Description |
|---|---|---|
| `--words` | `6` | Number of words, from 1 to 20 |
| `--separator` | `-` | Separator, up to eight characters |
| `--capitalize` | `true` | Capitalize words; set `false` to preserve lowercase |
| `--server URL` | CLI config | Server URL for remote commands |
| `--token TOKEN` | CLI config | Optional Bearer token |
| `--token-file FILE` | CLI config | Read a token from a user controlled file |
| `--lang CODE` | CLI config | Select a supported client language; unsupported codes fall back to English |
| `--output table|json|plain` | `table` | Format remote responses as readable tables/lists, indented JSON, or server plain text |
| `--color auto|yes|no` | `auto` | Color table headings; `NO_COLOR` disables color |
| `--debug[=BOOL]` | `false` | Append request method, redacted URL, attempt, and status to the private CLI log |
| `--quiet[=BOOL]` | `false` | Suppress routine status output |
| `--verbose[=BOOL]` | `false` | Enable verbose status output |
| `--update check|yes` | — | Check for a published client update or download, verify, replace, and restart with it |
| `--config FILE` | `cli.yml` | Select a YAML config by name, path, or path with `.yml`/`.yaml` extension |
| `--help`, `-h` | — | Show usage |
| `--version`, `-v` | — | Show version |
| `--shell completions [SHELL]` | shell from `$SHELL` | Print built-in completion definitions |
| `--shell init [SHELL]` | shell from `$SHELL` | Print a command to load completions in the current shell |

Both binaries print a parseable first version line, followed by build time, Go version, and target OS/architecture. The version command does not add color.

With no command, the client starts its first-run setup prompts when no server is configured, then opens the interactive Bubble Tea TUI in a terminal or the native GUI on a local display when built with the `gui` tag. In a non-interactive environment it exits with usage code `2` and names `server.primary` as the configuration path. `display.mode` accepts `auto`, `tui`, or `gui`; `tui.enabled` enables or disables the TUI. The default static CLI has no GUI dependencies. `make local` and `make build` also create separate GUI-tagged desktop artifacts for Linux, macOS, and Windows; these artifacts require the platform's dynamic loader and graphics runtime. GUI command screens remain under development; use the CLI or TUI for complete command access.

The CLI creates its user-scoped config, data, cache, and log directories on startup with mode `0700` on Unix. It creates `cli.log` with mode `0600`; credentials stay in `cli.yml`, also mode `0600`. Unix paths honor `XDG_CONFIG_HOME`, `XDG_DATA_HOME`, `XDG_CACHE_HOME`, and `XDG_STATE_HOME`; Windows uses `%APPDATA%` and `%LOCALAPPDATA%`. `--config staging` selects `staging.yml` (or an existing `staging.yaml`) under the config directory; absolute and `~/` paths are also supported. Passing a valid `--server` saves it only when the selected config has no valid primary URL. Authentication priority is `--token`, `WORDLIST_TOKEN`, `auth.token`, then `--token-file` or `auth.token_file`. Tokens supplied with `--token` are saved only when the configured token is empty or invalid; valid owner tokens use the `own_` prefix and 32 alphanumeric characters. Missing configuration is initialized with mode `0600` defaults. Output settings live under `output.format`, `output.color`, `output.quiet`, and `output.verbose`; environment overrides are `WORDLIST_OUTPUT_FORMAT`, `WORDLIST_OUTPUT_COLOR`, `WORDLIST_QUIET`, and `WORDLIST_VERBOSE`. `WORDLIST_DEBUG` overrides `debug`. Remote `table` and `json` output negotiate JSON, while `plain` requests server text. HTTPS is used when the configured server uses HTTPS.

The `server` config supports `api_version` (default `v1`), `timeout` (default `30s`), `retry` (default `3` retries), and `retry_delay` (default `1s`). `update.auto` defaults to false, with `check_interval: per_invocation` and `channel: stable`. The client reads `/api/autodiscover` on startup when connected to a server and on `--update check`; it verifies a downloaded binary against the published SHA-256 before staging an atomic replacement. The server must publish the current platform in `cli_versions` and serve `/cli/binaries/wordList-cli-{os}-{arch}` for updates to be available. The client retries transport failures, HTTP 408/429, and server errors; it uses the configured API version for versioned requests.

Errors are written to standard error and use the documented exit codes: `0` success, `1` general error, `2` usage error, `65` invalid input, `66` resource not found, `69` connection error, `77` authentication or permission error, and `78` configuration error. In a non-interactive environment, starting without a command exits with usage code `2`; use `--help` to print usage.

Both binaries provide programmable completions for Bash, Zsh, Fish, and PowerShell. `sh`, `dash`, and `ksh` support only the harmless `--shell init` no-op; requesting completions for those shells exits with usage code `2`. For example, load them for the current Bash session with `eval "$(wordList-cli --shell init bash)"`.
