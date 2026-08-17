package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/core"
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/transport/nats"
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/otel/tracer"
)

type Role struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Permissions []string  `json:"permissions" db:"permissions"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type RoleCreatedEvent struct {
	Role      Role      `json:"role"`
	Timestamp time.Time `json:"timestamp"`
}

type RoleUpdatedEvent struct {
	Role      Role      `json:"role"`
	Timestamp time.Time `json:"timestamp"`
}

type RoleDeletedEvent struct {
	RoleID    string    `json:"role_id"`
	Timestamp time.Time `json:"timestamp"`
}

type RoleListRequest struct {
	Limit     int       `json:"limit"`
	Offset    int       `json:"offset"`
	Timestamp time.Time `json:"timestamp"`
}

type RoleGetRequest struct {
	RoleID    string    `json:"role_id"`
	Timestamp time.Time `json:"timestamp"`
}

type RoleHandler struct {
	publisher     *nats.Publisher
	subjectPrefix string
}

func NewRoleHandler(pub *nats.Publisher, subjectPrefix string) *RoleHandler {
	return &RoleHandler{
		publisher:     pub,
		subjectPrefix: subjectPrefix,
	}
}

func (s *RoleHandler) HandleCreate(ctx context.Context, msg *core.Message) error {
	ctx, span := tracer.StartConsumer(ctx, "Role.HandleCreate")
	defer span.End()
	ctx = core.InjectContext(ctx, msg.Headers)

	var event RoleCreatedEvent
	if err := json.Unmarshal(msg.Data, &event); err != nil {
		span.RecordError(err)
		return err
	}

	if event.Role.Name == "" {
		err := fmt.Errorf("invalid role data: missing required fields")
		span.RecordError(err)
		return err
	}

	span.SetAttributes(
		tracer.StringAttr("role.id", event.Role.ID),
		tracer.StringAttr("tenant.id", msg.Headers.Get(core.HeaderTenantID)),
	)

	outMsg := core.NewMessage(msg.Data)
	outMsg.Subject = s.subjectPrefix + ".role.db.create"
	outMsg.Headers = core.ExtractHeaders(ctx, outMsg.Headers)
	outMsg.Headers.Set("X-Business-Validated", "true")

	if err := s.publisher.Publish(ctx, outMsg.Subject, outMsg); err != nil {
		span.RecordError(err)
		return fmt.Errorf("publish error: %w", err)
	}
	return nil
}

func (s *RoleHandler) HandleGet(ctx context.Context, msg *core.Message) error {
	ctx, span := tracer.StartConsumer(ctx, "Role.HandleGet")
	defer span.End()
	ctx = core.InjectContext(ctx, msg.Headers)

	var event RoleGetRequest
	if err := json.Unmarshal(msg.Data, &event); err != nil {
		span.RecordError(err)
		return err
	}

	if event.RoleID == "" {
		err := fmt.Errorf("invalid request: missing role ID")
		span.RecordError(err)
		return err
	}

	span.SetAttributes(
		tracer.StringAttr("role.id", event.RoleID),
		tracer.StringAttr("tenant.id", msg.Headers.Get(core.HeaderTenantID)),
	)

	outMsg := core.NewMessage(msg.Data)
	outMsg.Subject = s.subjectPrefix + ".role.db.get"
	outMsg.Headers = core.ExtractHeaders(ctx, outMsg.Headers)
	outMsg.Headers.Set("X-Business-Validated", "true")

	if err := s.publisher.Publish(ctx, outMsg.Subject, outMsg); err != nil {
		span.RecordError(err)
		return fmt.Errorf("publish error: %w", err)
	}
	return nil
}

func (s *RoleHandler) HandleUpdate(ctx context.Context, msg *core.Message) error {
	ctx, span := tracer.StartConsumer(ctx, "Role.HandleUpdate")
	defer span.End()
	ctx = core.InjectContext(ctx, msg.Headers)

	var event RoleUpdatedEvent
	if err := json.Unmarshal(msg.Data, &event); err != nil {
		span.RecordError(err)
		return err
	}

	if event.Role.ID == "" {
		err := fmt.Errorf("invalid role data: missing ID")
		span.RecordError(err)
		return err
	}

	span.SetAttributes(
		tracer.StringAttr("role.id", event.Role.ID),
		tracer.StringAttr("tenant.id", msg.Headers.Get(core.HeaderTenantID)),
	)

	outMsg := core.NewMessage(msg.Data)
	outMsg.Subject = s.subjectPrefix + ".role.db.update"
	outMsg.Headers = core.ExtractHeaders(ctx, outMsg.Headers)
	outMsg.Headers.Set("X-Business-Validated", "true")

	if err := s.publisher.Publish(ctx, outMsg.Subject, outMsg); err != nil {
		span.RecordError(err)
		return fmt.Errorf("publish error: %w", err)
	}
	return nil
}

func (s *RoleHandler) HandleDelete(ctx context.Context, msg *core.Message) error {
	ctx, span := tracer.StartConsumer(ctx, "Role.HandleDelete")
	defer span.End()
	ctx = core.InjectContext(ctx, msg.Headers)

	var event RoleDeletedEvent
	if err := json.Unmarshal(msg.Data, &event); err != nil {
		span.RecordError(err)
		return err
	}

	if event.RoleID == "" {
		err := fmt.Errorf("invalid request: missing role ID")
		span.RecordError(err)
		return err
	}

	span.SetAttributes(
		tracer.StringAttr("role.id", event.RoleID),
		tracer.StringAttr("tenant.id", msg.Headers.Get(core.HeaderTenantID)),
	)

	outMsg := core.NewMessage(msg.Data)
	outMsg.Subject = s.subjectPrefix + ".role.db.delete"
	outMsg.Headers = core.ExtractHeaders(ctx, outMsg.Headers)
	outMsg.Headers.Set("X-Business-Validated", "true")

	if err := s.publisher.Publish(ctx, outMsg.Subject, outMsg); err != nil {
		span.RecordError(err)
		return fmt.Errorf("publish error: %w", err)
	}
	return nil
}

func (s *RoleHandler) HandleList(ctx context.Context, msg *core.Message) error {
	ctx, span := tracer.StartConsumer(ctx, "Role.HandleList")
	defer span.End()
	ctx = core.InjectContext(ctx, msg.Headers)

	var event RoleListRequest
	if err := json.Unmarshal(msg.Data, &event); err != nil {
		span.RecordError(err)
		return err
	}

	if event.Limit < 0 || event.Offset < 0 {
		err := fmt.Errorf("invalid pagination parameters")
		span.RecordError(err)
		return err
	}

	span.SetAttributes(
		tracer.StringAttr("tenant.id", msg.Headers.Get(core.HeaderTenantID)),
	)

	outMsg := core.NewMessage(msg.Data)
	outMsg.Subject = s.subjectPrefix + ".role.db.list"
	outMsg.Headers = core.ExtractHeaders(ctx, outMsg.Headers)
	outMsg.Headers.Set("X-Business-Validated", "true")

	if err := s.publisher.Publish(ctx, outMsg.Subject, outMsg); err != nil {
		span.RecordError(err)
		return fmt.Errorf("publish error: %w", err)
	}
	return nil
}
