package store

import "time"

// Fact represents a stored piece of information
type Fact struct {
	ID        int64     `json:"id"`
	Content   string    `json:"content"`
	Tags      []string  `json:"tags,omitempty"`
	SourceDir string    `json:"source_dir,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Instance represents a running claude_pp instance
type Instance struct {
	ID            string    `json:"id"`
	PID           int       `json:"pid"`
	WorkingDir    string    `json:"working_dir"`
	StartedAt     time.Time `json:"started_at"`
	LastHeartbeat time.Time `json:"last_heartbeat"`
}

// Message represents a message between instances
type Message struct {
	ID           int64      `json:"id"`
	FromInstance string     `json:"from_instance"`
	ToInstance   string     `json:"to_instance"`
	Content      string     `json:"content"`
	CreatedAt    time.Time  `json:"created_at"`
	ReadAt       *time.Time `json:"read_at,omitempty"`
}

// Store defines the interface for data persistence
type Store interface {
	// Facts
	AddFact(content string, tags []string, sourceDir string) (int64, error)
	QueryFacts(query string, tags []string, sourceDir string, limit int) ([]Fact, error)
	GetFact(id int64) (*Fact, error)
	DeleteFact(id int64) error
	CountFacts(sourceDir string) (total int, local int, err error)

	// Instances
	RegisterInstance(id string, pid int, workingDir string) error
	UpdateHeartbeat(id string) error
	UnregisterInstance(id string) error
	GetInstances() ([]Instance, error)
	GetInstance(id string) (*Instance, error)
	CleanupStaleInstances(maxAge time.Duration) error

	// Messages
	SendMessage(fromInstance, toInstance, content string) (int64, error)
	GetMessages(instanceID string, unreadOnly bool) ([]Message, error)
	MarkMessageRead(id int64) error

	// Lifecycle
	Close() error
}
