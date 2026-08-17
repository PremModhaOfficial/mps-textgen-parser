package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"

	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/core/pool/workerpool"
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events"
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/core"
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/otel/logger"
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/otel/tracer"
)

type User struct {
	ID        string    `json:"id" db:"id"`
	Password  string    `json:"password" db:"password"`
	UserName  string    `json:"user_name" db:"user_name"`
	UserMail  string    `json:"user_mail" db:"user_mail"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
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

type UserGetRequest struct {
	UserID    string    `json:"user_id"`
	Timestamp time.Time `json:"timestamp"`
}

type UserListRequest struct {
	Limit     int       `json:"limit"`
	Offset    int       `json:"offset"`
	Timestamp time.Time `json:"timestamp"`
}

type UserHandler struct {
	publisher     *events.Publisher
	subjectPrefix string
}

func NewUserHandler(pub *events.Publisher, subjectPrefix string) *UserHandler {
	return &UserHandler{
		publisher:     pub,
		subjectPrefix: subjectPrefix,
	}
}

func (s *UserHandler) HandleCreate(req core.Request) {
	ctx := req.Context()
	ctx, span := tracer.StartConsumer(ctx, "User.HandleCreate")
	defer span.End()
	ctx = core.InjectContext(ctx, req.Headers())
	logger.Info(ctx, "request received", logger.String("handler", "User.Create"), logger.Int("bytes", len(req.Data())))

	var event UserCreatedEvent
	if err := json.Unmarshal(req.Data(), &event); err != nil {
		span.RecordError(err)
		_ = req.RespondError("400", "invalid JSON: "+err.Error(), nil)
		return
	}

	logger.Debug(ctx, "hook dispatched", logger.String("hook", "preCreateM"), logger.String("phase", "pre"),
		logger.String("mode", "async"))
	_ = workerpool.AsyncWithCtx(ctx, func(c context.Context) {
		s.preCreateM(c, span, &event)
	})

	logger.Debug(ctx, "hook executing", logger.String("hook", "preCreateD"), logger.String("phase", "pre"),
		logger.String("mode", "sync"))
	if err := s.preCreateD(ctx, span, &event); err != nil {
		logger.Error(ctx, "hook failed", logger.String("hook", "preCreateD"), logger.Err(err))
		span.RecordError(err)
		_ = req.RespondError("400", "pre-hook: "+err.Error(), nil)
		return
	}

	logger.Debug(ctx, "hook dispatched", logger.String("hook", "preCreateC"), logger.String("phase", "pre"),
		logger.String("mode", "async"))
	_ = workerpool.AsyncWithCtx(ctx, func(c context.Context) {
		s.preCreateC(c, span, &event)
	})

	logger.Debug(ctx, "hook executing", logger.String("hook", "preCreateB"), logger.String("phase", "pre"),
		logger.String("mode", "sync"))
	if err := s.preCreateB(ctx, span, &event); err != nil {
		logger.Error(ctx, "hook failed", logger.String("hook", "preCreateB"), logger.Err(err))
		span.RecordError(err)
		_ = req.RespondError("400", "pre-hook: "+err.Error(), nil)
		return
	}

	logger.Debug(ctx, "hook dispatched", logger.String("hook", "preCreateA"), logger.String("phase", "pre"),
		logger.String("mode", "async"))
	_ = workerpool.AsyncWithCtx(ctx, func(c context.Context) {
		s.preCreateA(c, span, &event)
	})

	if event.User.UserName == "" || event.User.UserMail == "" {
		err := fmt.Errorf("invalid user data: missing required fields")
		span.RecordError(err)
		_ = req.RespondError("400", err.Error(), nil)
		return
	}

	span.SetAttributes(
		tracer.StringAttr("user.id", event.User.ID),
		tracer.StringAttr("tenant.id", req.Header(core.HeaderTenantID)),
	)

	dalSubject := s.subjectPrefix + ".user.db.create"
	outMsg := &nats.Msg{Data: req.Data()}
	outMsg.Header = core.ExtractHeaders(ctx, nil)
	outMsg.Header.Set("X-Business-Validated", "true")

	dalCtx, dalCancel := context.WithTimeout(ctx, 10*time.Second)
	defer dalCancel()

	reply, err := s.publisher.Request(dalCtx, dalSubject, outMsg)
	if err != nil {
		span.RecordError(err)
		_ = req.RespondError("500", "DAL request error: "+err.Error(), nil)
		return
	}

	logger.Info(ctx, "DAL reply received", logger.String("handler", "User.create"), logger.Int("bytes", len(reply.Data)))
	responseData := reply.Data
	logger.Debug(ctx, "hook dispatched", logger.String("hook", "postCreateNotifyAdmin"), logger.String("phase", "post"),
		logger.String("mode", "async"))
	_ = workerpool.AsyncWithCtx(ctx, func(c context.Context) {
		s.postCreateNotifyAdmin(c, span, &event, responseData)
	})

	logger.Debug(ctx, "hook dispatched", logger.String("hook", "postCreateSendWelcome"), logger.String("phase", "post"),
		logger.String("mode", "async"))
	_ = workerpool.AsyncWithCtx(ctx, func(c context.Context) {
		s.postCreateSendWelcome(c, span, &event, responseData)
	})

	_ = req.Respond(responseData)
}

func (s *UserHandler) HandleUpdate(req core.Request) {
	ctx := req.Context()
	ctx, span := tracer.StartConsumer(ctx, "User.HandleUpdate")
	defer span.End()
	ctx = core.InjectContext(ctx, req.Headers())
	logger.Info(ctx, "request received", logger.String("handler", "User.Update"), logger.Int("bytes", len(req.Data())))

	var event UserUpdatedEvent
	if err := json.Unmarshal(req.Data(), &event); err != nil {
		span.RecordError(err)
		_ = req.RespondError("400", "invalid JSON: "+err.Error(), nil)
		return
	}

	logger.Debug(ctx, "hook dispatched", logger.String("hook", "preUpdateNotifyCache"), logger.String("phase", "pre"),
		logger.String("mode", "async"))
	_ = workerpool.AsyncWithCtx(ctx, func(c context.Context) {
		s.preUpdateNotifyCache(c, span, &event)
	})

	logger.Debug(ctx, "hook dispatched", logger.String("hook", "preUpdateAuditLog"), logger.String("phase", "pre"),
		logger.String("mode", "async"))
	_ = workerpool.AsyncWithCtx(ctx, func(c context.Context) {
		s.preUpdateAuditLog(c, span, &event)
	})

	if event.User.ID == "" {
		err := fmt.Errorf("invalid user data: missing ID")
		span.RecordError(err)
		_ = req.RespondError("400", err.Error(), nil)
		return
	}

	span.SetAttributes(
		tracer.StringAttr("user.id", event.User.ID),
		tracer.StringAttr("tenant.id", req.Header(core.HeaderTenantID)),
	)

	dalSubject := s.subjectPrefix + ".user.db.update"
	outMsg := &nats.Msg{Data: req.Data()}
	outMsg.Header = core.ExtractHeaders(ctx, nil)
	outMsg.Header.Set("X-Business-Validated", "true")

	dalCtx, dalCancel := context.WithTimeout(ctx, 10*time.Second)
	defer dalCancel()

	reply, err := s.publisher.Request(dalCtx, dalSubject, outMsg)
	if err != nil {
		span.RecordError(err)
		_ = req.RespondError("500", "DAL request error: "+err.Error(), nil)
		return
	}

	logger.Info(ctx, "DAL reply received", logger.String("handler", "User.update"), logger.Int("bytes", len(reply.Data)))
	responseData := reply.Data
	logger.Debug(ctx, "hook dispatched", logger.String("hook", "postUpdateInvalidateCache"), logger.String("phase", "post"),
		logger.String("mode", "async"))
	_ = workerpool.AsyncWithCtx(ctx, func(c context.Context) {
		s.postUpdateInvalidateCache(c, span, &event, responseData)
	})

	_ = req.Respond(responseData)
}

func (s *UserHandler) HandleDelete(req core.Request) {
	ctx := req.Context()
	ctx, span := tracer.StartConsumer(ctx, "User.HandleDelete")
	defer span.End()
	ctx = core.InjectContext(ctx, req.Headers())
	logger.Info(ctx, "request received", logger.String("handler", "User.Delete"), logger.Int("bytes", len(req.Data())))

	var event UserDeletedEvent
	if err := json.Unmarshal(req.Data(), &event); err != nil {
		span.RecordError(err)
		_ = req.RespondError("400", "invalid JSON: "+err.Error(), nil)
		return
	}

	logger.Debug(ctx, "hook executing", logger.String("hook", "preDeleteXyz"), logger.String("phase", "pre"),
		logger.String("mode", "sync"))
	if err := s.preDeleteXyz(ctx, span, &event); err != nil {
		logger.Error(ctx, "hook failed", logger.String("hook", "preDeleteXyz"), logger.Err(err))
		span.RecordError(err)
		_ = req.RespondError("400", "pre-hook: "+err.Error(), nil)
		return
	}

	logger.Debug(ctx, "hook executing", logger.String("hook", "preDeleteAbc"), logger.String("phase", "pre"),
		logger.String("mode", "sync"))
	if err := s.preDeleteAbc(ctx, span, &event); err != nil {
		logger.Error(ctx, "hook failed", logger.String("hook", "preDeleteAbc"), logger.Err(err))
		span.RecordError(err)
		_ = req.RespondError("400", "pre-hook: "+err.Error(), nil)
		return
	}

	logger.Debug(ctx, "hook executing", logger.String("hook", "preDeleteValidateTen"), logger.String("phase", "pre"),
		logger.String("mode", "sync"))
	if err := s.preDeleteValidateTen(ctx, span, &event); err != nil {
		logger.Error(ctx, "hook failed", logger.String("hook", "preDeleteValidateTen"), logger.Err(err))
		span.RecordError(err)
		_ = req.RespondError("400", "pre-hook: "+err.Error(), nil)
		return
	}

	if event.UserID == "" {
		err := fmt.Errorf("invalid request: missing user ID")
		span.RecordError(err)
		_ = req.RespondError("400", err.Error(), nil)
		return
	}

	span.SetAttributes(
		tracer.StringAttr("user.id", event.UserID),
		tracer.StringAttr("tenant.id", req.Header(core.HeaderTenantID)),
	)

	dalSubject := s.subjectPrefix + ".user.db.delete"
	outMsg := &nats.Msg{Data: req.Data()}
	outMsg.Header = core.ExtractHeaders(ctx, nil)
	outMsg.Header.Set("X-Business-Validated", "true")

	dalCtx, dalCancel := context.WithTimeout(ctx, 10*time.Second)
	defer dalCancel()

	reply, err := s.publisher.Request(dalCtx, dalSubject, outMsg)
	if err != nil {
		span.RecordError(err)
		_ = req.RespondError("500", "DAL request error: "+err.Error(), nil)
		return
	}

	logger.Info(ctx, "DAL reply received", logger.String("handler", "User.delete"), logger.Int("bytes", len(reply.Data)))
	responseData := reply.Data
	logger.Debug(ctx, "hook executing", logger.String("hook", "postDeleteMustDelete"), logger.String("phase", "post"),
		logger.String("mode", "sync"))
	responseData = s.postDeleteMustDelete(ctx, span, &event, responseData)
	_ = req.Respond(responseData)
}

func (s *UserHandler) HandleGet(req core.Request) {
	ctx := req.Context()
	ctx, span := tracer.StartConsumer(ctx, "User.HandleGet")
	defer span.End()
	ctx = core.InjectContext(ctx, req.Headers())
	logger.Info(ctx, "request received", logger.String("handler", "User.Get"), logger.Int("bytes", len(req.Data())))

	var event UserGetRequest
	if err := json.Unmarshal(req.Data(), &event); err != nil {
		span.RecordError(err)
		_ = req.RespondError("400", "invalid JSON: "+err.Error(), nil)
		return
	}

	logger.Debug(ctx, "hook dispatched", logger.String("hook", "preGetTraceAcces"), logger.String("phase", "pre"),
		logger.String("mode", "async"))
	_ = workerpool.AsyncWithCtx(ctx, func(c context.Context) {
		s.preGetTraceAcces(c, span, &event)
	})

	if event.UserID == "" {
		err := fmt.Errorf("invalid request: missing user ID")
		span.RecordError(err)
		_ = req.RespondError("400", err.Error(), nil)
		return
	}

	span.SetAttributes(
		tracer.StringAttr("user.id", event.UserID),
		tracer.StringAttr("tenant.id", req.Header(core.HeaderTenantID)),
	)

	dalSubject := s.subjectPrefix + ".user.db.get"
	outMsg := &nats.Msg{Data: req.Data()}
	outMsg.Header = core.ExtractHeaders(ctx, nil)
	outMsg.Header.Set("X-Business-Validated", "true")

	dalCtx, dalCancel := context.WithTimeout(ctx, 10*time.Second)
	defer dalCancel()

	reply, err := s.publisher.Request(dalCtx, dalSubject, outMsg)
	if err != nil {
		span.RecordError(err)
		_ = req.RespondError("500", "DAL request error: "+err.Error(), nil)
		return
	}

	logger.Info(ctx, "DAL reply received", logger.String("handler", "User.get"), logger.Int("bytes", len(reply.Data)))
	responseData := reply.Data
	logger.Debug(ctx, "hook executing", logger.String("hook", "postGetFilterPII"), logger.String("phase", "post"),
		logger.String("mode", "sync"))
	responseData = s.postGetFilterPII(ctx, span, &event, responseData)
	logger.Debug(ctx, "hook dispatched", logger.String("hook", "postGetCacheResult"), logger.String("phase", "post"),
		logger.String("mode", "async"))
	_ = workerpool.AsyncWithCtx(ctx, func(c context.Context) {
		s.postGetCacheResult(c, span, &event, responseData)
	})

	logger.Debug(ctx, "hook executing", logger.String("hook", "postGetEnrichData"), logger.String("phase", "post"),
		logger.String("mode", "sync"))
	responseData = s.postGetEnrichData(ctx, span, &event, responseData)
	_ = req.Respond(responseData)
}

func (s *UserHandler) HandleList(req core.Request) {
	ctx := req.Context()
	ctx, span := tracer.StartConsumer(ctx, "User.HandleList")
	defer span.End()
	ctx = core.InjectContext(ctx, req.Headers())
	logger.Info(ctx, "request received", logger.String("handler", "User.List"), logger.Int("bytes", len(req.Data())))

	var event UserListRequest
	if err := json.Unmarshal(req.Data(), &event); err != nil {
		span.RecordError(err)
		_ = req.RespondError("400", "invalid JSON: "+err.Error(), nil)
		return
	}

	if event.Limit < 0 || event.Offset < 0 {
		err := fmt.Errorf("invalid pagination parameters")
		span.RecordError(err)
		_ = req.RespondError("400", err.Error(), nil)
		return
	}

	span.SetAttributes(
		tracer.StringAttr("tenant.id", req.Header(core.HeaderTenantID)),
	)

	dalSubject := s.subjectPrefix + ".user.db.list"
	outMsg := &nats.Msg{Data: req.Data()}
	outMsg.Header = core.ExtractHeaders(ctx, nil)
	outMsg.Header.Set("X-Business-Validated", "true")

	dalCtx, dalCancel := context.WithTimeout(ctx, 10*time.Second)
	defer dalCancel()

	reply, err := s.publisher.Request(dalCtx, dalSubject, outMsg)
	if err != nil {
		span.RecordError(err)
		_ = req.RespondError("500", "DAL request error: "+err.Error(), nil)
		return
	}

	logger.Info(ctx, "DAL reply received", logger.String("handler", "User.list"), logger.Int("bytes", len(reply.Data)))
	_ = req.Respond(reply.Data)
}

type UserRolesAssignedEvent struct {
	UserID    string    `json:"user_id"`
	RolesID   string    `json:"roles_id"`
	Timestamp time.Time `json:"timestamp"`
}

type UserRolesListRequest struct {
	UserID    string    `json:"user_id"`
	Limit     int       `json:"limit"`
	Offset    int       `json:"offset"`
	Timestamp time.Time `json:"timestamp"`
}

type UserRolesRemovedEvent struct {
	UserID    string    `json:"user_id"`
	RolesID   string    `json:"roles_id"`
	Timestamp time.Time `json:"timestamp"`
}

type UserRolesHandler struct {
	publisher     *events.Publisher
	subjectPrefix string
}

func NewUserRolesHandler(pub *events.Publisher, subjectPrefix string) *UserRolesHandler {
	return &UserRolesHandler{
		publisher:     pub,
		subjectPrefix: subjectPrefix,
	}
}

func (s *UserRolesHandler) HandleAssign(req core.Request) {
	ctx := req.Context()
	ctx, span := tracer.StartConsumer(ctx, "UserRoles.HandleAssign")
	defer span.End()
	logger.Info(ctx, "request received", logger.String("handler", "UserRoles.Assign"), logger.Int("bytes", len(req.Data())))
	ctx = core.InjectContext(ctx, req.Headers())

	var event UserRolesAssignedEvent
	if err := json.Unmarshal(req.Data(), &event); err != nil {
		span.RecordError(err)
		_ = req.RespondError("400", "invalid JSON: "+err.Error(), nil)
		return
	}

	if event.UserID == "" || event.RolesID == "" {
		err := fmt.Errorf("invalid data: missing user or roles ID")
		span.RecordError(err)
		_ = req.RespondError("400", err.Error(), nil)
		return
	}

	span.SetAttributes(
		tracer.StringAttr("user.id", event.UserID),
		tracer.StringAttr("roles.id", event.RolesID),
		tracer.StringAttr("tenant.id", req.Header(core.HeaderTenantID)),
	)

	if err := s.preAssignHook(ctx, span, &event); err != nil {
		span.RecordError(err)
		_ = req.RespondError("400", "pre-hook: "+err.Error(), nil)
		return
	}

	dalSubject := s.subjectPrefix + ".user.roles.db.assign"
	outMsg := &nats.Msg{Data: req.Data()}
	outMsg.Header = core.ExtractHeaders(ctx, nil)
	outMsg.Header.Set("X-Business-Validated", "true")

	dalCtx, dalCancel := context.WithTimeout(ctx, 10*time.Second)
	defer dalCancel()

	reply, err := s.publisher.Request(dalCtx, dalSubject, outMsg)
	if err != nil {
		span.RecordError(err)
		_ = req.RespondError("500", "DAL request error: "+err.Error(), nil)
		return
	}

	logger.Info(ctx, "DAL reply received", logger.String("handler", "UserRoles.assign"), logger.Int("bytes", len(reply.Data)))

	responseData := s.postAssignHook(ctx, span, &event, reply.Data)
	_ = req.Respond(responseData)
}

func (s *UserRolesHandler) HandleList(req core.Request) {
	ctx := req.Context()
	ctx, span := tracer.StartConsumer(ctx, "UserRoles.HandleList")
	defer span.End()
	logger.Info(ctx, "request received", logger.String("handler", "UserRoles.List"), logger.Int("bytes", len(req.Data())))
	ctx = core.InjectContext(ctx, req.Headers())

	var event UserRolesListRequest
	if err := json.Unmarshal(req.Data(), &event); err != nil {
		span.RecordError(err)
		_ = req.RespondError("400", "invalid JSON: "+err.Error(), nil)
		return
	}

	if event.UserID == "" {
		err := fmt.Errorf("invalid request: missing user ID")
		span.RecordError(err)
		_ = req.RespondError("400", err.Error(), nil)
		return
	}

	span.SetAttributes(
		tracer.StringAttr("user.id", event.UserID),
		tracer.StringAttr("tenant.id", req.Header(core.HeaderTenantID)),
	)

	if err := s.preListHook(ctx, span, &event); err != nil {
		span.RecordError(err)
		_ = req.RespondError("400", "pre-hook: "+err.Error(), nil)
		return
	}

	dalSubject := s.subjectPrefix + ".user.roles.db.list"
	outMsg := &nats.Msg{Data: req.Data()}
	outMsg.Header = core.ExtractHeaders(ctx, nil)
	outMsg.Header.Set("X-Business-Validated", "true")

	dalCtx, dalCancel := context.WithTimeout(ctx, 10*time.Second)
	defer dalCancel()

	reply, err := s.publisher.Request(dalCtx, dalSubject, outMsg)
	if err != nil {
		span.RecordError(err)
		_ = req.RespondError("500", "DAL request error: "+err.Error(), nil)
		return
	}

	logger.Info(ctx, "DAL reply received", logger.String("handler", "UserRoles.list"), logger.Int("bytes", len(reply.Data)))

	responseData := s.postListHook(ctx, span, &event, reply.Data)
	_ = req.Respond(responseData)
}

func (s *UserRolesHandler) HandleRemove(req core.Request) {
	ctx := req.Context()
	ctx, span := tracer.StartConsumer(ctx, "UserRoles.HandleRemove")
	defer span.End()
	logger.Info(ctx, "request received", logger.String("handler", "UserRoles.Remove"), logger.Int("bytes", len(req.Data())))
	ctx = core.InjectContext(ctx, req.Headers())

	var event UserRolesRemovedEvent
	if err := json.Unmarshal(req.Data(), &event); err != nil {
		span.RecordError(err)
		_ = req.RespondError("400", "invalid JSON: "+err.Error(), nil)
		return
	}

	if event.UserID == "" || event.RolesID == "" {
		err := fmt.Errorf("invalid data: missing user or roles ID")
		span.RecordError(err)
		_ = req.RespondError("400", err.Error(), nil)
		return
	}

	span.SetAttributes(
		tracer.StringAttr("user.id", event.UserID),
		tracer.StringAttr("roles.id", event.RolesID),
		tracer.StringAttr("tenant.id", req.Header(core.HeaderTenantID)),
	)

	if err := s.preRemoveHook(ctx, span, &event); err != nil {
		span.RecordError(err)
		_ = req.RespondError("400", "pre-hook: "+err.Error(), nil)
		return
	}

	dalSubject := s.subjectPrefix + ".user.roles.db.remove"
	outMsg := &nats.Msg{Data: req.Data()}
	outMsg.Header = core.ExtractHeaders(ctx, nil)
	outMsg.Header.Set("X-Business-Validated", "true")

	dalCtx, dalCancel := context.WithTimeout(ctx, 10*time.Second)
	defer dalCancel()

	reply, err := s.publisher.Request(dalCtx, dalSubject, outMsg)
	if err != nil {
		span.RecordError(err)
		_ = req.RespondError("500", "DAL request error: "+err.Error(), nil)
		return
	}

	logger.Info(ctx, "DAL reply received", logger.String("handler", "UserRoles.remove"), logger.Int("bytes", len(reply.Data)))

	responseData := s.postRemoveHook(ctx, span, &event, reply.Data)
	_ = req.Respond(responseData)
}
