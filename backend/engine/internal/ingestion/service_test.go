package ingestion

import (
	"testing"

	"github.com/nebula-stream/engine/internal/bus"
)

func TestWorkflowNameFromEvent(t *testing.T) {
	fromMeta := workflowNameFromEvent(bus.EventEnvelope{Topic: "workflow.hello", Meta: map[string]string{"workflow": "priority"}})
	if fromMeta != "priority" {
		t.Fatalf("expected meta workflow, got %q", fromMeta)
	}

	fromTopic := workflowNameFromEvent(bus.EventEnvelope{Topic: "workflow.image.start"})
	if fromTopic != "image" {
		t.Fatalf("expected workflow from topic, got %q", fromTopic)
	}
}
