# Security Policy

## Reporting Security Vulnerabilities

If you discover a security vulnerability in Claude PP, please report it responsibly:

1. **DO NOT** open a public GitHub issue for security vulnerabilities
2. Email the maintainer directly with details of the vulnerability
3. Include steps to reproduce the issue
4. Allow reasonable time for a fix before public disclosure

## Security Measures

### Data Storage

- All data is stored **locally** in `~/.claude_pp/claude_pp.db`
- SQLite database with WAL mode for safe concurrent access
- No data is transmitted to external servers
- No cloud storage or synchronization

### SQL Injection Prevention

- All database queries use **parameterized statements**
- FTS5 queries are sanitized to prevent injection attacks
- User input is never interpolated directly into SQL

### Input Validation

- Fact content limited to **1MB** maximum
- Message content limited to **64KB** maximum
- Tag length limited to **100 characters**
- Maximum **50 tags** per fact
- Query results limited to **1000 entries**

### Privacy

- **No telemetry data is collected** (telemetry is disabled by default)
- No personal information is stored or transmitted
- Machine identification uses one-way SHA256 hash (if telemetry enabled)
- Users can opt-out with `CLAUDE_PP_NO_TELEMETRY=1` or `DO_NOT_TRACK=1`

### MCP Protocol Security

- Communication happens over stdin/stdout (local only)
- JSON-RPC 2.0 protocol with proper error handling
- No network listeners or exposed ports
- Instance communication is local (same machine only)

## Best Practices for Users

1. **Keep the binary updated** - Install latest releases for security fixes
2. **Protect your data directory** - Ensure `~/.claude_pp/` has appropriate permissions
3. **Don't share database files** - They may contain sensitive project information
4. **Review stored facts periodically** - Remove any sensitive information

## Configuration Security

### Recommended File Permissions

```bash
# Data directory
chmod 700 ~/.claude_pp

# Database file
chmod 600 ~/.claude_pp/claude_pp.db
```

### Environment Variables

| Variable | Purpose | Default |
|----------|---------|---------|
| `CLAUDE_PP_NO_TELEMETRY` | Disable telemetry | - |
| `DO_NOT_TRACK` | Disable telemetry (standard) | - |
| `HOME` | User home directory | System default |

## Secure Development

When contributing to Claude PP:

1. Never commit secrets, API keys, or credentials
2. Use parameterized SQL queries
3. Validate and sanitize all user input
4. Follow the principle of least privilege
5. Add tests for security-critical code paths

## Third-Party Dependencies

| Dependency | Purpose | Security Notes |
|------------|---------|----------------|
| `go-sqlite3` | Database | CGO-based, audited SQLite |
| `cobra` | CLI framework | No security concerns |
| `uuid` | Instance IDs | Cryptographically random |
| `go-toml` | Config parsing | Safe parsing only |

## Version History

| Version | Security Updates |
|---------|-----------------|
| 0.1.0 | Initial release with security measures |

---

**Last Updated**: January 2025
