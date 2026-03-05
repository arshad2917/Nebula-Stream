package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestRunDeployValidation(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.yaml")
	raw := []byte("version: v1\nname: demo\ntriggers:\n  - type: manual\nsteps:\n  - id: s1\n    type: builtin.log\n")
	if err := os.WriteFile(path, raw, 0o600); err != nil {
		t.Fatalf("write workflow: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		w.WriteHeader(http.StatusAccepted)
	}))
	defer server.Close()

	if err := run([]string{"deploy", "-f", path, "--engine-url", server.URL}); err != nil {
		t.Fatalf("deploy should succeed: %v", err)
	}
}

func TestRunUnknownCommand(t *testing.T) {
	if err := run([]string{"unknown"}); err == nil {
		t.Fatal("expected unknown command error")
	}
}
