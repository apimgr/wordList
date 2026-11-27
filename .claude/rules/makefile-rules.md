# Makefile Rules (PART 25)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

## CRITICAL - NEVER DO

- Never use Makefile targets from CI/CD workflows.
- Never run Go compilation or tests directly on the host.

## CRITICAL - ALWAYS DO

- Provide the specified local development, build, test, and cleanup targets using the Go container and persistent caches.

## Key rules summary

Use temporary output for dev builds and avoid destructive broad cleanup. See AI.md PART 25.
