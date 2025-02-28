# Security Policy

## Supported Versions

Intake follows semantic versioning. Currently supported versions:

| Version | Supported          |
| ------- | ------------------ |
| 1.x.x   | :white_check_mark: |
| < 1.0   | :x:                |

## Security Features

Intake provides a foundation for implementing security features:

1. **Panic Handler Support**: Support for custom panic handlers via the `PanicHandler` field
2. **Flexible Middleware System**: Easily add security middleware for:
   - Authentication
   - Authorization
   - Rate limiting
   - CORS
   - Input validation
   - XSS protection
   - CSRF protection
3. **Content Type Headers**: Automatic setting of appropriate content type headers for JSON and XML responses

## Security Best Practices

When using Intake, follow these security best practices:

1. **Always use HTTPS** in production
2. **Implement authentication** before exposing sensitive endpoints
3. **Validate all inputs** using appropriate middleware
4. **Set appropriate timeouts** on the HTTP server
5. **Use panic recovery middleware** in production
6. **Configure CORS** appropriately for your application
7. **Set security headers** using middleware

## Example Security Middleware

```go
func securityHeaders(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("X-Frame-Options", "DENY")
        w.Header().Set("X-Content-Type-Options", "nosniff")
        w.Header().Set("X-XSS-Protection", "1; mode=block")
        w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
        next(w, r)
    }
}
```

## Reporting a Vulnerability

We take security vulnerabilities seriously. If you discover a security vulnerability within Intake, please follow these steps:

1. **Do NOT** open a public GitHub issue
2. Email security@example.com with details about the vulnerability
3. Include:
   - Description of the vulnerability
   - Steps to reproduce
   - Potential impact
   - Suggested fix (if any)

### What to Expect

- **Initial Response**: Within 48 hours
- **Status Update**: Within 1 week
- **Security Fix**: As soon as possible, depending on severity

We appreciate your help in making Intake secure for everyone. Responsible disclosure allows us to fix vulnerabilities before they are publicly known.

## Security Updates

Security updates are released as patch versions. We recommend always using the latest patch version of your current minor release.

Subscribe to GitHub releases to be notified of new security updates.
