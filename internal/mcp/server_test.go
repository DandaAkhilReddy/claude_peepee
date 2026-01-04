package mcp

import (
	"encoding/json"
	"testing"
)

func TestJSONRPCRequestParsing(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantID  interface{}
		wantMet string
		wantErr bool
	}{
		{
			name:    "valid initialize request",
			input:   `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}`,
			wantID:  float64(1),
			wantMet: "initialize",
			wantErr: false,
		},
		{
			name:    "valid tools/list request",
			input:   `{"jsonrpc":"2.0","id":2,"method":"tools/list"}`,
			wantID:  float64(2),
			wantMet: "tools/list",
			wantErr: false,
		},
		{
			name:    "string id",
			input:   `{"jsonrpc":"2.0","id":"abc","method":"test"}`,
			wantID:  "abc",
			wantMet: "test",
			wantErr: false,
		},
		{
			name:    "invalid json",
			input:   `{invalid`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req JSONRPCRequest
			err := json.Unmarshal([]byte(tt.input), &req)

			if (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if req.ID != tt.wantID {
					t.Errorf("ID = %v, want %v", req.ID, tt.wantID)
				}
				if req.Method != tt.wantMet {
					t.Errorf("Method = %v, want %v", req.Method, tt.wantMet)
				}
			}
		})
	}
}

func TestJSONRPCResponseSerialization(t *testing.T) {
	tests := []struct {
		name     string
		response JSONRPCResponse
		wantErr  bool
	}{
		{
			name: "success response",
			response: JSONRPCResponse{
				JSONRPC: "2.0",
				ID:      1,
				Result: map[string]interface{}{
					"status": "ok",
				},
			},
			wantErr: false,
		},
		{
			name: "error response",
			response: JSONRPCResponse{
				JSONRPC: "2.0",
				ID:      1,
				Error: &RPCError{
					Code:    -32601,
					Message: "Method not found",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.response)
			if (err != nil) != tt.wantErr {
				t.Errorf("Marshal error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verify we can unmarshal it back
				var resp JSONRPCResponse
				if err := json.Unmarshal(data, &resp); err != nil {
					t.Errorf("Failed to unmarshal response: %v", err)
				}
			}
		})
	}
}

func TestRPCErrorCodes(t *testing.T) {
	tests := []struct {
		name    string
		code    int
		message string
	}{
		{"parse error", -32700, "Parse error"},
		{"invalid request", -32600, "Invalid Request"},
		{"method not found", -32601, "Method not found"},
		{"invalid params", -32602, "Invalid params"},
		{"internal error", -32603, "Internal error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &RPCError{
				Code:    tt.code,
				Message: tt.message,
			}

			if err.Code != tt.code {
				t.Errorf("Code = %d, want %d", err.Code, tt.code)
			}
		})
	}
}

func TestFormatResult(t *testing.T) {
	tests := []struct {
		name   string
		input  interface{}
		wantOK bool
	}{
		{
			name:   "simple map",
			input:  map[string]interface{}{"key": "value"},
			wantOK: true,
		},
		{
			name:   "nested structure",
			input:  map[string]interface{}{"facts": []string{"fact1", "fact2"}},
			wantOK: true,
		},
		{
			name:   "string",
			input:  "simple string",
			wantOK: true,
		},
		{
			name:   "number",
			input:  42,
			wantOK: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatResult(tt.input)
			if (len(result) > 0) != tt.wantOK {
				t.Errorf("formatResult() returned empty for valid input")
			}
		})
	}
}

func TestMCPVersion(t *testing.T) {
	if MCPVersion != "2024-11-05" {
		t.Errorf("MCPVersion = %s, want 2024-11-05", MCPVersion)
	}
}

func TestServerInfo(t *testing.T) {
	if ServerName != "claude_pp" {
		t.Errorf("ServerName = %s, want claude_pp", ServerName)
	}

	if ServerVersion != "0.1.0" {
		t.Errorf("ServerVersion = %s, want 0.1.0", ServerVersion)
	}
}
