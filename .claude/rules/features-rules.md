# Feature Rules (PART 17-22)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

## CRITICAL - NEVER DO

- Never add an external scheduler for application tasks or expose internal metrics publicly.
- Never store backups without required encryption or self-update binaries without checksum verification.

## CRITICAL - ALWAYS DO

- Implement required email, scheduler, GeoIP, metrics, backup, and update behavior according to their PARTs.
- Make optional integrations opt-in as specified and report only public-safe health data.

## Key rules summary

Do not infer optional features from this summary; check the complete PART. See AI.md PART 17 through PART 22.
