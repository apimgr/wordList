# Configuration

The server reads `server.yml` from the OS-specific configuration directory and creates a default file on first start. `server.yml` holds node-specific settings; other settings belong in the relevant `yaml.d/` area file, such as page content in `yaml.d/pages.yml`. Drop-ins are merged in byte-wise filename order after `server.yml`, and `secrets.yml` is merged last. Later values replace earlier values. The server migrates older all-in-one configurations into area files on startup. Configuration files, drop-ins, and secrets are private to the service account; `yaml.d/` is owner-only. The server creates `secrets.yml` with an `adm_` operator token on first start and prints that token once to the startup console. No web-based configuration interface is provided. Server command-line options are limited to the flags documented in the [CLI reference](cli.md); there is no generic CLI setting override.

## Current settings

| YAML key | Default | Purpose |
|---|---:|---|
| `server.port` | random `64000`–`64999` (`80` in an application container) | Listen port |
| `server.address` | `[::]` | Dual-stack listen address |
| `server.pidfile` | `true` | Write a PID file outside application containers; `--pid` selects its path |
| `server.daemonize` | `false` | Detach manual starts; service managers and application containers always keep the server in the foreground |
| `server.fqdn` | empty | First-run detected hostname recorded for operator reference; URL construction does not use it |
| `server.baseurl` | `/` | Optional path prefix for deployments mounted below a host root |
| `server.trusted_proxies.additional` | empty | Additional proxy IPs, CIDRs, or DNS names allowed to supply forwarded headers; DNS names refresh every five minutes |
| `server.trusted_proxies.client_ip_header` | `auto` | Trusted-proxy client address header; auto resolves a right-to-left proxy chain |
| `server.mode` | `production` | Application mode; `--maintenance mode` writes this setting to `yaml.d/maintenance.yml` |
| `server.api_version` | `v1` | Canonical API prefix; accepts `v1` or `v2` |
| `server.branding.title` | `WordList` | Public application title |
| `server.branding.tagline` | empty | Optional public application tagline |
| `server.branding.description` | built-in | Public application description and SEO fallback |
| `server.update.branch` | `stable` | Release channel for server update checks (`stable`, `beta`, or `devel`) |
| `server.update.auto_install` | `false` | Whether the scheduled update task stages eligible releases |
| `server.update.defer_days` | `0` | Days an update must be published before scheduled checks consider it |
| `server.update.cli_min_version` | empty | Minimum CLI version accepted by server autodiscovery |
| `server.backup.retention.max_backups` | `1` | Daily full backups to retain; valid range is 1–365 |
| `server.backup.retention.keep_weekly` | `0` | Weekly backups to retain; zero disables weekly retention |
| `server.backup.retention.keep_monthly` | `0` | Monthly backups to retain; zero disables monthly retention |
| `server.backup.retention.keep_yearly` | `0` | Yearly backups to retain; zero disables yearly retention |
| `server.backup.retention.max_total_size` | `10%` | Optional hard cap for backup storage |
| `server.backup.disk_threshold` | `90` | Maximum backup-volume usage percentage before scheduled backup is skipped |
| `server.restore.max_archive_size` | `1G` | Maximum compressed archive size accepted by restore and verification |
| `server.restore.max_expanded_size` | `4G` | Maximum combined expanded archive size |
| `server.restore.max_files` | `10000` | Maximum number of files in a restored archive |
| `server.restore.max_file_size` | `512M` | Maximum size of an individual restored file |
| `server.restore.max_ratio` | `100` | Maximum expanded-to-compressed ratio |
| `server.scheduler.timezone` | `America/New_York` | IANA timezone for scheduled task times |
| `server.scheduler.catch_up_window` | `1h` | Maximum age of a missed scheduled run that is replayed at startup |
| `server.scheduler.tasks.database_check.schedule` | `@every 1m` | Probe the database read and write handles using bounded read-only queries |
| `server.database.driver` | `sqlite` | Database driver; accepts `sqlite`, `sqlite2`, `sqlite3`, `libsql`, or `turso` |
| `server.database.url` | empty | Optional SQLite file path or remote `libsql://` / `https://` connection URL; the URL scheme selects the driver |
| `server.database.token` | empty | Separate libSQL auth token, stored in `secrets.yml` |
| `server.database.cache.enabled` | `false` | Enable the local SQLite read replica for a remote libSQL database |
| `server.database.cache.path` | `{data_dir}/db/cache.db` | Owner-only SQLite cache file path |
| `server.database.cache.sync_interval` | `30s` | Background full reconciliation interval for cached tables |
| `server.database.cache.offline_mode` | `true` | When enabled, permit eligible cached reads in read-only maintenance while remote libSQL is unavailable |
| `server.database.cache.max_age` | `1h` | Maximum cache age while remote is reachable; older tables are read remotely |
| `server.scheduler.tasks.update_check.schedule` | `0 6 * * *` | Check the configured release channel daily; notify operators once per eligible version |
| `server.notifications.database_error.enabled` | `true` | Enable alerts after consecutive failed database probes |
| `server.notifications.database_error.failures` | `3` | Consecutive failed probes before the first database alert |
| `server.notifications.database_error.repeat` | `24h` | Minimum interval between repeated alerts while probes continue to fail |
| `server.scheduler.tasks.disk_check.schedule` | `@every 5m` | Check the data, log, and backup filesystems once per volume and notify on low space |
| `server.notifications.disk_low.enabled` | `true` | Enable low disk checks, warning logs, and ready-relay email delivery together |
| `server.notifications.disk_low.threshold_percent` | `85` | Used-space percentage that triggers a low-disk notification |
| `server.notifications.disk_low.repeat` | `24h` | Minimum interval between repeated warnings for a volume that remains above threshold |
| `server.notifications.email.events.update_available` | `false` | Send the optional new-release notification to the configured server admin |
| `server.notifications.email.events.scheduler_error` | `true` | Email terminal failures from scheduled tasks without a dedicated subsystem event |
| `server.notifications.email.events.backup_failed` | `true` | Send the dedicated backup failure email and suppress a duplicate scheduler alert |
| `server.notifications.email.events.startup` / `shutdown` | `false` | Send optional lifecycle messages |
| `server.notifications.email.events.backup_complete` | `false` | Send an optional successful backup message |
| `server.notifications.email.events.ssl_expiring` / `ssl_renewed` / `ssl_renewal_failed` | `true` / `false` / `true` | Control certificate-expiry and renewal messages |
| `server.notifications.email.events.security_alert` | `true` | Send security alert messages to the server admin |
| `server.notifications.email.events.update_installed` / `update_failed` / `update_staged` | `true` | Control update lifecycle messages |
| `server.security.reports.retention_days` | `365` | Days to retain closed security reports; the daily `security_cleanup` task purges older closed reports and cannot be disabled |
| `server.security.allowlist` | `[]` | IP addresses or CIDRs that bypass configured IP blocks, rate limits, and GeoIP country policy; authentication remains required |
| `server.security.blocked_ips` | `[]` | Static IP addresses or CIDRs denied before rate limiting; an overlapping allowlist entry wins |
| `server.security.blocklists.enabled` / `auto_update` | `true` / `true` | Enforce built-in FireHOL and configured sources; download at startup and daily when automatic updates are enabled |
| `server.security.blocklists.action` | `reject` | Respond with 403 or silently drop blocked requests (`drop`) |
| `server.security.blocklists.reject_message` / `log_blocked` | translated forbidden / `true` | Optional rejection copy and blocked-request event logging |
| `server.security.blocklists.extra_sources` | `[]` | Additional HTTPS IP or domain sources; `id`, `type`, `url`, and `enabled` fields |
| `server.scheduler.tasks.blocklist_update.schedule` | `0 4 * * *` | Daily blocklist update at 04:00 in the configured scheduler timezone |
| `server.geoip.enabled` | `true` | Enable local ip-location-db ASN, Country, and City lookups |
| `server.geoip.deny_countries` / `allow_countries` | `[]` / `[]` | ISO alpha-2 country policy; a non-empty allow list takes precedence, unknown data fails open |
| `server.geoip.presets` | `{}` | Named operator-authored country lists for manual reuse; presets are never enforced automatically |
| `server.scheduler.tasks.geoip_update.schedule` | `0 3 * * 0` | Weekly GeoIP database update at 03:00 Sunday in the configured scheduler timezone |

GeoIP downloads, status behavior, and dataset attribution are documented in the [GeoIP guide](geoip.md).
Blocklist source formats, freshness, and matching policy are documented in the [blocklists guide](blocklists.md).
SQLite and remote libSQL setup, URL credential handling, cache behavior, and offline maintenance are documented in the [database guide](database.md).

### Internal metrics

`server.metrics` enables authenticated Prometheus, Grafana, and Loki endpoints. The three generated bearer tokens are stored under `server.metrics.auth.tokens` in `secrets.yml`; they are never included in configuration or health responses. Scrapes are available at `/server/metrics` and the default `/metrics` alias. Keep both paths restricted to monitoring networks.

| YAML key | Default | Purpose |
|---|---:|---|
| `server.metrics.enabled` | `true` | Enable metrics endpoints |
| `server.metrics.root.enabled` | `true` | Mount the root `/metrics` aliases |
| `server.metrics.auth.allow_unauthenticated` | `false` | Skip token checks on all metrics routes; use only on firewalled networks |
| `server.metrics.auth.tokens.prometheus` / `grafana` / `loki` | generated | Separate bearer tokens, stored in `secrets.yml` |
| `server.metrics.include_system` / `include_runtime` | `true` / `true` | Include system and Go runtime metric groups |
| `server.metrics.loki.max_entries` / `max_age` | `1000` / `1h` | Bound the recent log entries returned by the Loki endpoint |
| `server.metrics.duration_buckets` | built-in histogram buckets | HTTP duration histogram boundaries in seconds |
| `server.metrics.size_buckets` | built-in histogram buckets | HTTP body-size histogram boundaries in bytes |
| `server.pages.help.content` | empty Markdown | Help content shown on `/server/help`; accepts one shared string or a language map with `default`/first-entry fallback |
| `server.scheduler.tasks.backup_daily.enabled` | `true` | Enable the daily full and cumulative incremental backup task |
| `server.scheduler.tasks.backup_hourly.enabled` | `false` | Enable hourly incremental backups |
| `server.scheduler.tasks.{task}.retry_on_fail` | `false` | Retry a failed task up to three times with exponential delay |
| `server.scheduler.tasks.{task}.retry_delay` | `5m` | Delay before the first retry; later delays double |
| `server.healthz.root.enabled` | `false` | Mount `/healthz` as a direct alias of `/server/healthz` for integration compatibility |
| `server.rate_limit.read.requests` | `120` | Read requests per minute per IP or valid bearer token |
| `server.rate_limit.read.window` | `60` seconds | Sliding read-request window |
| `server.rate_limit.write.requests` | `10` | Write requests per minute per IP or valid bearer token |
| `server.rate_limit.write.window` | `60` seconds | Sliding write-request window |
| `server.rate_limit.health.requests` | `120` | Health/status requests per minute per IP or valid bearer token |
| `server.rate_limit.health.window` | `60` seconds | Sliding health-request window |
| `server.rate_limit.global_burst` | `240` | Absolute per-client ceiling across all endpoint classes, persisted in SQLite |
| `server.rate_limit.strict.requests` | `10` | Token issue, reveal, rotate, restore, and revoke attempts per sliding window |
| `server.rate_limit.strict.window` | `900` seconds | Strict token-action window, persisted across restarts |
| `server.rate_limit.auth_failures.requests` | `5` | Failed bearer-authentication attempts per client in a sliding window |
| `server.rate_limit.auth_failures.window` | `900` seconds | Failed-authentication lockout window |
| `server.logs.audit.enabled` | `true` | Enable append-only JSON Lines audit log |
| `server.logs.audit.filename` | `audit.log` | Audit file name under the configured log directory |
| `server.logs.audit.rotate` | `daily` | Rotation cadence; daily rotation runs at midnight UTC |
| `server.logs.audit.keep` | `none` | Retain no closed audit files by default; supported retention values include `Nd`, `Nw`, `Nm`, and `forever` |
| `server.logs.audit.compress` | `false` | Gzip retained rotated audit files |
| `server.logs.audit.events.*` | `true` | Enable configuration, security, backup, and server event categories |
| `server.logs.audit.include_user_agent` | `true` | Include sanitized request user-agent data |
| `server.i18n.default_language` | `en` | Default language for operator notifications; unsupported values fall back to English |
| `server.i18n.available_languages` | `en, es, fr, de, it, ar` | Enabled embedded locales; all six required languages must remain enabled |
| `server.tokens.max_linked` | `25` | Maximum active linked tokens per owner |
| `server.tokens.cookie_max_age` | `1y` | Sliding owner-cookie lifetime |
| `server.tokens.idle_purge_days` | `0` | Idle owner purge delay; zero disables purge. The built-in cleanup task runs every 15 minutes |
| `server.logging.level` | `info` | Configured log level |
| `server.logging.access_format` | `apache` | Configured access-log format |
| `server.contact.dpo.name` | empty | Optional DPO name shown on `/server/dpo` |
| `server.contact.dpo.email` | empty | DPO email; falls back to `server.contact.general.email` |
| `server.contact.admin.email` | empty | Private operator mailbox; never shown on public pages |
| `server.contact.abuse.email` | empty | Public abuse-report address; falls back to the general address |
| `server.contact.general.email` | empty | Public DPO fallback address; not shown if unset |
| `server.notifications.email.smtp.host` | empty | SMTP relay hostname; empty triggers first-run/reload relay autodetection unless `SMTP_HOST` is set |
| `server.notifications.email.smtp.port` | `587` | SMTP relay port; `465` uses implicit TLS in `auto` mode |
| `server.notifications.email.smtp.username` | empty | Optional SMTP authentication username |
| `server.notifications.email.smtp.tls` | `auto` | SMTP encryption mode: `auto`, `starttls`, `tls`, or `none` |
| `server.notifications.email.from.name` | app title | Display name for outbound mail |
| `server.notifications.email.from.email` | empty | Sender address for outbound mail |
| `server.notifications.email.reply_to` | empty | Reply address and default test-email recipient |
| `server.privacy.data.sold` | `false` | Whether the privacy page reports that data may be sold or shared |
| `server.privacy.data.stored_on_server` | `true` | Whether submitted data is stored on this server |
| `server.privacy.data.sharing` | analytics, email, user initiated | Conditions and data types that may be shared |
| `server.privacy.retention.period` | built-in policy | Configurable retention statement displayed on the privacy page |
| `server.privacy.cookies.*` | essential, preferences, analytics | Cookie descriptions displayed on the privacy page |
| `server.privacy.third_party.services` | empty | Third-party services and policy links displayed on the privacy page |
| `server.privacy.consent.show_until_acknowledged` | `true` | Show the server-rendered consent banner until the visitor makes a choice |
| `server.privacy.consent.message` / `message_if_sold` | built-in notice | Consent banner copy for the not-sold and sold states |
| `server.privacy.consent.policy.*` | Privacy Policy, `/server/privacy`, version `1` | Policy link shown in the consent banner |
| `server.privacy.consent.buttons.*` | Decline / I Agree | Consent banner button labels |
| `server.privacy.consent.show_preferences` / `preferences_text` | `true` / Manage Preferences | Whether the granular cookie preferences control is offered |
| `web.ui.theme` | `dark` | Initial display theme in `yaml.d/branding.yml` |
| `web.robots.allow` | `/`, `/api/v1` | Robots allow paths |
| `web.robots.deny` | empty | Robots disallow paths |
| `server.contact.security.email` | `security@wordlist.apimgr.us` | Optional secondary security contact in `security.txt` |
| `web.security.report_url` | `https://github.com/apimgr/wordList/security/advisories/new` | Primary private vulnerability reporting URL listed first in `security.txt` |
| `web.security.keyservers` | `keys.openpgp.org` | Keyserver hosts that receive generated or rotated public keys; an empty list disables publishing |
| `web.security.policy` | empty | Optional Markdown disclosure policy; empty uses the built-in policy page |
| `cli.binary_download.dir` | `{data_dir}/cli` | Operator-provisioned CLI release files served to CLI updaters |
| `cli.binary_download.require_auth` | `false` | Require a valid bearer token before serving CLI release binaries |
| `server.cors.allowed_origins` | `*` | Allowed CORS origins |

The About page can be extended, the built-in Privacy Policy can be extended, and the Terms of Service can be replaced in `yaml.d/branding.yml`. Extensions are sanitized before rendering. A language map uses the requested language, then `default`; when neither is available, the built-in content remains. Allowed formats are `md` (default), `html`, and `txt`; the About API returns the operator source when configured, while the Terms API returns the configured replacement.

Branding content expands `{current_year}`, `{project_name}`, `{project_org}`, `{project_version}`, and `{build_datetime}`. `{onion_address}` and `{i2p_address}` currently expand to empty when no published address is available.

```yaml
server:
  branding:
    content:
      privacy:
        format: md
        body:
          default: |
            ## Additional privacy details

            Contact the privacy office for questions about this deployment.
          fr: |
            ## Informations complémentaires

            Contactez le bureau de la vie privée pour toute question.
      terms:
        format: md
        body:
          default: |
            ## Terms for this deployment

            Use of this instance is subject to these additional terms.
```

The MCP help page's project section can be overridden in `yaml.d/pages.yml`:

```yaml
server:
  pages:
    help:
      mcp: |
        Use `getPassphraseByCount` to generate a passphrase with a chosen word count.

        Example prompt: Generate an eight-word passphrase separated by hyphens.
```

The Markdown is rendered with raw HTML disabled and sanitized output. Changes to `server.yml`, `secrets.yml`, recognized `.yml` or `.yaml` drop-ins, and the `assets/` directory are watched and reloaded. Invalid values fall back to their defaults; an unparseable configuration file keeps its last good copy. Listener address/port, TLS, database, cache, daemon, PID-file, and Tor settings require a restart. Rotate `server.security.encryption_key` with `wordList --maintenance secret rotate encryption_key`; do not change it by editing the key directly, because rotation must re-encrypt existing protected values. Settings not represented by the configuration schema are not available.

## Overlay-network configuration

Tor hosting starts automatically when the Tor binary is available. Its hidden-service address is stored in the application data directory and is not published by autodiscovery unless `expose` is enabled. Outbound HTTP requests only use Tor when explicitly requested and `use_network` is enabled; a request for Tor routing fails if the Tor SOCKS service is unavailable.

| YAML key | Default | Purpose |
|---|---:|---|
| `server.tor.binary` | auto-detect | Optional Tor executable path |
| `server.tor.expose` | `false` | Include the onion address in autodiscovery and health output |
| `server.tor.contact_email` | empty | Contact address shown on onion responses; never falls back to clearnet contact addresses |
| `server.tor.use_network` | `false` | Enable Tor-routed outbound HTTP clients |
| `server.tor.max_circuits` | `32` | Maximum pending client circuits; valid range 1–128 |
| `server.tor.startup_timeout` | `60` seconds | Wait for Tor control-port startup; valid range 10–300 |
| `server.tor.circuit_timeout` | `60` seconds | Tor circuit build timeout; valid range 10–300 |
| `server.tor.bootstrap_timeout` | `180` seconds | Maximum initial network bootstrap wait; valid range 30–600 |
| `server.tor.safe_logging` | `true` | Enable Tor's SafeLogging setting |
| `server.tor.max_streams_per_circuit` | `100` | Maximum hidden-service streams per circuit; valid range 10–500 |
| `server.tor.close_circuit_on_stream_limit` | `true` | Close a circuit when its stream limit is reached |
| `server.tor.bandwidth_rate` / `bandwidth_burst` | `1 MB` / `2 MB` | Tor bandwidth rate and burst limits |
| `server.tor.max_monthly_bandwidth` | `100 GB` | Tor accounting limit, or `unlimited` |
| `server.tor.num_intro_points` | `3` | Hidden-service introduction points; valid range 3–10 |
| `server.tor.virtual_port` | `80` | Onion service port exposed to visitors |
| `server.i2p.enabled` | container-dependent | Enable an eepsite; defaults on in containers and off elsewhere |
| `server.i2p.expose` | `false` | Include the I2P address in autodiscovery and health output |
| `server.i2p.binary` | auto-detect | Optional dedicated `i2pd` executable path |
| `server.i2p.sam_address` | `127.0.0.1:7656` | External SAM bridge address when a managed `i2pd` binary is unavailable |
| `server.i2p.virtual_port` | `80` | Eepsite port shown to visitors |
| `server.i2p.inbound_length` / `outbound_length` | `3` / `3` | Tunnel hop counts; valid range 0–7 |
| `server.i2p.inbound_quantity` / `outbound_quantity` | `5` / `5` | Parallel tunnel counts; valid range 1–16 |
| `server.i2p.signature_type` | `7` | SAM destination signature type |
| `server.i2p.bootstrap_timeout` | `300` seconds | Maximum I2P destination/tunnel startup wait; valid range 30–600 |

Set `I2P_ENABLED=true|false` to override the configured I2P activation. The `--i2p` / `--i2p=false` flags provide a command-line override. Tor settings take effect after restart; I2P settings reload dynamically.

When I2P is enabled, the built-in `i2p_health` task checks the provider every 10 minutes and queues a restart if an active provider has stopped responding. It is not registered when I2P is disabled.

## Environment variables

These are the supported environment variables. Server settings do not have a generic `{PROJECT_NAME}_{SECTION}_{KEY}` environment override. For settings with both a flag and environment variable, the command-line flag takes precedence.

| Variable | Scope and behavior |
|---|---|
| `NO_COLOR` | Runtime; a non-empty value disables ANSI color and emoji output unless an explicit setting overrides it. |
| `TERM` | Runtime; `TERM=dumb` disables ANSI escapes and forces CLI mode. |
| `DOMAIN` | Runtime; comma-separated public hostnames used for URL and host resolution. |
| `PORT`, `LISTEN`, `MODE`, `DEBUG`, `PID_FILE`, `I2P_ENABLED` | Runtime; listen port/address, application mode, debug flag, PID path, and I2P activation. `--port`, `--address`, `--mode`, `--debug`, `--pid`, and `--i2p` take precedence where applicable. |
| `SMTP_HOST`, `SMTP_PORT`, `SMTP_USERNAME`, `SMTP_PASSWORD`, `SMTP_FROM_NAME`, `SMTP_FROM_EMAIL`, `SMTP_TLS` | Runtime SMTP settings; credentials belong in `secrets.yml`. If `SMTP_HOST` is set, relay autodetection is skipped. |
| `CONFIG_DIR`, `DATA_DIR`, `LOG_DIR`, `BACKUP_DIR`, `CACHE_DIR`, `DATABASE_DIR` | Directory overrides applied at every startup. The matching `--config`, `--data`, `--log`, `--backup`, or `--cache` flag takes precedence. `DATABASE_DIR` has no CLI flag. |
| `APPLICATION_NAME`, `APPLICATION_TAGLINE` | First-run initialization only; seed the branding settings, then are ignored. |
| `DATABASE_DRIVER`, `DATABASE_URL`, `CACHE_URL` | First-run initialization only; seed the database/cache connection. Later environment changes do not replace the saved configuration. Credential-bearing URLs are stored in `secrets.yml`. |

Directory environment variables are read on every start, not from `server.yml`. Database and cache connection variables are initialization-only. Do not assume a setting is environment-overridable unless it appears in this table.

`SMTP_PASSWORD` is a supported runtime variable, but storing relay credentials in `secrets.yml` is preferred. When no SMTP host is configured and `SMTP_HOST` is unset, the server probes candidate relays; failure to find one leaves email disabled without preventing startup.

The built-in scheduler creates a daily full backup at 02:00 in `server.scheduler.timezone` and replaces the cumulative daily incremental. Set `server.scheduler.tasks.backup_hourly.enabled: true` to enable hourly incrementals. Both tasks use the configured backup retention, disk threshold, encryption password, and restore verification limits. Set a task's `retry_on_fail: true` and optional `retry_delay` to retry failures up to three times with exponential backoff. Restore validates the full archive before changing files, uses rollback files to recover previous destinations if installation fails, and preserves the target installation's operator token.

SMTP credentials belong in `secrets.yml` under `server.notifications.email.smtp.password`. `SMTP_HOST`, `SMTP_PORT`, `SMTP_USERNAME`, `SMTP_PASSWORD`, `SMTP_TLS`, `SMTP_FROM_NAME`, and `SMTP_FROM_EMAIL` override their corresponding settings for container deployments. A configured relay is checked with EHLO at startup and after supported config reloads. If no host is configured and `SMTP_HOST` is unset, startup probes loopback, Docker bridge, default gateway, FQDN, global IPv4, and the `mail.`/`smtp.` FQDNs on ports 25, 465, and 587. The first relay that completes the TLS-aware SMTP handshake is saved to `server.yml`; if none responds, email stays disabled without interrupting startup. Templates default to embedded messages; valid overrides belong in `{config_dir}/template/email/{name}.txt` and are read when used so edits apply immediately. Use `wordList --maintenance email test --dry-run` to validate and preview every template without sending, or `wordList --maintenance email test [--to ADDRESS]` to send a test message through a ready relay.

Absolute URLs are built for each request. A valid `DOMAIN` environment variable (a comma-separated host list) constrains untrusted request hosts; otherwise a valid request `Host` can be used. Forwarded host, scheme, and port headers are honored only when the TCP peer is in the built-in trusted private/loopback ranges or `server.trusted_proxies.additional`. For example:

```yaml
server:
  trusted_proxies:
    additional:
      - 203.0.113.10
      - 198.51.100.0/24
      - proxy.example.com
```

Set `DOMAIN=example.com,www.example.com` when the public host should be fixed. The first valid entry is the fallback hostname. URL construction strips the default HTTP and HTTPS ports.

For a path-mounted deployment, set `server.baseurl: /wordlist` or send a validated prefix from a trusted proxy. Rendered page links, forms, absolute URLs, and redirects include the prefix; requests arriving with that prefix are routed after it is removed.

The DPO page always links to the privacy policy. On clearnet it displays the configured DPO name and email when present; otherwise it uses the general contact email, then links to the security page if no address is configured. On Tor it hides the configured DPO name and clearnet email and shows only `server.tor.contact_email`, if set. The server does not invent or expose an administrator address.

Runtime assets are served from the configuration `assets/` directory and are reloaded when that directory changes. Encryption-key rotation is available through `wordList --maintenance secret rotate encryption_key`; listener address and port, TLS, database, cache, daemon, PID-file, and Tor changes require a restart.

The generated operator token is used by protected token-management operations and MCP operator tools. Generated metrics scraper tokens and webhook signing secrets are stored in `secrets.yml` and are never printed. Sensitive maintenance operations have explicit authorization rules: setup is first-run-only or root; restore requires an empty database, root, or the operator token; mode changes require root or the service user's operator token; PGP private-key import/export/delete and installation/encryption-secret rotation require root or the operator token plus typed confirmation; linked-token revocation requires root or the operator token plus confirmation. See `wordList --maintenance --help` for available maintenance commands.

The audit log is the authoritative JSON Lines record and uses a SHA-256 chain with `audit.anchor`; SQLite `audit_log` rows are a best-effort query mirror. Token lifecycle, server startup/shutdown, generated-secret, integrity, and strict rate-limit events are recorded when their category is enabled. Daily rotation preserves the chain anchor, applies the configured retention, and optionally compresses retained archives.
