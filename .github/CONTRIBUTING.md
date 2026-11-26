# Contributing

## Local setup and checks

Go builds, tests, and linters run in the `casjaysdev/go:latest` container through the project Makefile. Docker is required; do not install or run a host Go toolchain.

```sh
make dev
make test
make lint
make build
./tests/run_tests.sh
```

Documentation is built with MkDocs using `docs/requirements.txt` and `mkdocs.yml`; use the project's documentation container workflow in `docs/development.md`.

## Branch and pull request workflow

Keep `main` as the only persistent remote branch. Make changes on `main`, or use a temporary local branch that you merge into `main` and delete before pushing; do not push working branches. Check the current branch before committing or pushing. GitHub also runs CI for pull requests targeting `main`. Include the commands run and their actual results. Update tests, documentation, and configuration examples whenever behavior changes. Keep AI.md read-only; implement its requirements and record specification defects in the project-maintained API template issue list.

## Security reports

Do not open a public issue for a vulnerability. Report source-code vulnerabilities through [GitHub private vulnerability reporting](https://github.com/apimgr/wordList/security/advisories/new). See [the security policy](SECURITY.md) for response and disclosure details. The deployed instance contact form is for reports about that deployment, not the source repository.
