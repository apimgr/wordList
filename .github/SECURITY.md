# Security policy

## Supported versions

Security fixes are provided for the latest stable release and the current `main` branch. Older releases are unsupported; upgrade to the latest stable release before reporting an issue that may already be fixed.

## Reporting a vulnerability

Use [GitHub private vulnerability reporting](https://github.com/apimgr/wordList/security/advisories/new) or the repository's **Security → Report a vulnerability** flow. This is the primary reporting channel. Do not file vulnerabilities as public issues or include exploit details in a public discussion.

For a deployed WordList instance, the [security page](/server/security) links to its current report form and policy. The instance also publishes `/.well-known/security.txt`. Its [contact form](/server/contact?security_id=...) accepts deployment reports, but it is not the channel for vulnerabilities in this source repository. The configured security email is a secondary contact only.

## Response and disclosure

The maintainers aim to acknowledge reports within 48 hours and provide an initial assessment within seven days. We coordinate a fix and disclosure timeline with the reporter, prioritizing active exploitation and severity. Please allow up to 90 days for remediation before public disclosure; we will discuss extensions when a fix needs more time.

Reports must not include generated passphrases, access tokens, or other live credentials. Use a minimal reproduction and redact logs before attaching them.
