package mcp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/DandaAkhilReddy/claude_pp/internal/store"
)

const (
	MCPVersion = "2024-11-05"
	ServerName = "claude_pp"
	ServerVersion = "0.1.0"
)

// Server implements the MCP protocol over stdio
type Server struct {
	store      store.Store
	instanceID string
	workingDir string
	mu         sync.Mutex
}

// NewServer creates a new MCP server
func NewServer(s store.Store, instanceID, workingDir string) *Server {
	return &Server{
		store:      s,
		instanceID: instanceID,
		workingDir: workingDir,
	}
}

// JSONRPCRequest represents an incoming JSON-RPC request
type JSONRPCRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

// JSONRPCResponse represents an outgoing JSON-RPC response
type JSONRPCResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   *RPCError   `json:"error,omitempty"`
}

// RPCError represents a JSON-RPC error
type RPCError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Run starts the MCP server
func (s *Server) Run() error {
	scanner := bufio.NewScanner(os.Stdin)
	// Increase buffer size for large messages
	buf := make([]byte, 0, 1024*1024)
	scanner.Buffer(buf, 10*1024*1024)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		var req JSONRPCRequest
		if err := json.Unmarshal([]byte(line), &req); err != nil {
			s.sendError(nil, -32700, "Parse error", err.Error())
			continue
		}

		s.handleRequest(&req)
	}

	return scanner.Err()
}

func (s *Server) handleRequest(req *JSONRPCRequest) {
	switch req.Method {
	case "initialize":
		s.handleInitialize(req)
	case "notifications/initialized":
		// No response needed for notifications
	case "tools/list":
		s.handleToolsList(req)
	case "tools/call":
		s.handleToolsCall(req)
	default:
		s.sendError(req.ID, -32601, "Method not found", nil)
	}
}

func (s *Server) handleInitialize(req *JSONRPCRequest) {
	result := map[string]interface{}{
		"protocolVersion": MCPVersion,
		"capabilities": map[string]interface{}{
			"tools": map[string]interface{}{},
		},
		"serverInfo": map[string]interface{}{
			"name":    ServerName,
			"version": ServerVersion,
		},
	}
	s.sendResult(req.ID, result)
}

func (s *Server) handleToolsList(req *JSONRPCRequest) {
	tools := []map[string]interface{}{
		{
			"name":        "remember",
			"description": "Store a fact, decision, or piece of context that should persist across Claude Code sessions. Use this to save important information that you'll need to recall later.",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"fact": map[string]interface{}{
						"type":        "string",
						"description": "The fact, decision, or context to remember",
					},
					"tags": map[string]interface{}{
						"type":        "array",
						"items":       map[string]interface{}{"type": "string"},
						"description": "Optional tags to categorize this fact",
					},
				},
				"required": []string{"fact"},
			},
		},
		{
			"name":        "recall",
			"description": "Search for previously stored facts using keywords or tags. Returns matching facts that were saved with the remember tool.",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"query": map[string]interface{}{
						"type":        "string",
						"description": "Search query to find relevant facts",
					},
					"tags": map[string]interface{}{
						"type":        "array",
						"items":       map[string]interface{}{"type": "string"},
						"description": "Filter by specific tags",
					},
					"limit": map[string]interface{}{
						"type":        "integer",
						"description": "Maximum number of results to return (default: 20)",
					},
				},
			},
		},
		{
			"name":        "get_context",
			"description": "Get all relevant context for the current working directory. This includes local facts stored for this directory and recent global facts from other directories.",
			"inputSchema": map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{},
			},
		},
		{
			"name":        "list_instances",
			"description": "List all running claude_pp instances across different directories. Useful for discovering other active Claude Code sessions.",
			"inputSchema": map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{},
			},
		},
		{
			"name":        "send_message",
			"description": "Send a message to another running claude_pp instance. Use list_instances first to find available instances.",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"to_instance": map[string]interface{}{
						"type":        "string",
						"description": "The ID of the instance to send the message to",
					},
					"message": map[string]interface{}{
						"type":        "string",
						"description": "The message content to send",
					},
				},
				"required": []string{"to_instance", "message"},
			},
		},
		{
			"name":        "get_messages",
			"description": "Get messages sent to this instance from other instances. Messages are marked as read after retrieval.",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"unread_only": map[string]interface{}{
						"type":        "boolean",
						"description": "Only return unread messages (default: true)",
					},
				},
			},
		},
	}

	result := map[string]interface{}{
		"tools": tools,
	}
	s.sendResult(req.ID, result)
}

func (s *Server) handleToolsCall(req *JSONRPCRequest) {
	var params struct {
		Name      string                 `json:"name"`
		Arguments map[string]interface{} `json:"arguments"`
	}

	if err := json.Unmarshal(req.Params, &params); err != nil {
		s.sendError(req.ID, -32602, "Invalid params", err.Error())
		return
	}

	var result interface{}
	var err error

	switch params.Name {
	case "remember":
		result, err = s.handleRemember(params.Arguments)
	case "recall":
		result, err = s.handleRecall(params.Arguments)
	case "get_context":
		result, err = s.handleGetContext(params.Arguments)
	case "list_instances":
		result, err = s.handleListInstances(params.Arguments)
	case "send_message":
		result, err = s.handleSendMessage(params.Arguments)
	case "get_messages":
		result, err = s.handleGetMessages(params.Arguments)
	default:
		s.sendError(req.ID, -32602, "Unknown tool", params.Name)
		return
	}

	if err != nil {
		s.sendError(req.ID, -32603, "Tool execution error", err.Error())
		return
	}

	s.sendResult(req.ID, map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": formatResult(result),
			},
		},
	})
}

func (s *Server) sendResult(id interface{}, result interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()

	resp := JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Result:  result,
	}

	data, _ := json.Marshal(resp)
	fmt.Println(string(data))
}

func (s *Server) sendError(id interface{}, code int, message string, data interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()

	resp := JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Error: &RPCError{
			Code:    code,
			Message: message,
			Data:    data,
		},
	}

	data2, _ := json.Marshal(resp)
	fmt.Println(string(data2))
}

func formatResult(result interface{}) string {
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return fmt.Sprintf("%v", result)
	}
	return string(data)
}
