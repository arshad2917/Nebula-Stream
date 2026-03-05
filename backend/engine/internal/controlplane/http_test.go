package controlplane

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/nebula-stream/engine/internal/workflow"
)

func TestWorkflowDeployEndpoint(t *testing.T) {
	registry := workflow.NewRegistry(workflow.Definition{Name: "hello", Version: "v1", Triggers: []workflow.Trigger{{Type: "manual"}}, Steps: []workflow.Step{{ID: "s1", Type: "builtin.log"}}})
	srv := NewServer(registry)

	body := `version: v1
name: uploaded
triggers:
  - type: manual
steps:
  - id: s1
    type: builtin.log
`

	req := httptest.NewRequest(http.MethodPost, "/api/v1/workflows", strings.NewReader(body))
	res := httptest.NewRecorder()

	srv.Handler().ServeHTTP(res, req)

	if res.Code != http.StatusAccepted {
		t.Fatalf("unexpected status code: %d", res.Code)
	}

	active, ok := registry.Active()
	if !ok || active.Name != "uploaded" {
		t.Fatalf("expected uploaded workflow to be active, got %+v", active)
	}
}
