# Security Policy

## Reporting a Vulnerability

If you discover a security vulnerability, please report it by emailing the maintainer directly rather than opening a public issue.

**Do not disclose security vulnerabilities publicly until they have been addressed.**

## Supported Versions

| Version | Supported          |
| ------- | ------------------ |
| 2.x.x   | :white_check_mark: |
| 1.x.x   | :x:                |

## Security Best Practices

When using Flux:

- Keep your FluxFile in version control
- Review remote execution commands carefully
- Use environment variables for sensitive data
- Validate lockfile integrity with `flux --check-lock`
