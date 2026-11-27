# Interactive API documentation assets

The Swagger UI assets are the pinned `swagger-ui-dist` 5.36.0 release. The GraphiQL assets use `graphiql` 6.0.0, `@graphiql/toolkit` 1.0.0, React 19.3.1, React DOM 19.3.1, Monaco Editor 0.57.0, and Monaco GraphQL 2.0.1. `graphiql-vendor.mjs` and the three Monaco workers are minified local bundles built from those pinned package versions; third-party notices are in `src/server/static/vendor/graphiql/THIRD_PARTY_NOTICES.txt`. The server generates Monaco editor colors from the same Go theme palettes as the page CSS.

The UI libraries and styles are served locally. Page startup code remains in `src/server/static/js/docs-swagger.js` and `src/server/static/js/docs-graphql.js`.
