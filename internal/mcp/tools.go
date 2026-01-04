package mcp

import (
	"fmt"
	"time"

	"github.com/DandaAkhilReddy/claude_pp/internal/telemetry"
)

const (
	maxFactSize    = 1024 * 1024 // 1MB
	maxMessageSize = 64 * 1024   // 64KB
	maxTagLength   = 100
	maxTagCount    = 50
	staleTimeout   = 5 * time.Minute
)

// handleRemember stores a new fact
func (s *Server) handleRemember(args map[string]interface{}) (interface{}, error) {
	telemetry.TrackMCPTool("remember")

	fact, ok := args["fact"].(string)
	if !ok || fact == "" {
		return nil, fmt.Errorf("fact is required")
	}

	if len(fact) > maxFactSize {
		return nil, fmt.Errorf("fact exceeds maximum size of %d bytes", maxFactSize)
	}

	var tags []string
	if tagsRaw, ok := args["tags"].([]interface{}); ok {
		if len(tagsRaw) > maxTagCount {
			return nil, fmt.Errorf("too many tags (max %d)", maxTagCount)
		}
		for _, t := range tagsRaw {
			if tag, ok := t.(string); ok {
				if len(tag) > maxTagLength {
					return nil, fmt.Errorf("tag exceeds maximum length of %d characters", maxTagLength)
				}
				tags = append(tags, tag)
			}
		}
	}

	id, err := s.store.AddFact(fact, tags, s.workingDir)
	if err != nil {
		return nil, fmt.Errorf("failed to store fact: %w", err)
	}

	return map[string]interface{}{
		"success": true,
		"id":      id,
		"message": "Fact stored successfully",
	}, nil
}

// handleRecall searches for stored facts
func (s *Server) handleRecall(args map[string]interface{}) (interface{}, error) {
	telemetry.TrackMCPTool("recall")

	query, _ := args["query"].(string)

	var tags []string
	if tagsRaw, ok := args["tags"].([]interface{}); ok {
		for _, t := range tagsRaw {
			if tag, ok := t.(string); ok {
				tags = append(tags, tag)
			}
		}
	}

	limit := 20
	if l, ok := args["limit"].(float64); ok {
		limit = int(l)
	}

	facts, err := s.store.QueryFacts(query, tags, "", limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query facts: %w", err)
	}

	if len(facts) == 0 {
		return map[string]interface{}{
			"facts":   []interface{}{},
			"message": "No matching facts found",
		}, nil
	}

	results := make([]map[string]interface{}, len(facts))
	for i, f := range facts {
		results[i] = map[string]interface{}{
			"id":         f.ID,
			"content":    f.Content,
			"tags":       f.Tags,
			"source_dir": f.SourceDir,
			"created_at": f.CreatedAt.Format(time.RFC3339),
		}
	}

	return map[string]interface{}{
		"facts": results,
		"count": len(facts),
	}, nil
}

// handleGetContext retrieves relevant context for the current directory
func (s *Server) handleGetContext(args map[string]interface{}) (interface{}, error) {
	telemetry.TrackMCPTool("get_context")

	// Get local facts for this directory
	localFacts, err := s.store.QueryFacts("", nil, s.workingDir, 50)
	if err != nil {
		return nil, fmt.Errorf("failed to query local facts: %w", err)
	}

	// Get recent global facts from other directories
	allFacts, err := s.store.QueryFacts("", nil, "", 20)
	if err != nil {
		return nil, fmt.Errorf("failed to query global facts: %w", err)
	}

	// Filter out local facts from global results
	var globalFacts []map[string]interface{}
	for _, f := range allFacts {
		if f.SourceDir != s.workingDir {
			globalFacts = append(globalFacts, map[string]interface{}{
				"id":         f.ID,
				"content":    f.Content,
				"tags":       f.Tags,
				"source_dir": f.SourceDir,
				"created_at": f.CreatedAt.Format(time.RFC3339),
			})
		}
	}

	localResults := make([]map[string]interface{}, len(localFacts))
	for i, f := range localFacts {
		localResults[i] = map[string]interface{}{
			"id":         f.ID,
			"content":    f.Content,
			"tags":       f.Tags,
			"created_at": f.CreatedAt.Format(time.RFC3339),
		}
	}

	return map[string]interface{}{
		"working_dir":  s.workingDir,
		"instance_id":  s.instanceID,
		"local_facts":  localResults,
		"global_facts": globalFacts,
		"local_count":  len(localFacts),
		"global_count": len(globalFacts),
	}, nil
}

// handleListInstances returns all running instances
func (s *Server) handleListInstances(args map[string]interface{}) (interface{}, error) {
	telemetry.TrackMCPTool("list_instances")

	// Clean up stale instances first
	if err := s.store.CleanupStaleInstances(staleTimeout); err != nil {
		return nil, fmt.Errorf("failed to cleanup stale instances: %w", err)
	}

	instances, err := s.store.GetInstances()
	if err != nil {
		return nil, fmt.Errorf("failed to get instances: %w", err)
	}

	results := make([]map[string]interface{}, len(instances))
	for i, inst := range instances {
		results[i] = map[string]interface{}{
			"id":             inst.ID,
			"pid":            inst.PID,
			"working_dir":    inst.WorkingDir,
			"started_at":     inst.StartedAt.Format(time.RFC3339),
			"last_heartbeat": inst.LastHeartbeat.Format(time.RFC3339),
			"is_current":     inst.ID == s.instanceID,
		}
	}

	return map[string]interface{}{
		"instances":   results,
		"count":       len(instances),
		"current_id":  s.instanceID,
		"working_dir": s.workingDir,
	}, nil
}

// handleSendMessage sends a message to another instance
func (s *Server) handleSendMessage(args map[string]interface{}) (interface{}, error) {
	telemetry.TrackMCPTool("send_message")

	toInstance, ok := args["to_instance"].(string)
	if !ok || toInstance == "" {
		return nil, fmt.Errorf("to_instance is required")
	}

	message, ok := args["message"].(string)
	if !ok || message == "" {
		return nil, fmt.Errorf("message is required")
	}

	if len(message) > maxMessageSize {
		return nil, fmt.Errorf("message exceeds maximum size of %d bytes", maxMessageSize)
	}

	// Verify target instance exists
	inst, err := s.store.GetInstance(toInstance)
	if err != nil {
		return nil, fmt.Errorf("failed to verify instance: %w", err)
	}
	if inst == nil {
		return nil, fmt.Errorf("instance not found: %s", toInstance)
	}

	id, err := s.store.SendMessage(s.instanceID, toInstance, message)
	if err != nil {
		return nil, fmt.Errorf("failed to send message: %w", err)
	}

	return map[string]interface{}{
		"success":     true,
		"message_id":  id,
		"to_instance": toInstance,
		"to_dir":      inst.WorkingDir,
	}, nil
}

// handleGetMessages retrieves messages for this instance
func (s *Server) handleGetMessages(args map[string]interface{}) (interface{}, error) {
	telemetry.TrackMCPTool("get_messages")

	unreadOnly := true
	if u, ok := args["unread_only"].(bool); ok {
		unreadOnly = u
	}

	messages, err := s.store.GetMessages(s.instanceID, unreadOnly)
	if err != nil {
		return nil, fmt.Errorf("failed to get messages: %w", err)
	}

	// Mark messages as read
	for _, msg := range messages {
		if msg.ReadAt == nil {
			s.store.MarkMessageRead(msg.ID)
		}
	}

	results := make([]map[string]interface{}, len(messages))
	for i, msg := range messages {
		result := map[string]interface{}{
			"id":            msg.ID,
			"from_instance": msg.FromInstance,
			"content":       msg.Content,
			"created_at":    msg.CreatedAt.Format(time.RFC3339),
		}
		if msg.ReadAt != nil {
			result["read_at"] = msg.ReadAt.Format(time.RFC3339)
		}
		results[i] = result
	}

	return map[string]interface{}{
		"messages": results,
		"count":    len(messages),
	}, nil
}
