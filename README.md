<p align="center">
  <img src="https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go Version">
  <img src="https://img.shields.io/badge/License-MIT-green?style=for-the-badge" alt="License">
  <img src="https://img.shields.io/badge/MCP-Compatible-purple?style=for-the-badge" alt="MCP Compatible">
  <img src="https://img.shields.io/badge/Platform-macOS%20|%20Linux%20|%20Windows-blue?style=for-the-badge" alt="Platform">
</p>

<h1 align="center">Claude PP</h1>
<h3 align="center">Persistent Memory for Claude Code</h3>

<p align="center">
  <strong>Give Claude Code a brain that remembers across sessions</strong>
</p>

<p align="center">
  <a href="#-features">Features</a> •
  <a href="#-quick-start">Quick Start</a> •
  <a href="#-usage">Usage</a> •
  <a href="#-web-ui">Web UI</a> •
  <a href="#-examples">Examples</a>
</p>

---

## The Problem

Every time you start a new Claude Code session, you have to re-explain:
- Your project architecture
- Coding conventions
- Past decisions and why you made them
- Context from previous sessions

**Claude PP solves this.** It gives Claude Code persistent memory that survives across sessions.

---

## ✨ Features

| Feature | Description |
|---------|-------------|
| 🧠 **Persistent Memory** | Store facts, decisions, and context that survive forever |
| 🔍 **Smart Search** | Full-text search with tag filtering |
| 💬 **Multi-Instance Chat** | Coordinate between Claude instances in different directories |
| 🌐 **Web UI** | Beautiful interface to manage your knowledge base |
| 🔒 **100% Local** | All data stays on your machine |
| ⚡ **Instant Setup** | One command to get started |

---

## 🚀 Quick Start

### Install (30 seconds)

**macOS / Linux:**
```bash
curl -sSL https://raw.githubusercontent.com/DandaAkhilReddy/claude_pp/main/install.sh | sh
```

**Windows (PowerShell):**
```powershell
irm https://raw.githubusercontent.com/DandaAkhilReddy/claude_pp/main/install.ps1 | iex
```

**Go Install:**
```bash
go install github.com/DandaAkhilReddy/claude_pp@latest
```

### Setup for Claude Code

```bash
claude_pp setup
```

**That's it!** Claude Code now has persistent memory.

---

## 📖 Usage

### Store Knowledge

```bash
# Remember a fact
claude_pp remember "We use PostgreSQL 15 with TimescaleDB extension"

# Add tags for easy filtering
claude_pp remember "API rate limit is 100 req/min per user" -t api -t limits

# Store architectural decisions
claude_pp remember "Chose microservices for independent scaling" -t architecture -t decision
```

### Search Knowledge

```bash
# Search by keyword
claude_pp recall database

# Filter by tag
claude_pp recall -t architecture

# Combine search + tags
claude_pp recall authentication -t security

# Limit results
claude_pp recall -n 5
```

### Check Status

```bash
claude_pp status
```

Output:
```
Claude PP Status
================

Data directory: ~/.claude_pp
Working directory: /projects/myapp

Facts:
  Total: 42
  Local (this directory): 15

Running instances: 2
  - a1b2c3: /projects/myapp/frontend
  - d4e5f6: /projects/myapp/backend
```

---

## 🌐 Web UI

Claude PP includes a beautiful web interface to manage your knowledge base.

```bash
claude_pp ui
```

Then open http://localhost:8420 in your browser.

**Features:**
- View and search all stored facts
- Add new facts with tags
- Delete outdated information
- See running instances
- Real-time updates

---

## 🛠️ MCP Tools

When using Claude Code, these tools are automatically available:

| Tool | Description |
|------|-------------|
| `mcp__claude_pp__remember` | Store a fact with optional tags |
| `mcp__claude_pp__recall` | Search stored facts |
| `mcp__claude_pp__get_context` | Load all context for current directory |
| `mcp__claude_pp__list_instances` | Find other running Claude instances |
| `mcp__claude_pp__send_message` | Send message to another instance |
| `mcp__claude_pp__get_messages` | Receive messages from other instances |

---

## 💡 Examples

### 1. Project Setup Memory

```bash
# Store your tech stack
claude_pp remember "Frontend: React 18 + TypeScript + Vite" -t stack -t frontend
claude_pp remember "Backend: Go 1.21 + Gin + GORM" -t stack -t backend
claude_pp remember "Database: PostgreSQL 15 + Redis 7" -t stack -t database
claude_pp remember "Deployment: Docker + Kubernetes on AWS EKS" -t stack -t devops
```

### 2. Coding Conventions

```bash
claude_pp remember "Use kebab-case for file names" -t convention
claude_pp remember "All API responses use { data, error, meta } format" -t convention -t api
claude_pp remember "Tests go in __tests__ folder next to source" -t convention -t testing
```

### 3. Architecture Decisions

```bash
claude_pp remember "ADR-001: Use event sourcing for order history - need full audit trail" -t adr
claude_pp remember "ADR-002: Redis for sessions - need sub-ms latency" -t adr
claude_pp remember "ADR-003: Separate auth service - security isolation" -t adr
```

### 4. Monorepo Coordination

**Terminal 1 (Backend):**
```bash
claude_pp instances
# Shows: frontend instance abc123

claude_pp send abc123 "API contract updated - new field 'metadata' on User"
```

**Terminal 2 (Frontend):**
```bash
claude_pp messages abc123
# Shows: "API contract updated - new field 'metadata' on User"
```

---

## ⚙️ Configuration

### Setup Options

```bash
# Claude Code (default - global)
claude_pp setup

# Claude Code (project only)
claude_pp setup --project

# Pre-approve all commands
claude_pp setup --allow-all

# Other AI tools
claude_pp setup --opencode
claude_pp setup --codex
claude_pp setup --gemini
```

### Environment Variables

| Variable | Description |
|----------|-------------|
| `CLAUDE_PP_NO_TELEMETRY` | Disable telemetry (set to `1`) |
| `DO_NOT_TRACK` | Disable telemetry (standard) |

---

## 📁 Data Storage

All data is stored locally:

```
~/.claude_pp/
└── claude_pp.db    # SQLite database
```

**Your data never leaves your machine.**

---

## 🔧 Build from Source

```bash
git clone https://github.com/DandaAkhilReddy/claude_pp.git
cd claude_pp
make build
./claude_pp version
```

---

## 📋 Commands Reference

| Command | Description |
|---------|-------------|
| `claude_pp remember <fact>` | Store a fact |
| `claude_pp recall [query]` | Search facts |
| `claude_pp status` | Show status |
| `claude_pp instances` | List running instances |
| `claude_pp send <id> <msg>` | Send message |
| `claude_pp messages <id>` | View messages |
| `claude_pp setup` | Configure for AI tools |
| `claude_pp ui` | Open web interface |
| `claude_pp version` | Show version |

---

## 🤝 Contributing

We love contributions! Whether you're fixing bugs, adding features, or improving docs - all help is welcome.

### Quick Start for Contributors

```bash
# 1. Fork and clone
git clone https://github.com/YOUR_USERNAME/claude_pp.git
cd claude_pp

# 2. Install dependencies
go mod download

# 3. Build
make build

# 4. Run tests
make test

# 5. Make your changes and submit a PR!
```

### Ways to Contribute

| Type | Description |
|------|-------------|
| Bug Reports | Found a bug? Open an issue with steps to reproduce |
| Feature Requests | Have an idea? We'd love to hear it |
| Code | Fix bugs, add features, improve performance |
| Documentation | Improve README, add examples, fix typos |
| Testing | Add test cases, improve coverage |

### Good First Issues

New to the project? Look for issues labeled `good first issue` - they're perfect for getting started!

### Development Guidelines

- Follow Go best practices and `gofmt`
- Add tests for new features
- Keep commits focused and descriptive
- Update docs if needed

See [CONTRIBUTING.md](CONTRIBUTING.md) for detailed guidelines.

---

## 🔒 Security

See [SECURITY.md](SECURITY.md) for security policy and reporting vulnerabilities.

---

## 📄 License

MIT License - see [LICENSE](LICENSE) for details.

---

<p align="center">
  <strong>Claude PP</strong> - Because Claude Code deserves a memory
</p>
<p align="center">
  Made with care for the developer community
</p>
