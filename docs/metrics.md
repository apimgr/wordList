# Metrics package

`src/metrics` provides an isolated Prometheus registry for application
instrumentation. Its default namespace is `wordlist`, derived from the project
name by lowercasing it and replacing each run of non-alphanumeric characters
with one underscore. The package exposes
`Registry.Handler()` as Prometheus text exposition; the server owns endpoint
registration, bearer-token protection, configuration, and collector scheduling.

## Use

The server creates one registry at startup and passes build metadata through
`metrics.Options`. Its HTTP middleware records every request using bounded
route labels. A 30-second collector samples Go runtime and system data when
those groups are enabled. SQLite connections, authentication, tokens,
scheduler tasks, rate limits, passphrase generation, word lookups, Tor, and
bounded caches are wired to the registry.

The API is intentionally instance-based so the server can control registry
lifecycle and avoid global collectors:

- `New(Options) *Registry`
- `Handler() http.Handler`
- `HTTPMiddleware(http.Handler) http.Handler`
- `SetApplicationInfo(version, commit, buildDate, goVersion string)`
- `UpdateUptime()` and `UpdateRuntime()`
- `SetSystemStats(cpuPercent, memoryPercent, memoryUsed, memoryTotal, diskPercent, diskUsed, diskTotal float64)`
- `ObserveDBQuery(operation, table string, duration time.Duration, err error)`
- `SetDBConnections(open, inUse int)`, `RecordDBBusy()`, and `ObserveDBWriteWait(time.Duration)`
- `RecordCacheHit(cache string)`, `RecordCacheMiss(cache string)`, `RecordCacheEviction(cache string)`, and `SetCacheStats(cache string, items, bytes int64)`
- `RecordAuth(method, status string)` and `SetOwnerTokens(primary, linked int)`
- `RecordSchedulerStart(task string)` and `RecordSchedulerEnd(task string, duration time.Duration, err error)`
- `RecordRateLimit(limit, status string)`, `RecordPassphraseGenerated()`, and `RecordWordLookup()`
- `SetTorStatus(enabled, running, circuitEstablished bool)` and `RecordTorRequest()`

## Exported metric groups

All metrics use the `wordlist_` prefix by default. HTTP request count and
duration, request and response body sizes, active requests, database query and
error counts, query duration, pool connections, SQLite busy and write-wait
metrics, auth attempts, owner-token counts, scheduler runs/duration/running/last
run, rate-limit counts, application info/start/uptime, Go runtime stats, system
CPU/memory/data-directory disk stats, passphrase generation, word lookups, and
Tor availability/request metrics are available from the registry. Cache metrics
cover hits, misses, evictions, current item count, and bytes. Histogram buckets
are defined in the package and are stable instrumentation details.

HTTP labels use allow-listed methods, response status codes, and a fixed route
pattern list. Unknown routes collapse to `other`; query strings are ignored.
Database operation, table, cache name, auth method/status, scheduler task, and
rate-limit labels also use allow-lists. Cache names are limited to `blocklist`,
`geoip`, `rate_limit`, `tokens`, and `trusted_proxy_dns`; unknown names collapse
to `other`. Disk labels use the fixed value `data_dir`. The
package never labels metrics with client addresses, request or token IDs,
credentials, body contents, arbitrary paths, SQL text, or raw error messages.
The registry retains counters and gauges in memory only; it does not persist
per-user measurements or metric samples.

The server serves Prometheus text at `/server/metrics/prometheus` and provides
the Grafana dashboard and Loki log views at the corresponding service paths.
The same handlers are available through configured root and API aliases.
Metrics are internal: each service has its own bearer token, and the default
configuration does not allow unauthenticated access. The optional
`server.metrics.include_runtime` and `server.metrics.include_system` settings
control their respective metric groups.
