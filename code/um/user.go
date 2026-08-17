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

type User struct {
	ID        string    `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"-" db:"password"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type UserCreatedEvent struct {
	User      User      `json:"user"`
	Timestamp time.Time `json:"timestamp"`
}

type UserUpdatedEvent struct {
	User      User      `json:"user"`
	Timestamp time.Time `json:"timestamp"`
}

type UserDeletedEvent struct {
	UserID    string    `json:"user_id"`
	Timestamp time.Time `json:"timestamp"`
}

type UserListRequest struct {
	Limit     int       `json:"limit"`
	Offset    int       `json:"offset"`
	Timestamp time.Time `json:"timestamp"`
}

type UserGetRequest struct {
	UserID    string    `json:"user_id"`
	Timestamp time.Time `json:"timestamp"`
}

type UserRoleAssignedEvent struct {
	UserID    string    `json:"user_id"`
	RoleID    string    `json:"role_id"`
	Timestamp time.Time `json:"timestamp"`
}

type UserRoleRemovedEvent struct {
	UserID    string    `json:"user_id"`
	RoleID    string    `json:"role_id"`
	Timestamp time.Time `json:"timestamp"`
}

type UserRoleListRequest struct {
	UserID    string    `json:"user_id"`
	Limit     int       `json:"limit"`
	Offset    int       `json:"offset"`
	Timestamp time.Time `json:"timestamp"`
}

type UserHandler struct {
	publisher     *nats.Publisher
	subjectPrefix string
}

func NewUserHandler(pub *nats.Publisher, subjectPrefix string) *UserHandler {
	return &UserHandler{
		publisher:     pub,
		subjectPrefix: subjectPrefix,
	}
}

func (s *UserHandler) HandleCreate(ctx context.Context, msg *core.Message) error {
	ctx, span := tracer.StartConsumer(ctx, "User.HandleCreate")
	defer span.End()
	ctx = core.InjectContext(ctx, msg.Headers)

	var event UserCreatedEvent
	if err := json.Unmarshal(msg.Data, &event); err != nil {
		span.RecordError(err)
		return err
	}

	if event.User.Name == "" || event.User.Email == "" {
		err := fmt.Errorf("invalid user data: missing required fields")
		span.RecordError(err)
		return err
	}

	span.SetAttributes(
		tracer.StringAttr("user.id", event.User.ID),
		tracer.StringAttr("tenant.id", msg.Headers.Get(core.HeaderTenantID)),
	)

	outMsg := core.NewMessage(msg.Data)
	outMsg.Subject = s.subjectPrefix + ".user.db.create"
	outMsg.Headers = core.ExtractHeaders(ctx, outMsg.Headers)
	outMsg.Headers.Set("X-Business-Validated", "true")

	if err := s.publisher.Publish(ctx, outMsg.Subject, outMsg); err != nil {
		span.RecordError(err)
		return fmt.Errorf("publish error: %w", err)
	}
	return nil
}

func (s *UserHandler) HandleGet(ctx context.Context, msg *core.Message) error {
	ctx, span := tracer.StartConsumer(ctx, "User.HandleGet")
	defer span.End()
	ctx = core.InjectContext(ctx, msg.Headers)

	var event UserGetRequest
	if err := json.Unmarshal(msg.Data, &event); err != nil {
		span.RecordError(err)
		return err
	}

	if event.UserID == "" {
		err := fmt.Errorf("invalid request: missing user ID")
		span.RecordError(err)
		return err
	}

	span.SetAttributes(
		tracer.StringAttr("user.id", event.UserID),
		tracer.StringAttr("tenant.id", msg.Headers.Get(core.HeaderTenantID)),
	)

	outMsg := core.NewMessage(msg.Data)
	outMsg.Subject = s.subjectPrefix + ".user.db.get"
	outMsg.Headers = core.ExtractHeaders(ctx, outMsg.Headers)
	outMsg.Headers.Set("X-Business-Validated", "true")

	if err := s.publisher.Publish(ctx, outMsg.Subject, outMsg); err != nil {
		span.RecordError(err)
		return fmt.Errorf("publish error: %w", err)
	}
	return nil
}

func (s *UserHandler) HandleUpdate(ctx context.Context, msg *core.Message) error {
	ctx, span := tracer.StartConsumer(ctx, "User.HandleUpdate")
	defer span.End()
	ctx = core.InjectContext(ctx, msg.Headers)

	var event UserUpdatedEvent
	if err := json.Unmarshal(msg.Data, &event); err != nil {
		span.RecordError(err)
		return err
	}

	if event.User.ID == "" {
		err := fmt.Errorf("invalid user data: missing ID")
		span.RecordError(err)
		return err
	}

	span.SetAttributes(
		tracer.StringAttr("user.id", event.User.ID),
		tracer.StringAttr("tenant.id", msg.Headers.Get(core.HeaderTenantID)),
	)

	outMsg := core.NewMessage(msg.Data)
	outMsg.Subject = s.subjectPrefix + ".user.db.update"
	outMsg.Headers = core.ExtractHeaders(ctx, outMsg.Headers)
	outMsg.Headers.Set("X-Business-Validated", "true")

	if err := s.publisher.Publish(ctx, outMsg.Subject, outMsg); err != nil {
		span.RecordError(err)
		return fmt.Errorf("publish error: %w", err)
	}
	return nil
}

func (s *UserHandler) HandleDelete(ctx context.Context, msg *core.Message) error {
	ctx, span := tracer.StartConsumer(ctx, "User.HandleDelete")
	defer span.End()
	ctx = core.InjectContext(ctx, msg.Headers)

	var event UserDeletedEvent
	if err := json.Unmarshal(msg.Data, &event); err != nil {
		span.RecordError(err)
		return err
	}

	if event.UserID == "" {
		err := fmt.Errorf("invalid request: missing user ID")
		span.RecordError(err)
		return err
	}

	span.SetAttributes(
		tracer.StringAttr("user.id", event.UserID),
		tracer.StringAttr("tenant.id", msg.Headers.Get(core.HeaderTenantID)),
	)

	outMsg := core.NewMessage(msg.Data)
	outMsg.Subject = s.subjectPrefix + ".user.db.delete"
	outMsg.Headers = core.ExtractHeaders(ctx, outMsg.Headers)
	outMsg.Headers.Set("X-Business-Validated", "true")

	if err := s.publisher.Publish(ctx, outMsg.Subject, outMsg); err != nil {
		span.RecordError(err)
		return fmt.Errorf("publish error: %w", err)
	}
	return nil
}

func (s *UserHandler) HandleList(ctx context.Context, msg *core.Message) error {
	ctx, span := tracer.StartConsumer(ctx, "User.HandleList")
	defer span.End()
	ctx = core.InjectContext(ctx, msg.Headers)

	var event UserListRequest
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
	outMsg.Subject = s.subjectPrefix + ".user.db.list"
	outMsg.Headers = core.ExtractHeaders(ctx, outMsg.Headers)
	outMsg.Headers.Set("X-Business-Validated", "true")

	if err := s.publisher.Publish(ctx, outMsg.Subject, outMsg); err != nil {
		span.RecordError(err)
		return fmt.Errorf("publish error: %w", err)
	}
	return nil
}

type UserRoleHandler struct {
	publisher     *nats.Publisher
	subjectPrefix string
}

func NewUserRoleHandler(pub *nats.Publisher, subjectPrefix string) *UserRoleHandler {
	return &UserRoleHandler{
		publisher:     pub,
		subjectPrefix: subjectPrefix,
	}
}

func (s *UserRoleHandler) HandleAssign(ctx context.Context, msg *core.Message) error {
	ctx, span := tracer.StartConsumer(ctx, "UserRole.HandleAssign")
	defer span.End()
	ctx = core.InjectContext(ctx, msg.Headers)

	var event UserRoleAssignedEvent
	if err := json.Unmarshal(msg.Data, &event); err != nil {
		span.RecordError(err)
		return err
	}

	if event.UserID == "" || event.RoleID == "" {
		err := fmt.Errorf("invalid data: missing user or role ID")
		span.RecordError(err)
		return err
	}

	span.SetAttributes(
		tracer.StringAttr("user.id", event.UserID),
		tracer.StringAttr("role.id", event.RoleID),
		tracer.StringAttr("tenant.id", msg.Headers.Get(core.HeaderTenantID)),
	)

	outMsg := core.NewMessage(msg.Data)
	outMsg.Subject = s.subjectPrefix + ".user.role.db.assign"
	outMsg.Headers = core.ExtractHeaders(ctx, outMsg.Headers)
	outMsg.Headers.Set("X-Business-Validated", "true")

	if err := s.publisher.Publish(ctx, outMsg.Subject, outMsg); err != nil {
		span.RecordError(err)
		return fmt.Errorf("publish error: %w", err)
	}
	return nil
}

func (s *UserRoleHandler) HandleRemove(ctx context.Context, msg *core.Message) error {
	ctx, span := tracer.StartConsumer(ctx, "UserRole.HandleRemove")
	defer span.End()
	ctx = core.InjectContext(ctx, msg.Headers)

	var event UserRoleRemovedEvent
	if err := json.Unmarshal(msg.Data, &event); err != nil {
		span.RecordError(err)
		return err
	}

	if event.UserID == "" || event.RoleID == "" {
		err := fmt.Errorf("invalid data: missing user or role ID")
		span.RecordError(err)
		return err
	}

	span.SetAttributes(
		tracer.StringAttr("user.id", event.UserID),
		tracer.StringAttr("role.id", event.RoleID),
		tracer.StringAttr("tenant.id", msg.Headers.Get(core.HeaderTenantID)),
	)

	outMsg := core.NewMessage(msg.Data)
	outMsg.Subject = s.subjectPrefix + ".user.role.db.remove"
	outMsg.Headers = core.ExtractHeaders(ctx, outMsg.Headers)
	outMsg.Headers.Set("X-Business-Validated", "true")

	if err := s.publisher.Publish(ctx, outMsg.Subject, outMsg); err != nil {
		span.RecordError(err)
		return fmt.Errorf("publish error: %w", err)
	}
	return nil
}

func (s *UserRoleHandler) HandleList(ctx context.Context, msg *core.Message) error {
	ctx, span := tracer.StartConsumer(ctx, "UserRole.HandleList")
	defer span.End()
	ctx = core.InjectContext(ctx, msg.Headers)

	var event UserRoleListRequest
	if err := json.Unmarshal(msg.Data, &event); err != nil {
		span.RecordError(err)
		return err
	}

	if event.UserID == "" {
		err := fmt.Errorf("invalid request: missing user ID")
		span.RecordError(err)
		return err
	}

	span.SetAttributes(
		tracer.StringAttr("user.id", event.UserID),
		tracer.StringAttr("tenant.id", msg.Headers.Get(core.HeaderTenantID)),
	)

	outMsg := core.NewMessage(msg.Data)
	outMsg.Subject = s.subjectPrefix + ".user.role.db.list"
	outMsg.Headers = core.ExtractHeaders(ctx, outMsg.Headers)
	outMsg.Headers.Set("X-Business-Validated", "true")

	if err := s.publisher.Publish(ctx, outMsg.Subject, outMsg); err != nil {
		span.RecordError(err)
		return fmt.Errorf("publish error: %w", err)
	}
	return nil
}
