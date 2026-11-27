# IP and domain blocklists

Blocklist data lives under `{data_dir}/security/blocklists/`. The built-in
`firehol_level1` IP source is always active. Operator-added sources use the
`server.security.blocklists.extra_sources` list in the file-based
configuration and support the `ip` and `domain` types.

```yaml
server:
  security:
    blocklists:
      enabled: true
      auto_update: true
      action: reject
      log_blocked: true
      extra_sources:
        - id: phishing_domains
          type: domain
          url: https://example.net/blocklist.txt
          enabled: true
```

Source IDs are unique lowercase `a-z0-9_` names and cannot be
`firehol_level1`. Sources require HTTPS URLs. Each source is stored separately
as `{id}.txt` with its own `{id}.last_updated` marker, so a missing or invalid
source never replaces another source's data. Files are downloaded with a
fixed timeout and size limit, parsed before installation, written to a private
temporary file, and atomically renamed into place. A source with no valid
entries is rejected; failed refreshes retain the last valid in-memory list.

IP sources accept one IPv4 or IPv6 address/CIDR per line. Invalid entries,
duplicates, private and trusted ranges, loopback, link-local, unique-local,
and carrier-grade NAT ranges are skipped. This protects local operators from
accidentally locking themselves out with a broad downloaded range.

Domain sources accept lowercase hostnames and `*.example.com` wildcard
entries. Host and pattern comparison uses IDNA ASCII names, ignores URL ports,
userinfo, and trailing dots, and never performs reverse DNS. A wildcard
matches subdomains at a DNS label boundary, but not its apex. Wildcards over
public suffixes are rejected. Request integrations should provide Origin or
Referer hosts, submitted email domains, and URLs the application is asked to
fetch or store; the server's own Host header is not a match target.

The daily `blocklist_update` scheduler task refreshes enabled sources
independently. Each source has a fixed 72-hour freshness window. Stale lists
remain available and report `degraded`; a source with no usable copy reports
`error`. Blocklist data never makes startup fail or makes core health return
HTTP 503.
