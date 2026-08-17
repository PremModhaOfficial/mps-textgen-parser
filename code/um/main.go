package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/auth"
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/config"
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/core"
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/transport/nats"
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/otel"
)

func main() {
	log.Println("Starting UserManagement (Business Logic Layer)...")

	otel.InitFromEnv()

	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = "nats://localhost:4222"
	}

	cfg := config.DefaultConfig()
	cfg.Servers = []string{natsURL}

	creds := &auth.Credentials{Type: auth.TypeNone}
	conn := nats.NewConnection("um-client", cfg, creds)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := conn.Connect(ctx); err != nil {
		log.Fatalf("failed to connect to NATS: %v", err)
	}
	defer conn.Close(context.Background())

	publisher := nats.NewPublisher(conn)
	prefix := "um"

	userHandler := NewUserHandler(publisher, prefix)
	roleHandler := NewRoleHandler(publisher, prefix)
	userRoleHandler := NewUserRoleHandler(publisher, prefix)

	subscriber := nats.NewSubscriber(conn)

	subjects := map[string]core.MessageHandler{
		// User handlers
		"um.user.create": userHandler.HandleCreate,
		"um.user.get":    userHandler.HandleGet,
		"um.user.update": userHandler.HandleUpdate,
		"um.user.delete": userHandler.HandleDelete,
		"um.user.list":   userHandler.HandleList,
		// Role handlers
		"um.role.create": roleHandler.HandleCreate,
		"um.role.get":    roleHandler.HandleGet,
		"um.role.update": roleHandler.HandleUpdate,
		"um.role.delete": roleHandler.HandleDelete,
		"um.role.list":   roleHandler.HandleList,
		// User -> Role handlers
		"um.user.role.assign": userRoleHandler.HandleAssign,
		"um.user.role.remove": userRoleHandler.HandleRemove,
		"um.user.role.list":   userRoleHandler.HandleList,
	}

	for subject, handler := range subjects {
		if _, err := subscriber.Subscribe(context.Background(), subject, handler); err != nil {
			log.Fatalf("failed to subscribe to %s: %v", subject, err)
		}
	}

	log.Println("UserManagement listening on all subjects")

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	log.Println("UserManagement shutting down.")
}
