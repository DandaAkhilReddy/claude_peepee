# Claude PP - Persistent Memory for Claude Code

**Claude PP** is an MCP (Model Context Protocol) server that gives Claude Code a **persistent brain**. Store facts, decisions, and architectural notes that survive across sessions. Coordinate between multiple Claude Code instances working on different parts of your codebase.

## What It Does

| Feature | Description |
|---------|-------------|
| **Persistent Memory** | Store facts, decisions, and context that survive across Claude Code sessions. No more re-explaining your codebase every time. |
| **Multi-Instance Communication** | Running Claude Code in multiple directories? Instances can discover and message each other. |
| **Automatic Context** | Load relevant context based on your working directory. Start where you left off. |

## Quick Start

### 1. Install

**macOS/Linux:**
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

**Build from Source:**
```bash
git clone https://github.com/DandaAkhilReddy/claude_pp.git
cd claude_pp
make build
```

### 2. Setup for Claude Code

```bash
claude_pp setup
```

This adds Claude PP to your Claude Code configuration. That's it!

### 3. Start Using

In Claude Code, you now have access to these tools:
- `mcp__claude_pp__remember` - Store facts
- `mcp__claude_pp__recall` - Search facts
- `mcp__claude_pp__get_context` - Load all context
- `mcp__claude_pp__list_instances` - Find other instances
- `mcp__claude_pp__send_message` - Message other instances
- `mcp__claude_pp__get_messages` - Read messages

## CLI Commands

Claude PP also works as a standalone CLI:

### Remember - Store Facts

```bash
# Store a simple fact
claude_pp remember "This project uses PostgreSQL for persistence"

# Store with tags for easy filtering
claude_pp remember "API rate limit is 100 req/min" -t api -t config

# Store architectural decisions
claude_pp remember "We chose microservices over monolith for scalability" -t architecture -t decision
```

### Recall - Search Facts

```bash
# Search by keyword
claude_pp recall database
# Output:
# #1 [2026-01-04 04:58]
# Tags: database, architecture
# Dir: /home/user/myproject
# Project uses PostgreSQL for the database

# Filter by tag
claude_pp recall -t api

# Show only facts from current directory
claude_pp recall -l

# Limit results
claude_pp recall -n 5

# Combine filters
claude_pp recall authentication -t security -n 10
```

### Status - Check Your Memory

```bash
claude_pp status
# Output:
# Claude PP Status
# ================
#
# Data directory: /home/user/.claude_pp
# Working directory: /home/user/myproject
#
# Facts:
#   Total: 15
#   Local (this directory): 8
#
# Running instances: 2
#   - abc123: /home/user/myproject/frontend
#   - def456: /home/user/myproject/backend
```

### Instances - Multi-Instance Coordination

```bash
# List all running Claude Code instances
claude_pp instances
# Output:
# Running instances (2):
#
# ID: abc123
#   PID: 12345
#   Dir: /home/user/myproject/frontend
#   Started: 2026-01-04 10:30:00
#   Last heartbeat: 2026-01-04 10:35:00
#
# ID: def456
#   PID: 12346
#   Dir: /home/user/myproject/backend
#   Started: 2026-01-04 10:32:00
#   Last heartbeat: 2026-01-04 10:35:00
```

### Send & Messages - Cross-Instance Communication

```bash
# Send a message to another instance
claude_pp send abc123 "I'm updating the API types, hold off on frontend changes"
# Output: Message #1 sent to abc123 (/home/user/myproject/frontend)

# Check messages for an instance
claude_pp messages abc123
# Output:
# #1 [2026-01-04 10:36] from def456 (unread)
# Backend API changes complete, you can update types now

# Show all messages (including read)
claude_pp messages abc123 -a
```

## Setup Options

### Claude Code (Default)

```bash
# Global setup (recommended)
claude_pp setup

# Project-specific setup
claude_pp setup --project

# Pre-approve all Claude PP commands (no confirmation prompts)
claude_pp setup --allow-all
```

### Other AI Coding Tools

```bash
# OpenCode
claude_pp setup --opencode

# Codex CLI
claude_pp setup --codex

# Gemini CLI
claude_pp setup --gemini
```

## MCP Tools Reference

When using Claude PP through Claude Code, these MCP tools are available:

### `mcp__claude_pp__remember`

Store a fact, decision, or piece of context.

**Parameters:**
- `fact` (required): The information to store
- `tags` (optional): Array of tags for categorization

**Example:**
```
mcp__claude_pp__remember(
  fact="Database migrations run automatically on deploy",
  tags=["database", "deployment"]
)
```

### `mcp__claude_pp__recall`

Search for previously stored facts.

**Parameters:**
- `query` (optional): Search keywords
- `tags` (optional): Filter by tags
- `limit` (optional): Max results (default: 20)

**Example:**
```
mcp__claude_pp__recall(query="authentication", tags=["security"])
```

### `mcp__claude_pp__get_context`

Load all relevant context for the current directory. Returns local facts and recent global facts.

**Parameters:** None

### `mcp__claude_pp__list_instances`

List all running Claude PP instances across directories.

**Parameters:** None

### `mcp__claude_pp__send_message`

Send a message to another instance.

**Parameters:**
- `to_instance` (required): Instance ID from list_instances
- `message` (required): Message content

### `mcp__claude_pp__get_messages`

Get messages sent to this instance.

**Parameters:**
- `unread_only` (optional): Only unread messages (default: true)

## Example Use Cases

### 1. Remembering Architectural Decisions

```bash
claude_pp remember "We use event sourcing for the order system because we need full audit trails" -t architecture -t orders
claude_pp remember "Redis is used for caching with 1hr TTL for user profiles" -t cache -t performance
claude_pp remember "All API endpoints require JWT auth except /health and /metrics" -t api -t security
```

### 2. Project Conventions

```bash
claude_pp remember "File naming: kebab-case for files, PascalCase for components" -t conventions
claude_pp remember "All database columns use snake_case" -t conventions -t database
claude_pp remember "Error messages should be user-friendly, log technical details" -t conventions -t errors
```

### 3. Monorepo Coordination

Terminal 1 (Backend):
```bash
claude_pp instances
# See frontend instance: abc123
claude_pp send abc123 "Deploying new API version, expect brief downtime"
```

Terminal 2 (Frontend):
```bash
claude_pp messages abc123
# See message from backend about deployment
```

### 4. Session Continuity

At the start of each Claude Code session:
```
Use mcp__claude_pp__get_context to load previous context
```

This retrieves all facts stored for the current directory plus recent global facts.

## Data Storage

All data is stored locally in `~/.claude_pp/` using SQLite:

```
~/.claude_pp/
└── claude_pp.db    # SQLite database with facts, instances, messages
```

Your data never leaves your machine.

## Privacy & Telemetry

Claude PP includes optional telemetry to help improve the tool. **No personal data, file contents, or stored facts are ever collected.**

Telemetry is **disabled by default**. If enabled in the future, you can opt out:

```bash
export CLAUDE_PP_NO_TELEMETRY=1
# or
export DO_NOT_TRACK=1
```

## Requirements

- Go 1.21+ (for building from source)
- CGO enabled (for SQLite with FTS5 full-text search)
- Claude Code, OpenCode, Codex CLI, or Gemini CLI

## Troubleshooting

### "command not found: claude_pp"

Add the install directory to your PATH:
```bash
export PATH="$PATH:$HOME/.local/bin"
```

### "failed to open database"

Ensure the data directory exists:
```bash
mkdir -p ~/.claude_pp
```

### MCP tools not showing in Claude Code

Re-run setup:
```bash
claude_pp setup
```

Then restart Claude Code.

## License

MIT License - see [LICENSE](LICENSE) for details.

## Credits

Inspired by [Clauder](https://github.com/MaorBril/clauder) by Maor Bril.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

---

**Claude PP** - Because Claude Code deserves a memory.
