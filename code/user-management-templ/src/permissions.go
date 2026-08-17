package main

import (
 "context"
 "encoding/json"
 "fmt"
 "time"

 "github.com/nats-io/nats.go"

 "dev.azure.com/Motadata/NextGen/motadata-go-sdk/events"
 "dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/core"
 "dev.azure.com/Motadata/NextGen/motadata-go-sdk/otel/tracer"
 "dev.azure.com/Motadata/NextGen/motadata-go-sdk/otel/logger"
"dev.azure.com/Motadata/NextGen/motadata-go-sdk/core/pool/workerpool"
)

type Permissions struct {
 ID string `json:"id" db:"id"`
 Name string `json:"name" db:"name"`
 Moduele string `json:"moduele" db:"moduele"`
 Can string `json:"can" db:"can"`
}

type PermissionsCreatedEvent struct {
 Permissions Permissions `json:"permissions"`
 Timestamp time.Time `json:"timestamp"`
}

type PermissionsUpdatedEvent struct {
 Permissions Permissions `json:"permissions"`
 Timestamp time.Time `json:"timestamp"`
}

type PermissionsDeletedEvent struct {
 PermissionsID string `json:"permissions_id"`
 Timestamp time.Time `json:"timestamp"`
}

type PermissionsListRequest struct {
 Limit     int       `json:"limit"`
 Offset    int       `json:"offset"`
 Timestamp time.Time `json:"timestamp"`
}

type PermissionsGetRequest struct {
 PermissionsID string `json:"permissions_id"`
 Timestamp time.Time `json:"timestamp"`
}

type PermissionsHandler struct {
 publisher     *events.Publisher
 subjectPrefix string
}

func NewPermissionsHandler(pub *events.Publisher, subjectPrefix string) *PermissionsHandler {
 return &PermissionsHandler{
  publisher:     pub,
  subjectPrefix: subjectPrefix,
 }
}

func (s *PermissionsHandler) HandleCreate(req core.Request) {
 ctx := req.Context()
 ctx, span := tracer.StartConsumer(ctx, "Permissions.HandleCreate")
 defer span.End()
 ctx = core.InjectContext(ctx, req.Headers())
logger.Info(ctx, "request received", logger.String("handler", "Permissions.Create"), logger.Int("bytes", len(req.Data())))

 var event PermissionsCreatedEvent
 if err := json.Unmarshal(req.Data(), &event); err != nil {
  span.RecordError(err)
  _ = req.RespondError("400", "invalid JSON: " + err.Error(), nil)
  return
 }

 logger.Debug(ctx, "hook dispatched", logger.String("hook", "preCreateNotifyAdmin"), logger.String("phase", "pre"),
  logger.String("mode", "async"))
_ = workerpool.AsyncWithCtx(ctx, func(c context.Context) {
   s.preCreateNotifyAdmin(c, span, &event)
 })

 logger.Debug(ctx, "hook executing", logger.String("hook", "preCreateNotifyManager"), logger.String("phase", "pre"),
  logger.String("mode", "sync"))
 if err := s.preCreateNotifyManager(ctx, span, &event); err != nil {
  logger.Error(ctx, "hook failed", logger.String("hook", "preCreateNotifyManager"), logger.Err(err))
  span.RecordError(err)
  _ = req.RespondError("400", "pre-hook: " + err.Error(), nil)
  return
 }

 if event.Permissions.Name == "" || event.Permissions.Moduele == "" || event.Permissions.Can == "" {
  err := fmt.Errorf("invalid permissions data: missing required fields")
  span.RecordError(err)
  _ = req.RespondError("400", err.Error(), nil)
  return
 }

 span.SetAttributes(
  tracer.StringAttr("permissions.id", event.Permissions.ID),
  tracer.StringAttr("tenant.id", req.Header(core.HeaderTenantID)),
 )

 dalSubject := s.subjectPrefix + ".permissions.db.create"
 outMsg := &nats.Msg{Data: req.Data()}
 outMsg.Header = core.ExtractHeaders(ctx, nil)
 outMsg.Header.Set("X-Business-Validated", "true")

 dalCtx, dalCancel := context.WithTimeout(ctx, 10*time.Second)
 defer dalCancel()

 reply, err := s.publisher.Request(dalCtx, dalSubject, outMsg)
 if err != nil {
  span.RecordError(err)
  _ = req.RespondError("500", "DAL request error: " + err.Error(), nil)
  return
 }

logger.Info(ctx, "DAL reply received", logger.String("handler", "Permissions.create"), logger.Int("bytes", len(reply.Data)))
 responseData := reply.Data
 logger.Debug(ctx, "hook dispatched", logger.String("hook", "postCreateAudit"), logger.String("phase", "post"),
  logger.String("mode", "async"))
 _ = workerpool.AsyncWithCtx(ctx, func(c context.Context) {
  s.postCreateAudit(c, span, &event, responseData)
 })

 _ = req.Respond(responseData)
}

func (s *PermissionsHandler) HandleUpdate(req core.Request) {
 ctx := req.Context()
 ctx, span := tracer.StartConsumer(ctx, "Permissions.HandleUpdate")
 defer span.End()
 ctx = core.InjectContext(ctx, req.Headers())
logger.Info(ctx, "request received", logger.String("handler", "Permissions.Update"), logger.Int("bytes", len(req.Data())))

 var event PermissionsUpdatedEvent
 if err := json.Unmarshal(req.Data(), &event); err != nil {
  span.RecordError(err)
  _ = req.RespondError("400", "invalid JSON: " + err.Error(), nil)
  return
 }

 if event.Permissions.ID == "" {
  err := fmt.Errorf("invalid permissions data: missing ID")
  span.RecordError(err)
  _ = req.RespondError("400", err.Error(), nil)
  return
 }

 span.SetAttributes(
  tracer.StringAttr("permissions.id", event.Permissions.ID),
  tracer.StringAttr("tenant.id", req.Header(core.HeaderTenantID)),
 )

 dalSubject := s.subjectPrefix + ".permissions.db.update"
 outMsg := &nats.Msg{Data: req.Data()}
 outMsg.Header = core.ExtractHeaders(ctx, nil)
 outMsg.Header.Set("X-Business-Validated", "true")

 dalCtx, dalCancel := context.WithTimeout(ctx, 10*time.Second)
 defer dalCancel()

 reply, err := s.publisher.Request(dalCtx, dalSubject, outMsg)
 if err != nil {
  span.RecordError(err)
  _ = req.RespondError("500", "DAL request error: " + err.Error(), nil)
  return
 }

logger.Info(ctx, "DAL reply received", logger.String("handler", "Permissions.update"), logger.Int("bytes", len(reply.Data)))
 _ = req.Respond(reply.Data)
}

func (s *PermissionsHandler) HandleDelete(req core.Request) {
 ctx := req.Context()
 ctx, span := tracer.StartConsumer(ctx, "Permissions.HandleDelete")
 defer span.End()
 ctx = core.InjectContext(ctx, req.Headers())
logger.Info(ctx, "request received", logger.String("handler", "Permissions.Delete"), logger.Int("bytes", len(req.Data())))

 var event PermissionsDeletedEvent
 if err := json.Unmarshal(req.Data(), &event); err != nil {
  span.RecordError(err)
  _ = req.RespondError("400", "invalid JSON: " + err.Error(), nil)
  return
 }

 if event.PermissionsID == "" {
  err := fmt.Errorf("invalid request: missing permissions ID")
  span.RecordError(err)
  _ = req.RespondError("400", err.Error(), nil)
  return
 }

 span.SetAttributes(
  tracer.StringAttr("permissions.id", event.PermissionsID),
  tracer.StringAttr("tenant.id", req.Header(core.HeaderTenantID)),
 )

 dalSubject := s.subjectPrefix + ".permissions.db.delete"
 outMsg := &nats.Msg{Data: req.Data()}
 outMsg.Header = core.ExtractHeaders(ctx, nil)
 outMsg.Header.Set("X-Business-Validated", "true")

 dalCtx, dalCancel := context.WithTimeout(ctx, 10*time.Second)
 defer dalCancel()

 reply, err := s.publisher.Request(dalCtx, dalSubject, outMsg)
 if err != nil {
  span.RecordError(err)
  _ = req.RespondError("500", "DAL request error: " + err.Error(), nil)
  return
 }

logger.Info(ctx, "DAL reply received", logger.String("handler", "Permissions.delete"), logger.Int("bytes", len(reply.Data)))
 _ = req.Respond(reply.Data)
}

func (s *PermissionsHandler) HandleList(req core.Request) {
 ctx := req.Context()
 ctx, span := tracer.StartConsumer(ctx, "Permissions.HandleList")
 defer span.End()
 ctx = core.InjectContext(ctx, req.Headers())
logger.Info(ctx, "request received", logger.String("handler", "Permissions.List"), logger.Int("bytes", len(req.Data())))

 var event PermissionsListRequest
 if err := json.Unmarshal(req.Data(), &event); err != nil {
  span.RecordError(err)
  _ = req.RespondError("400", "invalid JSON: " + err.Error(), nil)
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

 dalSubject := s.subjectPrefix + ".permissions.db.list"
 outMsg := &nats.Msg{Data: req.Data()}
 outMsg.Header = core.ExtractHeaders(ctx, nil)
 outMsg.Header.Set("X-Business-Validated", "true")

 dalCtx, dalCancel := context.WithTimeout(ctx, 10*time.Second)
 defer dalCancel()

 reply, err := s.publisher.Request(dalCtx, dalSubject, outMsg)
 if err != nil {
  span.RecordError(err)
  _ = req.RespondError("500", "DAL request error: " + err.Error(), nil)
  return
 }

logger.Info(ctx, "DAL reply received", logger.String("handler", "Permissions.list"), logger.Int("bytes", len(reply.Data)))
 _ = req.Respond(reply.Data)
}

func (s *PermissionsHandler) HandleGet(req core.Request) {
 ctx := req.Context()
 ctx, span := tracer.StartConsumer(ctx, "Permissions.HandleGet")
 defer span.End()
 ctx = core.InjectContext(ctx, req.Headers())
logger.Info(ctx, "request received", logger.String("handler", "Permissions.Get"), logger.Int("bytes", len(req.Data())))

 var event PermissionsGetRequest
 if err := json.Unmarshal(req.Data(), &event); err != nil {
  span.RecordError(err)
  _ = req.RespondError("400", "invalid JSON: " + err.Error(), nil)
  return
 }

 if event.PermissionsID == "" {
  err := fmt.Errorf("invalid request: missing permissions ID")
  span.RecordError(err)
  _ = req.RespondError("400", err.Error(), nil)
  return
 }

 span.SetAttributes(
  tracer.StringAttr("permissions.id", event.PermissionsID),
  tracer.StringAttr("tenant.id", req.Header(core.HeaderTenantID)),
 )

 dalSubject := s.subjectPrefix + ".permissions.db.get"
 outMsg := &nats.Msg{Data: req.Data()}
 outMsg.Header = core.ExtractHeaders(ctx, nil)
 outMsg.Header.Set("X-Business-Validated", "true")

 dalCtx, dalCancel := context.WithTimeout(ctx, 10*time.Second)
 defer dalCancel()

 reply, err := s.publisher.Request(dalCtx, dalSubject, outMsg)
 if err != nil {
  span.RecordError(err)
  _ = req.RespondError("500", "DAL request error: " + err.Error(), nil)
  return
 }

logger.Info(ctx, "DAL reply received", logger.String("handler", "Permissions.get"), logger.Int("bytes", len(reply.Data)))
 _ = req.Respond(reply.Data)
}
