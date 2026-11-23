# Remaining AI.md alignment work

This list records verified implementation gaps that remain after synchronizing the project documentation and OpenAPI inventory. Completion entries must remain until the change that resolves them is committed, per AI.md PART 0.

## Priority order

1. Implement the mandatory server/API foundations and reconcile API route semantics, response envelopes, content negotiation, health output, security middleware, and project route inventory.
2. Implement the complete standard server page family and shared page shell, including preferences/help hubs, MCP, contact, DPO, and security pages with their required APIs and workflows.
3. Bring all rendered pages into PART 16 conformance: no-JS-first forms, accessible shared chrome, responsive theme tokens, PWA behavior, SEO, and removal of inline styles/scripts/handlers.
4. Implement shared i18n for server and CLI, the required supported locale set and parity checks, and translate every user-facing string.
5. Complete the required client behavior and align its command/config/update lifecycle with AI.md PART 32.
6. Complete required root/project infrastructure and supporting docs/tests: provider workflows, Docker variants, service templates, integration and E2E scripts, page conformance, route coverage, and API/docs parity.
7. Verify the complete implementation using the project-required container commands and report actual results.

## Completed in current working tree (not committed)

- Replaced README with a route- and behavior-focused guide.
- Added the required MkDocs/ReadTheDocs configuration and documentation pages for installation, configuration, REST/GraphQL, CLI, security, integrations, development, and translation status.
- Added the required root Claude loader, its AGENTS.md copy, Claude memory index, and PART-grouped Claude rule summaries.
- Removed the unsupported unversioned `/graphql` route and changed the GraphQL API routes to POST-only.
- Expanded the embedded OpenAPI document to cover the implemented JSON routes.
- Added the required IDEA.md project route inventory and removed its unsupported `/graphql` product claim.
