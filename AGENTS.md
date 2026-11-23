# Project SPEC

Project: WORDLIST
Role: Efficient loader for AI.md

`AI.md` is the read-only implementation specification and source of truth. `IDEA.md` defines product behavior and must not weaken mandatory requirements in `AI.md`.

Before changing code, read the relevant complete PARTs in `AI.md`. Use `IDEA.md` for project-specific behavior. Never edit `AI.md`.

Go builds, tests, and linters run in the `casjaysdev/go:latest` container through the documented Makefile workflow. Never run Go directly on the host. Keep generated/runtime data outside the repository.

Keep README, `docs/`, OpenAPI, GraphQL, and server-rendered pages synchronized with the implementation. Do not log or persist generated passphrases; generate them only with `crypto/rand`.

See `AI.md` PART 0 for the complete assistant workflow and `IDEA.md` for the product specification.
