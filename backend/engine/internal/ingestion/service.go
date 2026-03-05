package ingestion

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/nebula-stream/engine/internal/bus"
	"github.com/nebula-stream/engine/internal/engine"
	"github.com/nebula-stream/engine/internal/state"
	"github.com/nebula-stream/engine/internal/workflow"
)

type Service struct {
	busClient *bus.Client
	registry  *workflow.Registry
	runner    *engine.Runner
	store     state.Store
}

func NewService(busClient *bus.Client, registry *workflow.Registry, store state.Store) *Service {
	return &Service{
		busClient: busClient,
		registry:  registry,
		runner:    engine.NewRunner(),
		store:     store,
	}
}

func (s *Service) Start(ctx context.Context, subject string) error {
	if s == nil || s.busClient == nil {
		return fmt.Errorf("ingestion service requires an initialized bus client")
	}

	if _, err := s.busClient.Subscribe(subject, func(event bus.EventEnvelope) error {
		return s.handle(event)
	}); err != nil {
		return err
	}

	active, ok := s.registry.Active()
	if ok {
		log.Printf("ingestion subscribed subject=%s workflow=%s", subject, active.Name)
	} else {
		log.Printf("ingestion subscribed subject=%s workflow=<none>", subject)
	}

	<-ctx.Done()
	return nil
}

func (s *Service) handle(event bus.EventEnvelope) error {
	log.Printf("event received id=%s topic=%s payload=%dB", event.ID, event.Topic, len(event.Payload))

	def, err := s.resolveWorkflow(event)
	if err != nil {
		return err
	}

	results, err := s.runner.Execute(context.Background(), def, event)
	if err != nil {
		return err
	}

	if err := s.persistExecution(def.Name, event, results); err != nil {
		return err
	}

	log.Printf("workflow executed name=%s steps=%d", def.Name, len(results))

	return nil
}

func (s *Service) persistExecution(workflowName string, event bus.EventEnvelope, results map[string]engine.StepResult) error {
	if s.store == nil {
		return nil
	}

	record := map[string]any{
		"event_id":    event.ID,
		"workflow":    workflowName,
		"topic":       event.Topic,
		"executed_at": time.Now().UTC(),
		"results":     results,
	}

	raw, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("marshal execution record: %w", err)
	}

	if err := s.store.Save(executionKey(event.ID), raw); err != nil {
		return fmt.Errorf("save execution record: %w", err)
	}

	if err := s.store.Save(latestExecutionKey(workflowName), raw); err != nil {
		return fmt.Errorf("save latest execution record: %w", err)
	}

	return nil
}

func executionKey(eventID string) string {
	return fmt.Sprintf("execution:%s", eventID)
}

func latestExecutionKey(workflowName string) string {
	return fmt.Sprintf("workflow:%s:latest", workflowName)
}

func (s *Service) resolveWorkflow(event bus.EventEnvelope) (workflow.Definition, error) {
	if s.registry == nil {
		return workflow.Definition{}, fmt.Errorf("workflow registry is not initialized")
	}

	if name := workflowNameFromEvent(event); name != "" {
		if def, ok := s.registry.Get(name); ok {
			return def, nil
		}
		return workflow.Definition{}, fmt.Errorf("workflow %q not found", name)
	}

	def, ok := s.registry.Active()
	if !ok {
		return workflow.Definition{}, fmt.Errorf("no active workflow configured")
	}

	return def, nil
}

func workflowNameFromEvent(event bus.EventEnvelope) string {
	if name := event.Meta["workflow"]; name != "" {
		return name
	}

	parts := strings.Split(event.Topic, ".")
	if len(parts) >= 2 && parts[0] == "workflow" {
		return parts[1]
	}

	return ""
}
