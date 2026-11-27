# TLS certificate lifecycle

The SSL package's lifecycle manager inspects the configured domains against the certificate locations in AI.md PART 15. It selects the first readable certificate and key pair whose certificate matches that domain and is currently valid:

1. `/etc/letsencrypt/live/domain/` — system managed.
2. `/etc/letsencrypt/live/{domain}/` — system managed.
3. `{config_dir}/ssl/letsencrypt/{domain}/` — application managed.
4. `{config_dir}/ssl/local/{domain}/` — user managed.

System and local certificates are never renewed by the application. Managed Let's Encrypt certificates are due for renewal at 30 days remaining or less. A managed certificate that is unreadable, mismatched, expired, or not yet valid is reported as invalid and eligible for renewal. Missing certificates are eligible for initial issuance.

Overlay hostnames ending in `.onion` or `.i2p` are reported as not using TLS certificates and are never sent to the ACME issuer; those services use their overlay transport encryption.

`ssl.NewLifecycle` accepts the config directory, configured domains, and an optional clock. `Inspect` returns safe metadata only; it never includes certificate or private-key contents. `RenewDue(ctx, issuer)` asks the supplied `ssl.Issuer` to issue only due managed certificates, validates the complete chain and domain match, then installs `fullchain.pem` and `privkey.pem` with owner-only permissions. The issuer is responsible for the challenge protocol and must not log or expose DNS credentials. Issuance failures are returned as generic status errors so provider secrets cannot leak through logs.

Managed certificate directories are created with mode `0700`; private keys and staged files use `0600`. Installation stages and validates both files before replacing the current pair, retains the old pair for rollback during replacement, and syncs the containing directory where supported. Certificate inspection does not renew system-managed `/etc/letsencrypt` certificates or user-managed local certificates.

The application uses Lego for ACME issuance. `ssl.NewACMEIssuer` accepts `ACMEConfig`, validates the selected challenge, and returns an `ssl.Issuer` compatible with `Lifecycle.RenewDue`. It supports HTTP-01, TLS-ALPN-01, and DNS-01, including every DNS provider built into Lego. The application should use `ACMEIssuer.HTTPHandler` on its existing HTTP router and wrap its existing TLS configuration with `ACMEIssuer.TLSConfig`. Neither hook binds a socket. HTTP-01 tokens are served by the app listener; TLS-ALPN-01 challenge certificates are selected on the app listener only when the client offers `acme-tls/1`.

The ACME account key and registration are saved under `{config_dir}/ssl/letsencrypt/account/{production|staging}/account.json` with mode `0600` inside a mode `0700` directory. Certificates remain in `{config_dir}/ssl/letsencrypt/{domain}/`. Set the issuer's `Staging` option for Let's Encrypt staging; production is the default. DNS credentials are passed as a flat map of Lego environment variable names to values. The issuer injects them only around `dns.NewDNSChallengeProviderByName` under a process-wide mutex, restores the previous environment immediately, and returns sanitized errors. `ValidateDNSProvider` constructs the provider without initiating a DNS update. The configuration layer must encrypt these values in `secrets.yml` before passing them to the issuer and must never render or log them.

The dependency is `github.com/go-acme/lego/v5`. Its default build includes its DNS provider registry; builds must not use Lego's `lego_dnslight` tag because selecting providers by name would otherwise omit providers. The daily `ssl_renewal` task should call `RenewDue` at 03:00 and deliver `ssl_expiring`, `ssl_renewed`, and `ssl_renewal_failed` email events according to configuration.
