package ingestion

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/nebula-stream/engine/internal/bus"
	"github.com/nebula-stream/engine/internal/engine"
	"github.com/nebula-stream/engine/internal/workflow"
)

type Service struct {
	busClient *bus.Client
	registry  *workflow.Registry
	runner    *engine.Runner
}

func NewService(busClient *bus.Client, registry *workflow.Registry) *Service {
	return &Service{
		busClient: busClient,
		registry:  registry,
		runner:    engine.NewRunner(),
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

	log.Printf("workflow executed name=%s steps=%d", def.Name, len(results))

	return nil
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
