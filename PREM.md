# User Management MPS Template (NATS + OTEL)

This file contains the complete, single-file Go template for your JetBrains MPS DSL TextGen.
It implements user and role management by emitting NATS messages and tracing via OpenTelemetry, assuming the presence of a Data Access Layer (DAL) for data validation or queries.

## `usermanagement.go`

```go
package usermanagement

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/core"
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/transport/nats"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// ==========================================
// 1. Data Access Layer (Assumed Present)
// ==========================================

// DataAccessLayer represents the interface to your SQL database.
// Your MPS DSL will generate the actual SQL queries for this in another layer.
type DataAccessLayer interface {
	GetUserByID(ctx context.Context, id string) (*User, error)
	GetRoleByID(ctx context.Context, id string) (*Role, error)
}

// ==========================================
// 2. Domain Models (SQL Structures)
// ==========================================

type Role struct {
	ID          string   `json:"id" db:"id"`
	Name        string   `json:"name" db:"name"`
	Permissions []string `json:"permissions" db:"permissions"`
}

type User struct {
	ID       string   `json:"id" db:"id"`
	Username string   `json:"username" db:"username"`
	Email    string   `json:"email" db:"email"`
	RoleIDs  []string `json:"role_ids" db:"role_ids"`
}

// ==========================================
// 3. Event Payloads (NATS Messages)
// ==========================================

type UserCreatedEvent struct {
	User      User      `json:"user"`
	Timestamp time.Time `json:"timestamp"`
}

type RoleAssignedEvent struct {
	UserID    string    `json:"user_id"`
	RoleID    string    `json:"role_id"`
	Timestamp time.Time `json:"timestamp"`
}

// ==========================================
// 4. Service Implementation
// ==========================================

type UserService interface {
	CreateUser(ctx context.Context, user User) error
	AssignRole(ctx context.Context, userID, roleID string) error
}

// userService implements CQRS-style event emission. It validates via DAL
// but offloads the actual writes to NATS event subscribers.
type userService struct {
	publisher *nats.Publisher
	tracer    trace.Tracer
	dal       DataAccessLayer
}

func NewUserService(pub *nats.Publisher, tracer trace.Tracer, dal DataAccessLayer) UserService {
	return &userService{
		publisher: pub,
		tracer:    tracer,
		dal:       dal,
	}
}

func (s *userService) CreateUser(ctx context.Context, user User) error {
	// Start OTEL Span
	ctx, span := s.tracer.Start(ctx, "UserService.CreateUser")
	defer span.End()

	// Add semantic attributes for observability
	span.SetAttributes(
		attribute.String("user.id", user.ID),
		attribute.String("user.email", user.Email),
	)

	// Validate against DAL (Assumes DAL checks if email is already taken)
	// Example: _, err := s.dal.GetUserByEmail(ctx, user.Email)
	
	event := UserCreatedEvent{
		User:      user,
		Timestamp: time.Now().UTC(),
	}

	payload, err := json.Marshal(event)
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("failed to marshal UserCreatedEvent: %w", err)
	}

	msg := core.NewMessage(payload)
	msg.Subject = "motadata.iam.user.create"

	// Publish to NATS
	if err := s.publisher.Publish(ctx, msg.Subject, msg); err != nil {
		span.RecordError(err)
		return fmt.Errorf("failed to publish user creation event: %w", err)
	}

	return nil
}

func (s *userService) AssignRole(ctx context.Context, userID, roleID string) error {
	// Start OTEL Span
	ctx, span := s.tracer.Start(ctx, "UserService.AssignRole")
	defer span.End()

	span.SetAttributes(
		attribute.String("user.id", userID),
		attribute.String("role.id", roleID),
	)

	// Validate against DAL before publishing (Assumes DAL is present)
	_, err := s.dal.GetUserByID(ctx, userID)
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("user not found: %w", err)
	}

	_, err = s.dal.GetRoleByID(ctx, roleID)
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("role not found: %w", err)
	}

	event := RoleAssignedEvent{
		UserID:    userID,
		RoleID:    roleID,
		Timestamp: time.Now().UTC(),
	}

	payload, err := json.Marshal(event)
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("failed to marshal RoleAssignedEvent: %w", err)
	}

	msg := core.NewMessage(payload)
	msg.Subject = "motadata.iam.user.role.assign"

	// Publish to NATS
	if err := s.publisher.Publish(ctx, msg.Subject, msg); err != nil {
		span.RecordError(err)
		return fmt.Errorf("failed to publish role assignment event: %w", err)
	}

	return nil
}
```
