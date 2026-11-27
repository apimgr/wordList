# API Rules (PART 13, 14, 15)

⚠️ **These rules are NON-NEGOTIABLE. Violations are bugs.** ⚠️

## CRITICAL - NEVER DO

- Never add unversioned API routes outside the exhaustive exception list.
- Never leave a public user-facing API feature without its corresponding frontend surface.

## CRITICAL - ALWAYS DO

- Keep REST, OpenAPI, GraphQL, health output, and clients synchronized.
- Use canonical envelopes, error formats, request IDs, content negotiation, and TLS behavior.

## Key rules summary

Routes are lowercase, plural nouns, hyphenated, and have no trailing slash. See AI.md PART 13, PART 14, and PART 15.
