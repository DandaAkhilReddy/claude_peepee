# Claude PP - Persistent Memory MCP

This project provides Claude PP, an MCP server for persistent memory across Claude Code sessions.

## Available Tools

- **mcp__claude_pp__remember**: Store facts, decisions, and context for later recall
- **mcp__claude_pp__recall**: Search for previously stored facts using keywords or tags
- **mcp__claude_pp__get_context**: Load all relevant context for the current working directory
- **mcp__claude_pp__list_instances**: Discover other running Claude Code instances
- **mcp__claude_pp__send_message**: Send a message to another instance
- **mcp__claude_pp__get_messages**: Retrieve messages from other instances

## Workflow Recommendations

1. **Session Start**: Use `get_context` at the beginning of each session to restore previous state
2. **Important Decisions**: Store architectural decisions and important notes using `remember`
3. **Monorepo Work**: Use `list_instances` and messaging to coordinate between instances
4. **Periodic Checks**: Check for messages from other instances periodically

## Example Usage

```
# Store a fact
mcp__claude_pp__remember(fact="This project uses React 18 with TypeScript", tags=["frontend", "tech-stack"])

# Search for facts
mcp__claude_pp__recall(query="database", limit=10)

# Get all context for current directory
mcp__claude_pp__get_context()

# List running instances
mcp__claude_pp__list_instances()

# Send message to another instance
mcp__claude_pp__send_message(to_instance="abc123", message="API changes complete")

# Get messages
mcp__claude_pp__get_messages(unread_only=true)
```
