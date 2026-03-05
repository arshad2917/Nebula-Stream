package controlplane

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/nebula-stream/engine/internal/workflow"
)

type Server struct {
	registry *workflow.Registry
}

func NewServer(registry *workflow.Registry) *Server {
	return &Server{registry: registry}
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", s.handleHealth)
	mux.HandleFunc("/api/v1/workflows", s.handleWorkflows)
	return mux
}

func (s *Server) handleHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) handleWorkflows(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		s.handleWorkflowDeploy(w, r)
	case http.MethodGet:
		active, _ := s.registry.Active()
		writeJSON(w, http.StatusOK, map[string]any{
			"active":    active.Name,
			"workflows": s.registry.Names(),
		})
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleWorkflowDeploy(w http.ResponseWriter, r *http.Request) {
	raw, err := io.ReadAll(io.LimitReader(r.Body, 1<<20))
	if err != nil {
		writeErr(w, http.StatusBadRequest, fmt.Errorf("read request body: %w", err))
		return
	}

	def, err := workflow.ParseYAML(raw)
	if err != nil {
		writeErr(w, http.StatusBadRequest, fmt.Errorf("parse workflow yaml: %w", err))
		return
	}

	s.registry.Upsert(def)
	if r.URL.Query().Get("activate") != "false" {
		_ = s.registry.SetActive(def.Name)
	}

	writeJSON(w, http.StatusAccepted, map[string]string{
		"status":   "accepted",
		"workflow": def.Name,
	})
}

func writeErr(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, map[string]string{"error": err.Error()})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
