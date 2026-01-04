package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/DandaAkhilReddy/claude_peepee/internal/mcp"
	"github.com/DandaAkhilReddy/claude_peepee/internal/store"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var (
	remotePort int
	remoteHost string
)

var serveRemoteCmd = &cobra.Command{
	Use:   "serve-remote",
	Short: "Start the MCP server over HTTP (for Claude Connectors)",
	Long: `Start the Claude PeePee MCP server over HTTP for use with Claude's
Connectors feature. This allows you to connect Claude to your persistent
memory from anywhere.

The server uses Server-Sent Events (SSE) for real-time communication.

Example:
  claude_peepee serve-remote                    # Start on default port 8421
  claude_peepee serve-remote -p 3000            # Start on port 3000
  claude_peepee serve-remote --host 0.0.0.0     # Listen on all interfaces`,
	RunE: runServeRemote,
}

func init() {
	serveRemoteCmd.Flags().IntVarP(&remotePort, "port", "p", 8421, "Port to run the remote MCP server on")
	serveRemoteCmd.Flags().StringVar(&remoteHost, "host", "127.0.0.1", "Host to bind to (use 0.0.0.0 for all interfaces)")
}

// SSE-based MCP handler
type remoteMCPHandler struct {
	store      store.Store
	instanceID string
	workingDir string
	mcpServer  *mcp.Server
	mu         sync.RWMutex
	clients    map[string]chan []byte
}

func newRemoteMCPHandler(s store.Store, instanceID, workingDir string) *remoteMCPHandler {
	return &remoteMCPHandler{
		store:      s,
		instanceID: instanceID,
		workingDir: workingDir,
		mcpServer:  mcp.NewServer(s, instanceID, workingDir),
		clients:    make(map[string]chan []byte),
	}
}

func (h *remoteMCPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Enable CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	path := r.URL.Path

	switch {
	case path == "/sse" && r.Method == "GET":
		h.handleSSE(w, r)
	case path == "/message" && r.Method == "POST":
		h.handleMessage(w, r)
	case path == "/health":
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":     "ok",
			"version":    "0.1.0",
			"instanceID": h.instanceID,
		})
	case path == "/":
		h.handleInfo(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (h *remoteMCPHandler) handleInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	html := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <title>Claude PeePee - Remote MCP Server</title>
    <style>
        body { font-family: -apple-system, sans-serif; max-width: 800px; margin: 50px auto; padding: 20px; background: #1a1a2e; color: #e4e4e7; }
        h1 { color: #667eea; }
        code { background: #2d2d44; padding: 2px 8px; border-radius: 4px; color: #a78bfa; }
        .endpoint { background: #2d2d44; padding: 15px; border-radius: 8px; margin: 10px 0; border-left: 3px solid #667eea; }
        .steps { background: #2d2d44; padding: 20px; border-radius: 8px; margin: 20px 0; }
        .steps ol { margin: 0; padding-left: 20px; }
        .steps li { margin: 10px 0; }
        a { color: #667eea; }
    </style>
</head>
<body>
    <h1>🧠 Claude PeePee Remote MCP Server</h1>
    <p>This is the Claude PeePee MCP server running in remote mode. Connect it to Claude using the Connectors feature!</p>

    <h2>📡 API Endpoints</h2>
    <div class="endpoint">
        <strong>SSE Endpoint:</strong> <code>GET /sse</code><br>
        <small>Connect to receive MCP responses via Server-Sent Events</small>
    </div>
    <div class="endpoint">
        <strong>Message Endpoint:</strong> <code>POST /message?sessionId=xxx</code><br>
        <small>Send MCP requests</small>
    </div>
    <div class="endpoint">
        <strong>Health Check:</strong> <code>GET /health</code><br>
        <small>Check server status</small>
    </div>

    <h2>🔌 How to Connect via Claude Connectors</h2>
    <div class="steps">
        <ol>
            <li>Open <strong>Claude</strong> (Desktop or Web)</li>
            <li>Click the <strong>Settings</strong> icon (gear ⚙️)</li>
            <li>Navigate to <strong>"Connectors"</strong> in the sidebar</li>
            <li>Click <strong>"Add Custom Connector"</strong></li>
            <li>Fill in the details:
                <ul>
                    <li><strong>Name:</strong> <code>Claude PeePee</code></li>
                    <li><strong>URL:</strong> <code>http://%s:%d/sse</code></li>
                </ul>
            </li>
            <li>Click <strong>"Add"</strong> to save</li>
            <li>Start chatting! Claude now has persistent memory! 🎉</li>
        </ol>
    </div>

    <h2>✨ Available Tools</h2>
    <ul>
        <li><strong>remember</strong> - Store facts and context</li>
        <li><strong>recall</strong> - Search stored knowledge</li>
        <li><strong>get_context</strong> - Load all relevant context</li>
        <li><strong>list_instances</strong> - Find other Claude instances</li>
        <li><strong>send_message</strong> - Message other instances</li>
        <li><strong>get_messages</strong> - Receive messages</li>
    </ul>

    <p><small>Instance ID: %s | Working Dir: %s</small></p>
</body>
</html>`, remoteHost, remotePort, h.instanceID, h.workingDir)
	w.Write([]byte(html))
}

func (h *remoteMCPHandler) handleSSE(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "SSE not supported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	clientID := uuid.New().String()
	messageChan := make(chan []byte, 100)

	h.mu.Lock()
	h.clients[clientID] = messageChan
	h.mu.Unlock()

	defer func() {
		h.mu.Lock()
		delete(h.clients, clientID)
		h.mu.Unlock()
		close(messageChan)
	}()

	// Send initial endpoint event
	fmt.Fprintf(w, "event: endpoint\ndata: /message?sessionId=%s\n\n", clientID)
	flusher.Flush()

	ctx := r.Context()
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-messageChan:
			fmt.Fprintf(w, "event: message\ndata: %s\n\n", string(msg))
			flusher.Flush()
		}
	}
}

func (h *remoteMCPHandler) handleMessage(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("sessionId")
	if sessionID == "" {
		http.Error(w, "Missing sessionId", http.StatusBadRequest)
		return
	}

	h.mu.RLock()
	clientChan, exists := h.clients[sessionID]
	h.mu.RUnlock()

	if !exists {
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}

	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Process the MCP request
	response := h.processMCPRequest(request)

	// Send response through SSE
	responseBytes, _ := json.Marshal(response)
	select {
	case clientChan <- responseBytes:
	default:
		// Channel full, skip
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{"status": "accepted"})
}

func (h *remoteMCPHandler) processMCPRequest(request map[string]interface{}) map[string]interface{} {
	method, _ := request["method"].(string)
	id := request["id"]
	params, _ := request["params"].(map[string]interface{})

	response := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      id,
	}

	switch method {
	case "initialize":
		response["result"] = map[string]interface{}{
			"protocolVersion": "2024-11-05",
			"capabilities": map[string]interface{}{
				"tools": map[string]interface{}{},
			},
			"serverInfo": map[string]interface{}{
				"name":    "claude_peepee",
				"version": "0.1.0",
			},
		}

	case "tools/list":
		response["result"] = map[string]interface{}{
			"tools": mcp.GetToolDefinitions(),
		}

	case "tools/call":
		toolName, _ := params["name"].(string)
		toolArgs, _ := params["arguments"].(map[string]interface{})

		result, err := h.mcpServer.ExecuteTool(toolName, toolArgs)
		if err != nil {
			response["error"] = map[string]interface{}{
				"code":    -32603,
				"message": err.Error(),
			}
		} else {
			resultJSON, _ := json.MarshalIndent(result, "", "  ")
			response["result"] = map[string]interface{}{
				"content": []map[string]interface{}{
					{"type": "text", "text": string(resultJSON)},
				},
			}
		}

	default:
		response["error"] = map[string]interface{}{
			"code":    -32601,
			"message": "Method not found",
		}
	}

	return response
}

func runServeRemote(cmd *cobra.Command, args []string) error {
	dataDir := getDataDir()
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}

	s, err := store.NewSQLiteStore(dataDir)
	if err != nil {
		return fmt.Errorf("failed to open store: %w", err)
	}
	defer s.Close()

	workingDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	instanceID := uuid.New().String()[:8]
	pid := os.Getpid()

	if err := s.RegisterInstance(instanceID, pid, workingDir); err != nil {
		return fmt.Errorf("failed to register instance: %w", err)
	}

	// Handle graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		s.UnregisterInstance(instanceID)
		cancel()
		os.Exit(0)
	}()

	// Start heartbeat
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				s.UpdateHeartbeat(instanceID)
			}
		}
	}()

	handler := newRemoteMCPHandler(s, instanceID, workingDir)

	mux := http.NewServeMux()
	mux.Handle("/", handler)

	addr := fmt.Sprintf("%s:%d", remoteHost, remotePort)

	fmt.Printf("\n")
	fmt.Printf("  🧠 Claude PeePee Remote MCP Server\n")
	fmt.Printf("  ───────────────────────────────────\n")
	fmt.Printf("\n")
	fmt.Printf("  MCP Server URL:  http://%s/sse\n", addr)
	fmt.Printf("  Health Check:    http://%s/health\n", addr)
	fmt.Printf("\n")
	fmt.Printf("  ┌─────────────────────────────────────────────────────┐\n")
	fmt.Printf("  │  To connect Claude Desktop/Web:                     │\n")
	fmt.Printf("  │                                                     │\n")
	fmt.Printf("  │  1. Open Claude Settings                            │\n")
	fmt.Printf("  │  2. Go to 'Connectors' section                      │\n")
	fmt.Printf("  │  3. Click 'Add Custom Connector'                    │\n")
	fmt.Printf("  │  4. Enter:                                          │\n")
	fmt.Printf("  │     Name: Claude PeePee                             │\n")
	fmt.Printf("  │     URL:  http://%s/sse            │\n", addr)
	fmt.Printf("  └─────────────────────────────────────────────────────┘\n")
	fmt.Printf("\n")
	fmt.Printf("  Press Ctrl+C to stop\n\n")

	return http.ListenAndServe(addr, mux)
}
