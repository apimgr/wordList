# TODO.AI.md

## Admin web UI removal (2026-09-17)

- [x] **Removed forbidden `src/admin/` package.** `AI.md` PART 16 is explicit:
      "there is no admin web UI, no dashboard, and no settings pages (no auth,
      no user accounts)". Deleted `src/admin/handlers.go` and
      `src/admin/auth.go` (login/session/dashboard/settings routes), removed
      the `admin.NewHandler`/`RegisterRoutes` wiring and import from
      `src/server/server.go`, and removed the now-unused `AdminConfig`/
      `SessionConfig` types and fields from `src/config/config.go` (kept
      `WebSecurity.Admin`, which is an unrelated security-contact email).
      Also dropped the stale `/admin` entry from the default `WebRobots.Deny`
      list. Build/vet/test verified clean in Docker after removal.

- [x] **Renovate config and `daily.yml` were missing from the standard
      workflow set.** Added `renovate.json` and
      `.github/workflows/daily.yml` (matching sibling apimgr repos), so all
      5 standard workflows (`ci.yml`, `docker.yml`, `daily.yml`,
      `release.yml`, `beta.yml`) plus `renovate.json` now exist.

## CI/CD gaps found while adding standard workflow set (2026-09-17)

- [ ] **Zero test coverage vs. mandatory 60% gate.** `AI.md` PART 27 requires
      `ci.yml`'s `test` job to enforce a 60% coverage threshold, and `ci.yml`
      already implements that gate correctly (`THRESHOLD=60` in
      `.github/workflows/ci.yml`). However the project currently has no Go
      test files, so `go test -cover ./...` will report 0% and the `test` job
      will fail on the next push/PR until real tests are added. Add unit
      tests across `./src` until `go tool cover -func` reports >= 60% total.

- [ ] **Missing `docker/Dockerfile.dev`.** `AI.md` PART 27's Docker Workflow
      section specifies a standard+development image split
      (`docker/Dockerfile` for the standard image, `docker/Dockerfile.dev`
      for the `:devel` image with debug tooling). Only `docker/Dockerfile`
      exists in this repo. `.github/workflows/docker.yml`'s `build-devel` job
      has been guarded with `hashFiles('docker/Dockerfile.dev') != ''` so it
      currently no-ops instead of failing, but the `:devel` image itself is
      not being built/published. Create `docker/Dockerfile.dev` (alpine base,
      app binary + debug tooling) to close this gap and start producing the
      `:devel` tag.

- [x] **`main.go` didn't declare the LDFLAGS-injected build-info vars.**
      Fixed 2026-09-17: `CommitID`, `BuildEpoch`, and `OfficialSite`
      package-level vars added to `src/main.go`'s existing var block
      alongside `Version`/`BuildTime`, so the `release.yml`/`beta.yml`/
      `daily.yml` LDFLAGS `-X` targets now resolve. Still not surfaced in
      `--version` output — that remains open, folded into the CLI-flags
      item below.

- [ ] **go-lint flagged pre-existing Makefile/main.go convention
      violations** (not introduced by the CI/CD workflow additions):
      `Makefile` hardcodes PROJECTNAME/PROJECTORG instead of inferring
      from git remote, its LDFLAGS references `main.Commit`/`main.BuildDate`
      which don't match the actual var names (`CommitID`/`BuildTime`/
      `BuildEpoch`), uses `golang:alpine` instead of `casjaysdev/go:latest`,
      is missing `-e GOFLAGS=-buildvcs=false` on its Docker invocation and
      `-trimpath` on its own `go build` lines, and uses `macos`/`bsd`
      instead of `darwin`/`freebsd` in binary output names; `src/main.go`
      is missing `-h`/`-v` short flag forms, `--debug`, and `--color`
      (auto/yes/no, default auto), and gates emoji output (🛑 ✅ 📡 📊 📚
      📝 ❌ ⚠️) on nothing (should respect `NO_COLOR`) and doesn't surface
      `CommitID`/`BuildEpoch`/`OfficialSite` in `--version` output;
      `src/config/config.go` loads config via `os.ReadFile` at runtime
      instead of `go:embed`; `src/words/` should be singular (`src/word/`)
      per Go directory-naming convention. Full detail from the go-lint
      agent run on 2026-09-17.

- [ ] **PART 12 "Trusted Proxies" client-IP/FQDN gate not implemented.**
      `chi/v5` v5.3.0 deprecated `middleware.RealIP` (IP-spoofing risk,
      GHSA-3fxj-6jh8-hvhx) as part of the govulncheck fix on 2026-09-17;
      it was removed from `src/server/server.go` rather than replaced,
      since implementing the full `trusted_proxies`-gated resolution
      chain (`X-Forwarded-*` trust gate, original-peer preservation,
      `BuildURL(r, ...)`) described in AI.md PART 12 → "Trusted Proxies"
      and PART 8 → "Resolution Order" is a real feature, out of scope for
      that CI fix. Until built, the server falls back to `r.RemoteAddr`
      for logging (correct default for a no-trusted-proxy deployment per
      AI.md line 16044, but `{fqdn}`/`{proto}`/`{port}` reverse-proxy
      header detection and client-IP-based rate limiting/blocklists/GeoIP
      described elsewhere in AI.md are not implemented at all).
