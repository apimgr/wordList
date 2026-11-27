# Configuration Rules (PART 5, 6, 12)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

## CRITICAL - NEVER DO

- Never add a web settings panel or runtime configuration mutation API.
- Never store secrets in public configuration files or logs.

## CRITICAL - ALWAYS DO

- Load and validate configuration from the specified files, environment, and flags in the required precedence.
- Keep server configuration file-based and generate secure secrets into the secrets file.

## Key rules summary

Use safe defaults, atomic writes, path validation, and the documented application modes. See AI.md PART 5, PART 6, and PART 12.
