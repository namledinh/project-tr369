# Security Policy

## Supported Versions

The following versions are currently being supported with security updates:

| Version | Supported          |
| ------- | ------------------ |
| 1.x.x   | :white_check_mark: |
| < 1.0   | :x:                |

## Reporting a Vulnerability

We take security vulnerabilities seriously. If you discover a security vulnerability, please report it to us as follows:

### Where to Report

- **Email**: [security@yourcompany.com](mailto:security@yourcompany.com)
- **GitHub Security Advisories**: [Create a security advisory](https://github.com/namledinh/usp-management-backend/security/advisories/new)

### What to Include

When reporting a security vulnerability, please include the following information:

1. **Description**: A clear description of the vulnerability
2. **Steps to Reproduce**: Detailed steps to reproduce the issue
3. **Impact**: What an attacker could achieve by exploiting this vulnerability
4. **Affected Versions**: Which versions of the software are affected
5. **Suggested Fix**: If you have suggestions for how to fix the issue

### Response Timeline

- **Initial Response**: We will acknowledge receipt of your vulnerability report within 48 hours
- **Investigation**: We will investigate and validate the vulnerability within 7 days
- **Resolution**: We aim to resolve critical vulnerabilities within 30 days
- **Disclosure**: We follow responsible disclosure practices

### Security Best Practices

When contributing to this project, please follow these security best practices:

1. **Dependencies**: Keep dependencies up to date
2. **Secrets**: Never commit secrets, API keys, or passwords
3. **Input Validation**: Always validate and sanitize user inputs
4. **Authentication**: Implement proper authentication and authorization
5. **Logging**: Avoid logging sensitive information
6. **Error Handling**: Don't expose sensitive information in error messages

### Security Features

This project implements the following security features:

- Input validation and sanitization
- SQL injection prevention using parameterized queries
- Authentication and authorization middleware
- Rate limiting
- CORS configuration
- Security headers
- Dependency vulnerability scanning

## Bug Bounty Program

Currently, we do not have a formal bug bounty program. However, we appreciate security researchers who responsibly disclose vulnerabilities and will acknowledge their contributions.

## Contact

For any questions regarding this security policy, please contact:
- Email: [security@yourcompany.com](mailto:security@yourcompany.com)
- GitHub Issues: [Create an issue](https://github.com/namledinh/usp-management-backend/issues/new)