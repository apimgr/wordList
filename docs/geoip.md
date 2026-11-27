# GeoIP data and lookup module

The `src/geoip` package manages the four local MMDB files used by the server. Database files live under `{data_dir}/security/geoip/` and are downloaded from the canonical `sapics/ip-location-db` GitHub release assets. No database is embedded in the binary, no MaxMind account is required, and no GeoLite2 data is used.

| Package | Local file | Data |
|---|---|---|
| `iptoasn-asn` | `iptoasn-asn.mmdb` | ASN number and network organization |
| `iptoasn-country` | `iptoasn-country.mmdb` | Country ISO code |
| `dbip-city` IPv4 | `dbip-city-ipv4.mmdb` | City, country, state, postal code, coordinates, timezone |
| `dbip-city` IPv6 | `dbip-city-ipv6.mmdb` | City, country, state, postal code, coordinates, timezone |

`geoip.NewManager(directory, options)` opens valid files that already exist and remains usable when files are absent or invalid. Call `Update(ctx)` from the built-in weekly `geoip_update` task. It downloads and verifies each source independently, retries a checksum mismatch once, and replaces a live database atomically only after checksum and MMDB validation. A checksum mismatch never removes a valid live database. If no live file exists, a valid MMDB can be installed as unverified, with a warning and degraded status until a later verified update. A missing or failed source affects only that database and does not stop server startup.

The `.last_updated` marker records the successful verified update of the complete database set. Data older than 21 days stays available for lookups and reports degraded status. The country, ASN, and City databases use `github.com/oschwald/maxminddb-golang`, which accepts the custom ip-location-db `database_type` values. `geoip.LookupCountry` reports missing data or lookup errors to the caller; request middleware must fail open for GeoIP and continue all authentication and authorization checks.

Private, loopback, link-local, unspecified, and IPv4-mapped private addresses are not looked up or country-blocked. `geoip.CountryBlocked` is deny-only: a non-empty allow list takes precedence over the deny list, and an unknown country fails open. GeoIP never authenticates a caller or grants access.

The City dataset is licensed CC BY 4.0 and requires visible attribution on the About page and in the project license file:

<a href="https://db-ip.com/">IP Geolocation by DB-IP</a>

The ASN and Country datasets are PDDL from iptoasn.com and do not require attribution. Do not label AS holder names as WHOIS registrant information.
