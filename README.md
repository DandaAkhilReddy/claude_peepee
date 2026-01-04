<p align="center">
  <img src="https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go Version">
  <img src="https://img.shields.io/badge/License-MIT-green?style=for-the-badge" alt="License">
  <img src="https://img.shields.io/badge/MCP-Compatible-purple?style=for-the-badge" alt="MCP Compatible">
  <img src="https://img.shields.io/badge/Platform-macOS%20|%20Linux%20|%20Windows-blue?style=for-the-badge" alt="Platform">
  <img src="https://img.shields.io/badge/Tokens-Saved%20💰-gold?style=for-the-badge" alt="Token Efficient">
</p>

<h1 align="center">🧠 Claude PeePee</h1>
<h3 align="center">Persistent Memory for Claude Code</h3>

<p align="center">
  <strong>Give Claude Code a brain that remembers across sessions</strong>
</p>

<p align="center">
  <em>"Break down large context into small, reusable pieces"</em>
</p>

<p align="center">
  <a href="#-why-peepee">Why PeePee?</a> •
  <a href="#-features">Features</a> •
  <a href="#-quick-start">Quick Start</a> •
  <a href="#-token-savings">Token Savings</a> •
  <a href="#-examples">Examples</a>
</p>

---

## 🚽 Why "PeePee"?

Think about how a child's body works: food goes in, the body **absorbs the useful nutrients**, and the waste (pee-pee 🚽) gets flushed out. Simple, efficient, natural!

**Claude PeePee works the same way for your AI conversations:**

```
┌─────────────────────────────────────────────────────────────┐
│                    YOUR PROJECT CONTEXT                      │
│  ┌─────────────────────────────────────────────────────┐    │
│  │  Massive codebase, hundreds of decisions,            │    │
│  │  conventions, architecture notes, past discussions... │    │
│  └─────────────────────────────────────────────────────┘    │
│                           │                                  │
│                           ▼                                  │
│                  ┌─────────────────┐                        │
│                  │  Claude PeePee  │                        │
│                  │    (Filter)     │                        │
│                  └─────────────────┘                        │
│                     │           │                            │
│           ┌─────────┘           └─────────┐                  │
│           ▼                               ▼                  │
│   🚽 FLUSH OUT                    🧠 KEEP & STORE            │
│   ─────────────                   ─────────────              │
│   • Redundant context             • Key decisions            │
│   • Repeated explanations         • Architecture notes       │
│   • Token waste                   • Important facts          │
│   • Session bloat                 • Tagged knowledge         │
└─────────────────────────────────────────────────────────────┘
```

**The result?** Only the **useful nutrients** (relevant facts) flow into each conversation, while the **waste** (redundant tokens) gets flushed away!

- **Filter out** unnecessary repetition
- **Absorb** only what matters into persistent memory
- **Retrieve** just the relevant context when needed
- **Save tokens** = Save money 💰

> *"Flush the waste, keep the knowledge!"*

---

## 💰 Token Savings

**The Real Problem: Token Waste**

| Without Claude PeePee | With Claude PeePee |
|-------------------|----------------|
| Re-explain project structure every session | Stored once, recalled instantly |
| Copy-paste same conventions repeatedly | Tagged and searchable |
| Lose context between sessions | Persistent forever |
| ~2000+ tokens wasted per session | ~50 tokens to recall |

**Example Savings:**
```
Traditional approach:
  "Our project uses React 18 with TypeScript, Vite for bundling,
   PostgreSQL database, Redis caching, Docker deployment..."
  = 500+ tokens EVERY session

With Claude PeePee:
  claude_peepee recall -t stack
  = 50 tokens (just the search) + relevant results only
```

**Estimated savings: 80-90% fewer tokens on context!**

---

## 🎯 The Problem

Every time you start a new Claude Code session, you have to re-explain:
- Your project architecture
- Coding conventions
- Past decisions and why you made them
- Context from previous sessions

**Claude PeePee solves this.** It gives Claude Code persistent memory that survives across sessions.

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

## 🔬 How It Works

### Architecture Overview

```
┌─────────────────────────────────────────────────────────────────────┐
│                         Claude Code                                  │
│                              │                                       │
│                     MCP Protocol (JSON-RPC 2.0)                     │
│                         stdin/stdout                                 │
│                              │                                       │
│                              ▼                                       │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │                    Claude PeePee MCP Server                      │   │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐  │   │
│  │  │  Remember   │  │   Recall    │  │  Instance Messaging │  │   │
│  │  │  (Store)    │  │  (Search)   │  │   (Coordination)    │  │   │
│  │  └─────────────┘  └─────────────┘  └─────────────────────┘  │   │
│  │                         │                                    │   │
│  │                         ▼                                    │   │
│  │  ┌─────────────────────────────────────────────────────┐    │   │
│  │  │              SQLite Database (FTS5)                 │    │   │
│  │  │  ~/.claude_peepee/claude_peepee.db                  │    │   │
│  │  └─────────────────────────────────────────────────────┘    │   │
│  └─────────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────────┘
```

### Step-by-Step Flow

1. **You start Claude Code** → Claude PeePee MCP server starts automatically
2. **Claude needs context** → Calls `get_context` to load relevant facts
3. **You make decisions** → Claude stores them using `remember`
4. **Next session** → Previous facts are instantly available via `recall`

### What Data Is Saved

Claude PeePee stores everything in a local SQLite database with full-text search capabilities:

| Data Type | What's Stored | Example |
|-----------|--------------|---------|
| **Facts** | Any text information you want to remember | "API uses JWT authentication with RS256" |
| **Tags** | Comma-separated labels for filtering | "api, security, authentication" |
| **Source Directory** | Where the fact was created | "/home/user/myproject" |
| **Timestamps** | When created and last updated | "2024-01-15 10:30:00" |
| **Instances** | Running Claude Code sessions | ID, PID, working directory, heartbeat |
| **Messages** | Communication between instances | From, To, Content, Read status |

### Database Schema

```sql
-- Facts table (your persistent memory)
facts (
    id INTEGER PRIMARY KEY,
    content TEXT,           -- The actual fact/knowledge
    tags TEXT,              -- Comma-separated tags
    source_dir TEXT,        -- Directory where fact was created
    created_at DATETIME,
    updated_at DATETIME
)

-- Full-text search index (for fast keyword search)
facts_fts (content, tags)   -- FTS5 virtual table

-- Running instances (for multi-instance coordination)
instances (
    id TEXT PRIMARY KEY,    -- Unique instance ID
    pid INTEGER,            -- Process ID
    working_dir TEXT,       -- Current working directory
    started_at DATETIME,
    last_heartbeat DATETIME
)

-- Messages between instances
messages (
    id INTEGER PRIMARY KEY,
    from_instance TEXT,
    to_instance TEXT,
    content TEXT,
    created_at DATETIME,
    read_at DATETIME        -- NULL if unread
)
```

### The MCP Protocol

Claude PeePee uses the **Model Context Protocol (MCP)** - a standardized way for AI assistants to access external tools:

```
Claude Code                    Claude PeePee Server
    │                                │
    │ ─── initialize request ───────>│
    │ <── capabilities response ─────│
    │                                │
    │ ─── tools/list request ───────>│
    │ <── available tools ───────────│
    │                                │
    │ ─── tools/call (remember) ────>│
    │ <── success response ──────────│
    │                                │
    │ ─── tools/call (recall) ──────>│
    │ <── matching facts ────────────│
```

**Protocol Details:**
- Transport: `stdin/stdout` (fast, no network overhead)
- Format: JSON-RPC 2.0
- Tools exposed: 6 (remember, recall, get_context, list_instances, send_message, get_messages)

### How Token Savings Work

**Traditional approach** (expensive):
```
Every session, you type:
"This project uses React 18 with TypeScript, we follow these conventions..."
= 500+ tokens × every session = 💸💸💸
```

**With Claude PeePee** (efficient):
```
Session 1: claude_peepee remember "React 18 + TypeScript project" -t stack
           (stored once: ~10 tokens)

Session 2+: Claude calls get_context automatically
            Only relevant facts loaded (~50 tokens)

Savings: 90% fewer tokens!
```

### Data Location

```
~/.claude_peepee/
└── claude_peepee.db      # All your data in one SQLite file
    ├── facts             # Your stored knowledge
    ├── facts_fts         # Full-text search index
    ├── instances         # Running session registry
    └── messages          # Inter-instance messages
```

**Privacy:** All data stays 100% local. Nothing is sent to external servers.

---

## 🚀 Quick Start

### Install (30 seconds)

**macOS / Linux:**
```bash
curl -sSL https://raw.githubusercontent.com/DandaAkhilReddy/claude_peepee/main/install.sh | sh
```

**Windows (PowerShell):**
```powershell
irm https://raw.githubusercontent.com/DandaAkhilReddy/claude_peepee/main/install.ps1 | iex
```

**Go Install:**
```bash
go install github.com/DandaAkhilReddy/claude_peepee@latest
```

### Setup for Claude Code

```bash
claude_peepee setup
```

**That's it!** Claude Code now has persistent memory.

---

## 📖 Usage

### Store Knowledge

```bash
# Remember a fact
claude_peepee remember "We use PostgreSQL 15 with TimescaleDB extension"

# Add tags for easy filtering
claude_peepee remember "API rate limit is 100 req/min per user" -t api -t limits

# Store architectural decisions
claude_peepee remember "Chose microservices for independent scaling" -t architecture -t decision
```

### Search Knowledge

```bash
# Search by keyword
claude_peepee recall database

# Filter by tag
claude_peepee recall -t architecture

# Combine search + tags
claude_peepee recall authentication -t security

# Limit results
claude_peepee recall -n 5
```

### Check Status

```bash
claude_peepee status
```

Output:
```
Claude PeePee Status
================

Data directory: ~/.claude_peepee
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

Claude PeePee includes a beautiful web interface to manage your knowledge base.

```bash
claude_peepee ui
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
| `mcp__claude_peepee__remember` | Store a fact with optional tags |
| `mcp__claude_peepee__recall` | Search stored facts |
| `mcp__claude_peepee__get_context` | Load all context for current directory |
| `mcp__claude_peepee__list_instances` | Find other running Claude instances |
| `mcp__claude_peepee__send_message` | Send message to another instance |
| `mcp__claude_peepee__get_messages` | Receive messages from other instances |

---

## 💡 Examples

### 1. Project Setup Memory

```bash
# Store your tech stack
claude_peepee remember "Frontend: React 18 + TypeScript + Vite" -t stack -t frontend
claude_peepee remember "Backend: Go 1.21 + Gin + GORM" -t stack -t backend
claude_peepee remember "Database: PostgreSQL 15 + Redis 7" -t stack -t database
claude_peepee remember "Deployment: Docker + Kubernetes on AWS EKS" -t stack -t devops
```

### 2. Coding Conventions

```bash
claude_peepee remember "Use kebab-case for file names" -t convention
claude_peepee remember "All API responses use { data, error, meta } format" -t convention -t api
claude_peepee remember "Tests go in __tests__ folder next to source" -t convention -t testing
```

### 3. Architecture Decisions

```bash
claude_peepee remember "ADR-001: Use event sourcing for order history - need full audit trail" -t adr
claude_peepee remember "ADR-002: Redis for sessions - need sub-ms latency" -t adr
claude_peepee remember "ADR-003: Separate auth service - security isolation" -t adr
```

### 4. Monorepo Coordination

**Terminal 1 (Backend):**
```bash
claude_peepee instances
# Shows: frontend instance abc123

claude_peepee send abc123 "API contract updated - new field 'metadata' on User"
```

**Terminal 2 (Frontend):**
```bash
claude_peepee messages abc123
# Shows: "API contract updated - new field 'metadata' on User"
```

---

## ⚙️ Configuration

### Setup Options

```bash
# Claude Code (default - global)
claude_peepee setup

# Claude Code (project only)
claude_peepee setup --project

# Pre-approve all commands
claude_peepee setup --allow-all

# Other AI tools
claude_peepee setup --opencode
claude_peepee setup --codex
claude_peepee setup --gemini
```

### Environment Variables

| Variable | Description |
|----------|-------------|
| `CLAUDE_PEEPEE_NO_TELEMETRY` | Disable telemetry (set to `1`) |
| `DO_NOT_TRACK` | Disable telemetry (standard) |

---

## 📁 Data Storage

All data is stored locally:

```
~/.claude_peepee/
└── claude_peepee.db    # SQLite database
```

**Your data never leaves your machine.**

---

## 🔧 Build from Source

```bash
git clone https://github.com/DandaAkhilReddy/claude_peepee.git
cd claude_peepee
make build
./claude_peepee version
```

---

## 📋 Commands Reference

| Command | Description |
|---------|-------------|
| `claude_peepee remember <fact>` | Store a fact |
| `claude_peepee recall [query]` | Search facts |
| `claude_peepee status` | Show status |
| `claude_peepee instances` | List running instances |
| `claude_peepee send <id> <msg>` | Send message |
| `claude_peepee messages <id>` | View messages |
| `claude_peepee setup` | Configure for AI tools |
| `claude_peepee ui` | Open web interface |
| `claude_peepee version` | Show version |

---

## 🤝 Contributing

We love contributions! Whether you're fixing bugs, adding features, or improving docs - all help is welcome.

### Quick Start for Contributors

```bash
# 1. Fork and clone
git clone https://github.com/YOUR_USERNAME/claude_peepee.git
cd claude_peepee

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
  <strong>🧠 Claude PeePee</strong> - Because Claude Code deserves a memory
</p>
<p align="center">
  <em>Break it down. Store it smart. Save those tokens! 💰</em>
</p>
<p align="center">
  Made with 💜 for developers who hate repeating themselves
</p>
