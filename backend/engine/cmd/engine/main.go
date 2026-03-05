package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/nebula-stream/engine/internal/bus"
	"github.com/nebula-stream/engine/internal/config"
	"github.com/nebula-stream/engine/internal/controlplane"
	"github.com/nebula-stream/engine/internal/ingestion"
	"github.com/nebula-stream/engine/internal/workflow"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}
	log.Printf("nebula-engine bootstrap started node=%s nats=%s heartbeat=%ds", cfg.NodeID, cfg.NATSURL, cfg.HeartbeatSecs)

	def, err := workflow.ParseFile(cfg.WorkflowPath)
	if err != nil {
		log.Fatalf("parse workflow file %q: %v", cfg.WorkflowPath, err)
	}

	registry := workflow.NewRegistry(def)

	busClient, err := bus.Connect(cfg.NATSURL)
	if err != nil {
		log.Fatalf("connect bus: %v", err)
	}
	defer busClient.Close()

	api := controlplane.NewServer(registry)
	httpServer := &http.Server{
		Addr:    cfg.APIAddr,
		Handler: api.Handler(),
	}

	go func() {
		log.Printf("control plane listening addr=%s", cfg.APIAddr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("control plane server failed: %v", err)
		}
	}()

	svc := ingestion.NewService(busClient, registry)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := svc.Start(ctx, cfg.IngestSubject); err != nil {
		log.Fatalf("start ingestion: %v", err)
	}

	_ = httpServer.Shutdown(context.Background())

	log.Println("engine shutdown complete")
}
