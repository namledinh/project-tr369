# Security Policy

## Supported Versions

| Version | Supported          |
| ------- | ------------------ |
| 1.x.x   | :white_check_mark: |

## Reporting a Vulnerability

Please report security vulnerabilities by creating a private security advisory on GitHub or by emailing the maintainers directly.

### What to include in your report:

1. Description of the vulnerability
2. Steps to reproduce the issue
3. Potential impact assessment
4. Suggested fix (if available)

### Response timeline:

- **Acknowledgment**: Within 24 hours
- **Initial assessment**: Within 72 hours  
- **Status update**: Weekly until resolved

## Security Best Practices

This project follows these security practices:

- Regular dependency updates via Dependabot
- Automated security scanning with CodeQL
- Container image scanning
- Secret scanning enabled
- Branch protection rules
- Code review requirements

## Security Scanning

The following automated security tools are integrated:

- **CodeQL**: Static analysis for Go code
- **Gosec**: Go security checker
- **Dependency Review**: Checks for vulnerable dependencies
- **Container Scanning**: Docker image vulnerability scanning