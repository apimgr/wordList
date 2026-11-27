# Backend Rules (PART 9, 10, 11, 31)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

## CRITICAL - NEVER DO

- Never interpolate untrusted input into queries or commands.
- Never expose credentials, tokens, secrets, or private generated passphrases in logs or responses.

## CRITICAL - ALWAYS DO

- Use consistent API errors, bounded resource use, secure logging, and required middleware.
- Follow database, encryption, access-control, and overlay-network requirements where applicable.

## Key rules summary

Use parameterized queries and safe failure behavior. See AI.md PART 9, PART 10, PART 11, and PART 31.
