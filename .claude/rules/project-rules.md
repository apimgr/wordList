# Project Rules (PART 2, 3, 4)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

## CRITICAL - NEVER DO

- Never put runtime configuration, data, logs, or temporary output in repository source paths.
- Never add unsupported license text or remove required third-party attribution.

## CRITICAL - ALWAYS DO

- Keep source under `src/`, scripts under `scripts/`, tests under `tests/`, and deployment files under `docker/`.
- Use OS-specific config/data/cache/log directories and stable internal project identity.

## Key rules summary

Use the exact required project files and directory layout; inspect existing files before migration. See AI.md PART 2, PART 3, and PART 4.
