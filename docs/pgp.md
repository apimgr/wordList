# PGP key lifecycle

The server maintains one project security keypair for encrypting security
reports. The public key is served at `/.well-known/pgp-key.asc` and included in
`security.txt` while a key is active. When the security contact belongs to the
request host, the same key is available through Web Key Directory at
`/.well-known/openpgpkey/hu/{hash}`; the direct-method hash is z-base-32 of the
SHA-1 digest of the lowercased contact mailbox local part.

## Maintenance commands

All commands require the operator token or root. Rotate, import, export, and
delete also ask for a reason and typed confirmation.

```text
wordList --maintenance pgp generate
wordList --maintenance pgp rotate
wordList --maintenance pgp publish
wordList --maintenance pgp export public [PATH]
wordList --maintenance pgp export private PATH
wordList --maintenance pgp import FILE
wordList --maintenance pgp delete
```

Generated keys use Ed25519 signing and Curve25519 encryption, expire two years
after creation, and use the configured security contact as their identity.
Private material is encrypted with a key derived from the installation secret
and stored under `security/` with owner-only permissions. Private-key export
creates a new mode `0600` file, refuses an existing path, and is limited to one
export per operator per hour. Imported private keys must be unlocked, include
both required algorithms, and match the configured contact unless the operator
explicitly overrides the mismatch.

Rotation stages the new keypair before updating reports and key metadata in a
single database transaction. Startup recovery checks the database fingerprint
marker and either completes a committed rotation or discards uncommitted staged
files. Previous public keys are archived by fingerprint so pending publication
retries remain possible after rotation. The previous encrypted private key is
retained for 30 days.

## Keyserver publishing

Set `web.security.keyservers` in `yaml.d/web.yml` to bare DNS hostnames. The
default is `keys.openpgp.org`; an empty list disables publishing. The server
uses VKS JSON submission for that host and HKP form submission for other
hosts. Redirects are rejected. Requests have a 15-second timeout and honor
context cancellation.

Successful submissions are recorded in `security/keyservers.state`. Network
errors, HTTP 429, and 5xx responses retry after 1 minute with exponential
backoff up to 6 hours, for at most 8 attempts. Other 4xx responses are final.
The `keyserver_retry` task checks every 15 minutes. State and PGP keys are
included in backups so a restore can preserve retry history.

## Storage

`security/pgp.pub.asc` is the active ASCII-armored public key. The private key
is stored in `security/pgp.priv.asc.enc`; archived public keys use
`security/pgp.pub.{FINGERPRINT}.asc`. `security/keyservers.state` stores
per-host publication and retry timestamps. The backup command includes the
complete `security/` configuration directory and the security report database.
