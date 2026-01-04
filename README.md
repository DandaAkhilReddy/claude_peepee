# Claude PP - Persistent Memory for Claude Code

Claude PP is an MCP server that provides AI coding tools with **persistent memory** and **multi-instance communication**.

## Features

- **Persistent Memory**: Store facts, decisions, and architectural notes that survive across sessions. No more re-explaining your codebase every time.
- **Multi-Instance Communication**: Running Claude Code in multiple directories? Claude PP lets instances discover and message each other.
- **Automatic Context**: Load relevant context based on your working directory. Start where you left off.

## Installation

### Quick Install (macOS/Linux)

```bash
curl -sSL https://raw.githubusercontent.com/DandaAkhilReddy/claude_pp/main/install.sh | sh
```

### Quick Install (Windows)

```powershell
irm https://raw.githubusercontent.com/DandaAkhilReddy/claude_pp/main/install.ps1 | iex
```

### Go Install

```bash
go install github.com/DandaAkhilReddy/claude_pp@latest
```

### Manual Download

Download the latest binary for your platform from the [Releases](https://github.com/DandaAkhilReddy/claude_pp/releases) page:

- `claude_pp-darwin-arm64` - macOS (Apple Silicon)
- `claude_pp-darwin-amd64` - macOS (Intel)
- `claude_pp-linux-amd64` - Linux (x64)
- `claude_pp-linux-arm64` - Linux (ARM64)
- `claude_pp-windows-amd64.exe` - Windows (x64)

### Build from Source

```bash
git clone https://github.com/DandaAkhilReddy/claude_pp.git
cd claude_pp
make build
```

## Setup

### Claude Code (Default)

```bash
claude_pp setup
```

This configures Claude PP globally. For project-specific configuration:

```bash
claude_pp setup --project
```

To pre-approve all Claude PP commands:

```bash
claude_pp setup --allow-all
```

### OpenCode

```bash
claude_pp setup --opencode
```

### Codex CLI

```bash
claude_pp setup --codex
```

### Gemini CLI

```bash
claude_pp setup --gemini
```

## Usage

### CLI Commands

```bash
# Store a fact
claude_pp remember "This project uses PostgreSQL for persistence"

# Store with tags
claude_pp remember "API rate limit is 100 req/min" -t api -t config

# Search for facts
claude_pp recall database
claude_pp recall -t api

# List running instances
claude_pp instances

# Send message to another instance
claude_pp send <instance-id> "Found a bug in the auth module"

# View messages
claude_pp messages <instance-id>

# Check status
claude_pp status
```

### MCP Tools (in Claude Code)

When using Claude Code, these tools are available:

- `mcp__claude_pp__remember` - Store facts, decisions, and context
- `mcp__claude_pp__recall` - Search for previously stored facts
- `mcp__claude_pp__get_context` - Load all relevant context for current directory
- `mcp__claude_pp__list_instances` - Find other running Claude Code instances
- `mcp__claude_pp__send_message` - Send messages to other instances
- `mcp__claude_pp__get_messages` - Retrieve messages from other instances

## Example Use Cases

### Store Architectural Decisions

```
remember "We chose SQLite over PostgreSQL for simplicity - this is a single-user CLI tool"
remember "Authentication uses JWT tokens stored in httpOnly cookies" -t auth -t security
```

### Project Conventions

```
remember "All API endpoints follow REST conventions with /api/v1 prefix" -t api
remember "Use kebab-case for file names, camelCase for variables" -t style
```

### Monorepo Coordination

```
# In frontend directory
list_instances
send_message abc123 "I'm updating the API types, hold off on backend changes"

# In backend directory
get_messages
```

## Data Storage

All data is stored locally in `~/.claude_pp/` using SQLite. Your facts, decisions, and messages never leave your machine.

## Privacy & Telemetry

Claude PP collects anonymous usage data to help improve the tool. **No personal data, file contents, or facts are ever collected.**

To disable telemetry:

```bash
export CLAUDE_PP_NO_TELEMETRY=1
# or
export DO_NOT_TRACK=1
```

## License

MIT License - see [LICENSE](LICENSE) for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Credits

Inspired by [Clauder](https://github.com/MaorBril/clauder) by Maor Bril.
