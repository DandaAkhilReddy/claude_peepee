package ui

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/DandaAkhilReddy/claude_peepee/internal/store"
)

// Server represents the web UI server
type Server struct {
	store store.Store
	port  int
}

// NewServer creates a new UI server
func NewServer(s store.Store, port int) *Server {
	return &Server{
		store: s,
		port:  port,
	}
}

// Start starts the web UI server
func (s *Server) Start() error {
	mux := http.NewServeMux()

	// Serve the main page
	mux.HandleFunc("/", s.handleIndex)

	// API endpoints
	mux.HandleFunc("/api/facts", s.handleFacts)
	mux.HandleFunc("/api/facts/", s.handleFactByID)
	mux.HandleFunc("/api/instances", s.handleInstances)
	mux.HandleFunc("/api/status", s.handleStatus)

	addr := fmt.Sprintf(":%d", s.port)
	fmt.Printf("\n  Claude PeePee Web UI\n")
	fmt.Printf("  ─────────────────\n")
	fmt.Printf("  Local:   http://localhost%s\n\n", addr)
	fmt.Printf("  Press Ctrl+C to stop\n\n")

	return http.ListenAndServe(addr, mux)
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(indexHTML))
}

func (s *Server) handleFacts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		return
	}

	switch r.Method {
	case "GET":
		query := r.URL.Query().Get("q")
		tagsParam := r.URL.Query().Get("tags")
		var tags []string
		if tagsParam != "" {
			tags = strings.Split(tagsParam, ",")
		}

		facts, err := s.store.QueryFacts(query, tags, "", 100)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"facts": facts,
			"count": len(facts),
		})

	case "POST":
		var req struct {
			Content string   `json:"content"`
			Tags    []string `json:"tags"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		id, err := s.store.AddFact(req.Content, req.Tags, "")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"id":      id,
			"success": true,
		})

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleFactByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, DELETE, OPTIONS")

	if r.Method == "OPTIONS" {
		return
	}

	// Extract ID from path
	path := strings.TrimPrefix(r.URL.Path, "/api/facts/")
	id, err := strconv.ParseInt(path, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case "GET":
		fact, err := s.store.GetFact(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if fact == nil {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(fact)

	case "DELETE":
		if err := s.store.DeleteFact(id); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
		})

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleInstances(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	s.store.CleanupStaleInstances(5 * time.Minute)

	instances, err := s.store.GetInstances()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"instances": instances,
		"count":     len(instances),
	})
}

func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	total, _, _ := s.store.CountFacts("")
	s.store.CleanupStaleInstances(5 * time.Minute)
	instances, _ := s.store.GetInstances()

	json.NewEncoder(w).Encode(map[string]interface{}{
		"total_facts":    total,
		"instances":      len(instances),
		"version":        "0.1.0",
	})
}

const indexHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Claude PeePee - Knowledge Base</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, sans-serif;
            background: linear-gradient(135deg, #1a1a2e 0%, #16213e 100%);
            min-height: 100vh;
            color: #e4e4e7;
        }

        .container {
            max-width: 1200px;
            margin: 0 auto;
            padding: 2rem;
        }

        header {
            text-align: center;
            margin-bottom: 2rem;
        }

        h1 {
            font-size: 2.5rem;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            -webkit-background-clip: text;
            -webkit-text-fill-color: transparent;
            margin-bottom: 0.5rem;
        }

        .subtitle {
            color: #a1a1aa;
            font-size: 1.1rem;
        }

        .stats {
            display: flex;
            gap: 1rem;
            justify-content: center;
            margin-bottom: 2rem;
        }

        .stat-card {
            background: rgba(255, 255, 255, 0.05);
            border: 1px solid rgba(255, 255, 255, 0.1);
            border-radius: 12px;
            padding: 1.5rem 2rem;
            text-align: center;
        }

        .stat-number {
            font-size: 2rem;
            font-weight: 700;
            color: #667eea;
        }

        .stat-label {
            color: #a1a1aa;
            font-size: 0.9rem;
        }

        .search-box {
            display: flex;
            gap: 1rem;
            margin-bottom: 2rem;
        }

        input[type="text"] {
            flex: 1;
            padding: 1rem 1.5rem;
            border: 1px solid rgba(255, 255, 255, 0.1);
            border-radius: 12px;
            background: rgba(255, 255, 255, 0.05);
            color: #e4e4e7;
            font-size: 1rem;
        }

        input[type="text"]:focus {
            outline: none;
            border-color: #667eea;
        }

        input[type="text"]::placeholder {
            color: #71717a;
        }

        button {
            padding: 1rem 2rem;
            border: none;
            border-radius: 12px;
            font-size: 1rem;
            cursor: pointer;
            transition: all 0.2s;
        }

        .btn-primary {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
        }

        .btn-primary:hover {
            transform: translateY(-2px);
            box-shadow: 0 4px 20px rgba(102, 126, 234, 0.4);
        }

        .btn-danger {
            background: rgba(239, 68, 68, 0.2);
            color: #ef4444;
            padding: 0.5rem 1rem;
            font-size: 0.9rem;
        }

        .btn-danger:hover {
            background: rgba(239, 68, 68, 0.3);
        }

        .facts-grid {
            display: grid;
            gap: 1rem;
        }

        .fact-card {
            background: rgba(255, 255, 255, 0.05);
            border: 1px solid rgba(255, 255, 255, 0.1);
            border-radius: 12px;
            padding: 1.5rem;
            transition: all 0.2s;
        }

        .fact-card:hover {
            border-color: rgba(102, 126, 234, 0.5);
            transform: translateY(-2px);
        }

        .fact-header {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-bottom: 1rem;
        }

        .fact-id {
            color: #667eea;
            font-weight: 600;
        }

        .fact-date {
            color: #71717a;
            font-size: 0.9rem;
        }

        .fact-content {
            font-size: 1.1rem;
            line-height: 1.6;
            margin-bottom: 1rem;
        }

        .fact-tags {
            display: flex;
            gap: 0.5rem;
            flex-wrap: wrap;
        }

        .tag {
            background: rgba(102, 126, 234, 0.2);
            color: #667eea;
            padding: 0.25rem 0.75rem;
            border-radius: 20px;
            font-size: 0.85rem;
        }

        .fact-source {
            color: #71717a;
            font-size: 0.85rem;
            margin-top: 0.5rem;
        }

        .add-form {
            background: rgba(255, 255, 255, 0.05);
            border: 1px solid rgba(255, 255, 255, 0.1);
            border-radius: 12px;
            padding: 1.5rem;
            margin-bottom: 2rem;
        }

        .add-form h3 {
            margin-bottom: 1rem;
            color: #667eea;
        }

        .form-group {
            margin-bottom: 1rem;
        }

        textarea {
            width: 100%;
            padding: 1rem;
            border: 1px solid rgba(255, 255, 255, 0.1);
            border-radius: 12px;
            background: rgba(255, 255, 255, 0.05);
            color: #e4e4e7;
            font-size: 1rem;
            min-height: 100px;
            resize: vertical;
        }

        textarea:focus {
            outline: none;
            border-color: #667eea;
        }

        .empty-state {
            text-align: center;
            padding: 4rem 2rem;
            color: #71717a;
        }

        .empty-state h3 {
            margin-bottom: 0.5rem;
            color: #a1a1aa;
        }

        .tabs {
            display: flex;
            gap: 1rem;
            margin-bottom: 2rem;
        }

        .tab {
            padding: 0.75rem 1.5rem;
            border-radius: 8px;
            cursor: pointer;
            transition: all 0.2s;
            background: rgba(255, 255, 255, 0.05);
            border: 1px solid rgba(255, 255, 255, 0.1);
        }

        .tab:hover, .tab.active {
            background: rgba(102, 126, 234, 0.2);
            border-color: #667eea;
        }

        .instances-list {
            display: grid;
            gap: 1rem;
        }

        .instance-card {
            background: rgba(255, 255, 255, 0.05);
            border: 1px solid rgba(255, 255, 255, 0.1);
            border-radius: 12px;
            padding: 1.5rem;
        }

        .instance-id {
            font-family: monospace;
            color: #667eea;
            font-size: 1.1rem;
        }

        .instance-dir {
            color: #a1a1aa;
            margin-top: 0.5rem;
        }

        @media (max-width: 768px) {
            .stats {
                flex-direction: column;
            }

            .search-box {
                flex-direction: column;
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <header>
            <h1>🧠 Claude PeePee</h1>
            <p class="subtitle">Persistent Memory for Claude Code</p>
        </header>

        <div class="stats">
            <div class="stat-card">
                <div class="stat-number" id="total-facts">0</div>
                <div class="stat-label">Total Facts</div>
            </div>
            <div class="stat-card">
                <div class="stat-number" id="total-instances">0</div>
                <div class="stat-label">Running Instances</div>
            </div>
        </div>

        <div class="tabs">
            <div class="tab active" onclick="showTab('facts')">📚 Facts</div>
            <div class="tab" onclick="showTab('instances')">🖥️ Instances</div>
            <div class="tab" onclick="showTab('add')">➕ Add New</div>
        </div>

        <div id="facts-tab">
            <div class="search-box">
                <input type="text" id="search-input" placeholder="Search facts..." oninput="searchFacts()">
                <button class="btn-primary" onclick="searchFacts()">Search</button>
            </div>

            <div class="facts-grid" id="facts-list">
                <div class="empty-state">
                    <h3>Loading...</h3>
                </div>
            </div>
        </div>

        <div id="instances-tab" style="display: none;">
            <div class="instances-list" id="instances-list">
                <div class="empty-state">
                    <h3>Loading...</h3>
                </div>
            </div>
        </div>

        <div id="add-tab" style="display: none;">
            <div class="add-form">
                <h3>Add New Fact</h3>
                <div class="form-group">
                    <textarea id="new-fact-content" placeholder="Enter fact, decision, or context..."></textarea>
                </div>
                <div class="form-group">
                    <input type="text" id="new-fact-tags" placeholder="Tags (comma-separated, e.g., api, security)">
                </div>
                <button class="btn-primary" onclick="addFact()">Save Fact</button>
            </div>
        </div>
    </div>

    <script>
        let currentTab = 'facts';

        function showTab(tab) {
            document.querySelectorAll('.tab').forEach(t => t.classList.remove('active'));
            document.querySelectorAll('[id$="-tab"]').forEach(t => t.style.display = 'none');

            document.querySelector('.tabs').children[['facts', 'instances', 'add'].indexOf(tab)].classList.add('active');
            document.getElementById(tab + '-tab').style.display = 'block';

            currentTab = tab;

            if (tab === 'facts') loadFacts();
            if (tab === 'instances') loadInstances();
        }

        async function loadStatus() {
            try {
                const res = await fetch('/api/status');
                const data = await res.json();
                document.getElementById('total-facts').textContent = data.total_facts;
                document.getElementById('total-instances').textContent = data.instances;
            } catch (e) {
                console.error('Failed to load status:', e);
            }
        }

        async function loadFacts() {
            try {
                const query = document.getElementById('search-input').value;
                const res = await fetch('/api/facts?q=' + encodeURIComponent(query));
                const data = await res.json();

                const container = document.getElementById('facts-list');

                if (!data.facts || data.facts.length === 0) {
                    container.innerHTML = '<div class="empty-state"><h3>No facts found</h3><p>Add some facts to get started!</p></div>';
                    return;
                }

                container.innerHTML = data.facts.map(fact => {
                    const date = new Date(fact.created_at).toLocaleDateString('en-US', {
                        year: 'numeric', month: 'short', day: 'numeric'
                    });
                    const tags = fact.tags ? fact.tags.map(t => '<span class="tag">' + t + '</span>').join('') : '';
                    const source = fact.source_dir ? '<div class="fact-source">📁 ' + fact.source_dir + '</div>' : '';

                    return '<div class="fact-card">' +
                        '<div class="fact-header">' +
                        '<span class="fact-id">#' + fact.id + '</span>' +
                        '<span class="fact-date">' + date + '</span>' +
                        '</div>' +
                        '<div class="fact-content">' + escapeHtml(fact.content) + '</div>' +
                        '<div class="fact-tags">' + tags + '</div>' +
                        source +
                        '<button class="btn-danger" onclick="deleteFact(' + fact.id + ')" style="margin-top: 1rem;">Delete</button>' +
                        '</div>';
                }).join('');
            } catch (e) {
                console.error('Failed to load facts:', e);
            }
        }

        async function loadInstances() {
            try {
                const res = await fetch('/api/instances');
                const data = await res.json();

                const container = document.getElementById('instances-list');

                if (!data.instances || data.instances.length === 0) {
                    container.innerHTML = '<div class="empty-state"><h3>No running instances</h3><p>Start Claude Code with Claude PeePee to see instances here.</p></div>';
                    return;
                }

                container.innerHTML = data.instances.map(inst => {
                    const started = new Date(inst.started_at).toLocaleString();
                    const heartbeat = new Date(inst.last_heartbeat).toLocaleString();

                    return '<div class="instance-card">' +
                        '<div class="instance-id">🔗 ' + inst.id + '</div>' +
                        '<div class="instance-dir">📁 ' + inst.working_dir + '</div>' +
                        '<div class="instance-dir">PID: ' + inst.pid + '</div>' +
                        '<div class="instance-dir">Started: ' + started + '</div>' +
                        '<div class="instance-dir">Last heartbeat: ' + heartbeat + '</div>' +
                        '</div>';
                }).join('');
            } catch (e) {
                console.error('Failed to load instances:', e);
            }
        }

        async function addFact() {
            const content = document.getElementById('new-fact-content').value.trim();
            const tagsStr = document.getElementById('new-fact-tags').value.trim();

            if (!content) {
                alert('Please enter a fact');
                return;
            }

            const tags = tagsStr ? tagsStr.split(',').map(t => t.trim()).filter(t => t) : [];

            try {
                const res = await fetch('/api/facts', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ content, tags })
                });

                if (res.ok) {
                    document.getElementById('new-fact-content').value = '';
                    document.getElementById('new-fact-tags').value = '';
                    alert('Fact saved!');
                    loadStatus();
                    showTab('facts');
                } else {
                    alert('Failed to save fact');
                }
            } catch (e) {
                console.error('Failed to add fact:', e);
                alert('Failed to save fact');
            }
        }

        async function deleteFact(id) {
            if (!confirm('Are you sure you want to delete this fact?')) return;

            try {
                const res = await fetch('/api/facts/' + id, { method: 'DELETE' });
                if (res.ok) {
                    loadStatus();
                    loadFacts();
                } else {
                    alert('Failed to delete fact');
                }
            } catch (e) {
                console.error('Failed to delete fact:', e);
                alert('Failed to delete fact');
            }
        }

        function searchFacts() {
            loadFacts();
        }

        function escapeHtml(text) {
            const div = document.createElement('div');
            div.textContent = text;
            return div.innerHTML;
        }

        // Initial load
        loadStatus();
        loadFacts();

        // Refresh status every 30 seconds
        setInterval(loadStatus, 30000);
    </script>
</body>
</html>
`
