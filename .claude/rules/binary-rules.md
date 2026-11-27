# Binary Rules (PART 7, 8, 32)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

## CRITICAL - NEVER DO

- Never enable CGO or use a host Go toolchain.
- Never change the required server/client command and flag contracts.

## CRITICAL - ALWAYS DO

- Build static binaries with embedded runtime assets and required build metadata.
- Provide the required companion client and standard help, version, color, and NO_COLOR behavior.

## Key rules summary

Honor platform naming and CLI configuration/security rules. See AI.md PART 7, PART 8, and PART 32.
