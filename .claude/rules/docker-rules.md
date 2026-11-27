# Docker Rules (PART 26)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

## CRITICAL - NEVER DO

- Never copy or symlink host-built binaries into runtime images.
- Never run containers as root unless the specific requirement permits it.

## CRITICAL - ALWAYS DO

- Build a minimal static runtime image from the prescribed multi-stage workflow and provide production, development, and test Compose configurations.

## Key rules summary

Keep Docker assets in `docker/` and runtime data in mounted OS paths. See AI.md PART 26.
