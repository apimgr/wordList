# CI/CD Rules (PART 27)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

## CRITICAL - NEVER DO

- Never use Makefile targets in CI or use unpinned third-party actions.
- Never weaken test, coverage, secret-scan, dependency, image, or workflow-policy gates on one provider.

## CRITICAL - ALWAYS DO

- Provide the required provider workflows, secure permissions, and release validation.

## Key rules summary

Use provider-specific conventions and keep provider workflows in their correct directories. See AI.md PART 27.
