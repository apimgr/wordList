# Database

The server defaults to a local SQLite database at `{db_dir}/server.db`. Set `DATABASE_DIR` to change the SQLite directory. Configuration remains in `server.yml`, drop-ins, and `secrets.yml`; it is never stored in the database.

## Driver and URL

`server.database.driver` accepts `sqlite`, `sqlite2`, `sqlite3`, `libsql`, and `turso`. SQLite aliases resolve to the pure Go `modernc.org/sqlite` driver. `libsql` and `turso` use the pure Go, remote-only `libsql-client-go` driver. A configured URL takes precedence over the default SQLite path, and its scheme selects the driver. A URL/driver mismatch logs a warning and follows the URL.

```yaml
server:
  database:
    driver: libsql
    url: libsql://database.example.turso.io?authToken=${TURSO_AUTH_TOKEN}
    # Or use https://database.example.turso.io with a separate token.
    token: ${TURSO_AUTH_TOKEN}
```

SQLite URLs are file paths or `file:` URLs. Remote libSQL URLs use `libsql://` or `https://`. The `DATABASE_URL` environment variable seeds a first-run connection only; `DATABASE_DRIVER` seeds a first-run driver only when `DATABASE_URL` is unset. Later environment changes do not replace the saved connection. A URL query option selects the auth token over `server.database.token`.

Credential-bearing URLs and `server.database.token` are kept in `secrets.yml`. The server removes migrated credentials from `server.yml` and configuration drop-ins. Connection values are never included in database startup errors or maintenance health details.

## Optional local cache

For a remote primary, set `server.database.cache.enabled: true` to enable a separate SQLite read replica. The default is disabled. The cache path defaults to `{data_dir}/db/cache.db`; its file is mode `0600` and its directory is mode `0700`. It is opened with SQLite, never through the remote driver.

The cache contains owner, token, scheduler, backup, security report, security thanks, and CAPTCHA tables. It excludes generated secrets, audit data, and per-request counters. A completed synchronization records a UTC `synced_at` and `watermark` for every cached table in `cache_meta`; the local snapshot and metadata commit together. Tables with an `updated_at` column pull rows at or after the previous watermark, upsert those rows, and reconcile deleted keys. Tables without that column are reread in full.

Reads use a fresh cached table first. When its `max_age` expires, reads use the remote database while it is reachable. With `offline_mode: true` (the default), an existing snapshot can serve the eligible public read routes while remote access is down; those requests run against a read-only Store, and writes and features that require remote-only secrets remain unavailable. A cache without a complete successful sync is never served. With `offline_mode: false`, remote failure enters maintenance without cached read fallback.

Writes are committed remotely first. A failed remote write leaves the cache untouched. After a successful remote write, the cache is reconciled; if that update fails, rows that differ from the remote source are deleted from the cache. If that invalidation also fails, the cache is marked dirty and bypassed until a complete synchronization commits.

## Backup and restore

SQLite backups include `server.db`. A libSQL backup exports `app_secrets.json` from the configured remote database but never includes the remote database or the derived SQLite cache. Restore into libSQL applies those secrets only when the remote `app_secrets` table is empty. If the archive contains a SQLite database and the remote data tables are empty, the server seeds them from that file in one transaction. A remote database with existing application data is left intact. The restored SQLite source is staged in a private project temporary directory and removed after restore.
