# Metrics

WordList exposes Prometheus text, an importable Grafana dashboard, and recent Loki-compatible log streams. These are internal monitoring endpoints. Keep them behind a firewall or private reverse proxy even when bearer authentication is enabled.

## Endpoints

| Service | Paths | Response |
|---|---|---|
| Prometheus | `/server/metrics`, `/server/metrics/prometheus`, `/api/v1/server/metrics`, `/api/v1/server/metrics/prometheus`, `/api/metrics`, `/api/metrics/prometheus` | Prometheus text exposition |
| Grafana | The same path prefixes with `/grafana` | Importable dashboard JSON |
| Loki | The same path prefixes with `/loki` | Recent structured log entries in Loki stream JSON format |

The `/metrics`, `/metrics/prometheus`, `/metrics/grafana`, and `/metrics/loki` root aliases are enabled by default and can be disabled with `server.metrics.root.enabled`. Aliases invoke the same handler directly and do not redirect. All metrics routes return `404` when `server.metrics.enabled` is false.

Each service uses its own bearer token from `server.metrics.auth.tokens.{prometheus,grafana,loki}` in `secrets.yml`. Empty tokens are generated and saved at startup. Supply them in the HTTP `Authorization: Bearer …` header; query-string tokens are not accepted. Set `server.metrics.auth.allow_unauthenticated: true` only on a network protected by firewall rules. This setting bypasses token checks for all three services.

Example Prometheus scrape configuration:

```yaml
scrape_configs:
  - job_name: wordList
    metrics_path: /metrics
    static_configs:
      - targets: ["app.internal:8080"]
    authorization:
      credentials_file: /run/secrets/wordlist-prometheus-token
```

## Metric coverage

The registry contains application version and uptime, HTTP request counts/durations/sizes and active requests, database queries and pool state, cache activity and size, authentication and owner-token counts, scheduler activity, rate limits, passphrase generation, word lookup, Tor status/requests, and selected Go runtime and system measurements. System and runtime groups can be disabled with `server.metrics.include_system` and `server.metrics.include_runtime`.

Metrics use the `wordlist_` prefix. HTTP routes and other label values are bounded by allow-lists; unknown values collapse to `other`. Client addresses, request or token IDs, credentials, body contents, arbitrary paths, SQL text, and raw errors are not metric labels. Samples are held in memory and are not persisted per user.
