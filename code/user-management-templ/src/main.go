// motadata_user_management
package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/micro"

	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/config"
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events"
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/auth"
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/core"
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/otel"
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/otel/logger"
)

type microRequestAdapter struct {
	mr micro.Request
}

func (a *microRequestAdapter) Context() context.Context { return context.Background() }
func (a *microRequestAdapter) Subject() string          { return a.mr.Subject() }
func (a *microRequestAdapter) Reply() string            { return "" }
func (a *microRequestAdapter) Data() []byte             { return a.mr.Data() }
func (a *microRequestAdapter) Headers() nats.Header     { return nats.Header(a.mr.Headers()) }
func (a *microRequestAdapter) Header(key string) string { return a.mr.Headers().Get(key) }

func (a *microRequestAdapter) Respond(data []byte, opts ...core.RespondOption) error {
	return a.mr.Respond(data)
}
func (a *microRequestAdapter) RespondJSON(v any, opts ...core.RespondOption) error {
	return a.mr.RespondJSON(v)
}
func (a *microRequestAdapter) RespondError(code, desc string, data []byte, opts ...core.RespondOption) error {
	return a.mr.Error(code, desc, data)
}

func (a *microRequestAdapter) Ack() error                         { return nil }
func (a *microRequestAdapter) Nak() error                         { return nil }
func (a *microRequestAdapter) NakWithDelay(_ time.Duration) error { return nil }
func (a *microRequestAdapter) Term() error                        { return nil }
func (a *microRequestAdapter) InProgress() error                  { return nil }

func adaptRequest(mr micro.Request) core.Request {
	return &microRequestAdapter{mr: mr}
}

const prefix = "motadata"

func main() {
	otel.InitFromEnv()
	logcfg := logger.DefaultConfig()
	logcfg.ServiceName = "motadata_user_management"
	logcfg.ServiceVersion = "0.1.0"
	logcfg.Environment = "development"
	logcfg.Level = "debug"
	logInst := logger.MustInit(logcfg)
	defer logInst.Close()

	ctx0 := context.Background()
	logger.Info(ctx0, "Starting motadata_user_management")

	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = "DEFAULT_NATS_URL"
	}

	cfg := config.DefaultEventsConfig()
	cfg.Servers = []string{natsURL}

	creds := &auth.Credentials{Type: auth.TypeNone}
	conn := events.NewConnection("", cfg, creds)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := conn.Connect(ctx); err != nil {
		logger.Fatal(ctx, "failed to connect to NATS", logger.Err(err))

	}
	defer conn.Close(context.Background())

	publisher := events.NewPublisher(conn)

	rolesHandler := NewRolesHandler(publisher, prefix)
	permissionsHandler := NewPermissionsHandler(publisher, prefix)
	userHandler := NewUserHandler(publisher, prefix)
	rolesPermissionsHandler := NewRolesPermissionsHandler(publisher, prefix)
	userRolesHandler := NewUserRolesHandler(publisher, prefix)

	srv, err := micro.AddService(conn.Conn(), micro.Config{
		Name:    "motadata_user_management",
		Version: "0.1.0",
	})
	if err != nil {
		logger.Fatal(ctx, "failed to connect to NATS", logger.Err(err))

	}

	root := srv.AddGroup(prefix)

	// Roles endpoints
	rolesGroup := root.AddGroup("roles")
	rolesGroup.AddEndpoint("create", micro.HandlerFunc(func(mr micro.Request) { rolesHandler.HandleCreate(adaptRequest(mr)) }))
	rolesGroup.AddEndpoint("update", micro.HandlerFunc(func(mr micro.Request) { rolesHandler.HandleUpdate(adaptRequest(mr)) }))
	rolesGroup.AddEndpoint("list", micro.HandlerFunc(func(mr micro.Request) { rolesHandler.HandleList(adaptRequest(mr)) }))
	rolesGroup.AddEndpoint("delete", micro.HandlerFunc(func(mr micro.Request) { rolesHandler.HandleDelete(adaptRequest(mr)) }))
	rolesGroup.AddEndpoint("get", micro.HandlerFunc(func(mr micro.Request) { rolesHandler.HandleGet(adaptRequest(mr)) }))

	// Permissions endpoints
	permissionsGroup := root.AddGroup("permissions")
	permissionsGroup.AddEndpoint("create", micro.HandlerFunc(func(mr micro.Request) { permissionsHandler.HandleCreate(adaptRequest(mr)) }))
	permissionsGroup.AddEndpoint("update", micro.HandlerFunc(func(mr micro.Request) { permissionsHandler.HandleUpdate(adaptRequest(mr)) }))
	permissionsGroup.AddEndpoint("delete", micro.HandlerFunc(func(mr micro.Request) { permissionsHandler.HandleDelete(adaptRequest(mr)) }))
	permissionsGroup.AddEndpoint("list", micro.HandlerFunc(func(mr micro.Request) { permissionsHandler.HandleList(adaptRequest(mr)) }))
	permissionsGroup.AddEndpoint("get", micro.HandlerFunc(func(mr micro.Request) { permissionsHandler.HandleGet(adaptRequest(mr)) }))

	// User endpoints
	userGroup := root.AddGroup("user")
	userGroup.AddEndpoint("create", micro.HandlerFunc(func(mr micro.Request) { userHandler.HandleCreate(adaptRequest(mr)) }))
	userGroup.AddEndpoint("update", micro.HandlerFunc(func(mr micro.Request) { userHandler.HandleUpdate(adaptRequest(mr)) }))
	userGroup.AddEndpoint("delete", micro.HandlerFunc(func(mr micro.Request) { userHandler.HandleDelete(adaptRequest(mr)) }))
	userGroup.AddEndpoint("get", micro.HandlerFunc(func(mr micro.Request) { userHandler.HandleGet(adaptRequest(mr)) }))
	userGroup.AddEndpoint("list", micro.HandlerFunc(func(mr micro.Request) { userHandler.HandleList(adaptRequest(mr)) }))

	// Roles -> Permissions endpoints
	rolesPermissionsGroup := root.AddGroup("roles").AddGroup("permissions")
	rolesPermissionsGroup.AddEndpoint("assign", micro.HandlerFunc(func(mr micro.Request) { rolesPermissionsHandler.HandleAssign(adaptRequest(mr)) }))
	rolesPermissionsGroup.AddEndpoint("list", micro.HandlerFunc(func(mr micro.Request) { rolesPermissionsHandler.HandleList(adaptRequest(mr)) }))
	rolesPermissionsGroup.AddEndpoint("remove", micro.HandlerFunc(func(mr micro.Request) { rolesPermissionsHandler.HandleRemove(adaptRequest(mr)) }))

	// User -> Roles endpoints
	userRolesGroup := root.AddGroup("user").AddGroup("roles")
	userRolesGroup.AddEndpoint("assign", micro.HandlerFunc(func(mr micro.Request) { userRolesHandler.HandleAssign(adaptRequest(mr)) }))
	userRolesGroup.AddEndpoint("list", micro.HandlerFunc(func(mr micro.Request) { userRolesHandler.HandleList(adaptRequest(mr)) }))
	userRolesGroup.AddEndpoint("remove", micro.HandlerFunc(func(mr micro.Request) { userRolesHandler.HandleRemove(adaptRequest(mr)) }))

	logger.Info(ctx0, "motadata_user_management listening on all subjects")

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	logger.Info(ctx0, "motadata_user_management shutting down")
	if err := srv.Stop(); err != nil {
		logger.Error(ctx0, "error stopping service:", logger.Err(err))
	}
	logger.Info(ctx0, "motadata_user_management stopped.")
}
