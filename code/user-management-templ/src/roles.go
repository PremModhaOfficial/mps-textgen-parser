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

type Roles struct {
 ID string `json:"id" db:"id"`
 Desc string `json:"desc" db:"desc"`
 Name string `json:"name" db:"name"`
}

type RolesCreatedEvent struct {
 Roles Roles `json:"roles"`
 Timestamp time.Time `json:"timestamp"`
}

type RolesUpdatedEvent struct {
 Roles Roles `json:"roles"`
 Timestamp time.Time `json:"timestamp"`
}

type RolesListRequest struct {
 Limit     int       `json:"limit"`
 Offset    int       `json:"offset"`
 Timestamp time.Time `json:"timestamp"`
}

type RolesDeletedEvent struct {
 RolesID string `json:"roles_id"`
 Timestamp time.Time `json:"timestamp"`
}

type RolesGetRequest struct {
 RolesID string `json:"roles_id"`
 Timestamp time.Time `json:"timestamp"`
}

type RolesHandler struct {
 publisher     *events.Publisher
 subjectPrefix string
}

func NewRolesHandler(pub *events.Publisher, subjectPrefix string) *RolesHandler {
 return &RolesHandler{
  publisher:     pub,
  subjectPrefix: subjectPrefix,
 }
}

func (s *RolesHandler) HandleCreate(req core.Request) {
 ctx := req.Context()
 ctx, span := tracer.StartConsumer(ctx, "Roles.HandleCreate")
 defer span.End()
 ctx = core.InjectContext(ctx, req.Headers())
logger.Info(ctx, "request received", logger.String("handler", "Roles.Create"), logger.Int("bytes", len(req.Data())))

 var event RolesCreatedEvent
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

 if event.Roles.Name == "" {
  err := fmt.Errorf("invalid roles data: missing required fields")
  span.RecordError(err)
  _ = req.RespondError("400", err.Error(), nil)
  return
 }

 span.SetAttributes(
  tracer.StringAttr("roles.id", event.Roles.ID),
  tracer.StringAttr("tenant.id", req.Header(core.HeaderTenantID)),
 )

 dalSubject := s.subjectPrefix + ".roles.db.create"
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

logger.Info(ctx, "DAL reply received", logger.String("handler", "Roles.create"), logger.Int("bytes", len(reply.Data)))
 _ = req.Respond(reply.Data)
}

func (s *RolesHandler) HandleUpdate(req core.Request) {
 ctx := req.Context()
 ctx, span := tracer.StartConsumer(ctx, "Roles.HandleUpdate")
 defer span.End()
 ctx = core.InjectContext(ctx, req.Headers())
logger.Info(ctx, "request received", logger.String("handler", "Roles.Update"), logger.Int("bytes", len(req.Data())))

 var event RolesUpdatedEvent
 if err := json.Unmarshal(req.Data(), &event); err != nil {
  span.RecordError(err)
  _ = req.RespondError("400", "invalid JSON: " + err.Error(), nil)
  return
 }

 if event.Roles.ID == "" {
  err := fmt.Errorf("invalid roles data: missing ID")
  span.RecordError(err)
  _ = req.RespondError("400", err.Error(), nil)
  return
 }

 span.SetAttributes(
  tracer.StringAttr("roles.id", event.Roles.ID),
  tracer.StringAttr("tenant.id", req.Header(core.HeaderTenantID)),
 )

 dalSubject := s.subjectPrefix + ".roles.db.update"
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

logger.Info(ctx, "DAL reply received", logger.String("handler", "Roles.update"), logger.Int("bytes", len(reply.Data)))
 _ = req.Respond(reply.Data)
}

func (s *RolesHandler) HandleList(req core.Request) {
 ctx := req.Context()
 ctx, span := tracer.StartConsumer(ctx, "Roles.HandleList")
 defer span.End()
 ctx = core.InjectContext(ctx, req.Headers())
logger.Info(ctx, "request received", logger.String("handler", "Roles.List"), logger.Int("bytes", len(req.Data())))

 var event RolesListRequest
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

 dalSubject := s.subjectPrefix + ".roles.db.list"
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

logger.Info(ctx, "DAL reply received", logger.String("handler", "Roles.list"), logger.Int("bytes", len(reply.Data)))
 _ = req.Respond(reply.Data)
}

func (s *RolesHandler) HandleDelete(req core.Request) {
 ctx := req.Context()
 ctx, span := tracer.StartConsumer(ctx, "Roles.HandleDelete")
 defer span.End()
 ctx = core.InjectContext(ctx, req.Headers())
logger.Info(ctx, "request received", logger.String("handler", "Roles.Delete"), logger.Int("bytes", len(req.Data())))

 var event RolesDeletedEvent
 if err := json.Unmarshal(req.Data(), &event); err != nil {
  span.RecordError(err)
  _ = req.RespondError("400", "invalid JSON: " + err.Error(), nil)
  return
 }

 if event.RolesID == "" {
  err := fmt.Errorf("invalid request: missing roles ID")
  span.RecordError(err)
  _ = req.RespondError("400", err.Error(), nil)
  return
 }

 span.SetAttributes(
  tracer.StringAttr("roles.id", event.RolesID),
  tracer.StringAttr("tenant.id", req.Header(core.HeaderTenantID)),
 )

 dalSubject := s.subjectPrefix + ".roles.db.delete"
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

logger.Info(ctx, "DAL reply received", logger.String("handler", "Roles.delete"), logger.Int("bytes", len(reply.Data)))
 _ = req.Respond(reply.Data)
}

func (s *RolesHandler) HandleGet(req core.Request) {
 ctx := req.Context()
 ctx, span := tracer.StartConsumer(ctx, "Roles.HandleGet")
 defer span.End()
 ctx = core.InjectContext(ctx, req.Headers())
logger.Info(ctx, "request received", logger.String("handler", "Roles.Get"), logger.Int("bytes", len(req.Data())))

 var event RolesGetRequest
 if err := json.Unmarshal(req.Data(), &event); err != nil {
  span.RecordError(err)
  _ = req.RespondError("400", "invalid JSON: " + err.Error(), nil)
  return
 }

 if event.RolesID == "" {
  err := fmt.Errorf("invalid request: missing roles ID")
  span.RecordError(err)
  _ = req.RespondError("400", err.Error(), nil)
  return
 }

 span.SetAttributes(
  tracer.StringAttr("roles.id", event.RolesID),
  tracer.StringAttr("tenant.id", req.Header(core.HeaderTenantID)),
 )

 dalSubject := s.subjectPrefix + ".roles.db.get"
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

logger.Info(ctx, "DAL reply received", logger.String("handler", "Roles.get"), logger.Int("bytes", len(reply.Data)))
 _ = req.Respond(reply.Data)
}
type RolesPermissionsAssignedEvent struct {
 RolesID string `json:"roles_id"`
 PermissionsID string `json:"permissions_id"`
 Timestamp time.Time `json:"timestamp"`
}

type RolesPermissionsListRequest struct {
 RolesID string `json:"roles_id"`
 Limit     int       `json:"limit"`
 Offset    int       `json:"offset"`
 Timestamp time.Time `json:"timestamp"`
}

type RolesPermissionsRemovedEvent struct {
 RolesID string `json:"roles_id"`
 PermissionsID string `json:"permissions_id"`
 Timestamp time.Time `json:"timestamp"`
}

type RolesPermissionsHandler struct {
 publisher     *events.Publisher
 subjectPrefix string
}

func NewRolesPermissionsHandler(pub *events.Publisher, subjectPrefix string) *RolesPermissionsHandler {
 return &RolesPermissionsHandler{
  publisher:     pub,
  subjectPrefix: subjectPrefix,
 }
}

func (s *RolesPermissionsHandler) HandleAssign(req core.Request) {
 ctx := req.Context()
 ctx, span := tracer.StartConsumer(ctx, "RolesPermissions.HandleAssign")
 defer span.End()
logger.Info(ctx, "request received", logger.String("handler", "RolesPermissions.Assign"), logger.Int("bytes", len(req.Data())))
 ctx = core.InjectContext(ctx, req.Headers())

 var event RolesPermissionsAssignedEvent
 if err := json.Unmarshal(req.Data(), &event); err != nil {
  span.RecordError(err)
  _ = req.RespondError("400", "invalid JSON: " + err.Error(), nil)
  return
 }

 if event.RolesID == "" || event.PermissionsID == "" {
  err := fmt.Errorf("invalid data: missing roles or permissions ID")
  span.RecordError(err)
  _ = req.RespondError("400", err.Error(), nil)
  return
 }

 span.SetAttributes(
  tracer.StringAttr("roles.id", event.RolesID),
  tracer.StringAttr("permissions.id", event.PermissionsID),
  tracer.StringAttr("tenant.id", req.Header(core.HeaderTenantID)),
 )

 if err := s.preAssignHook(ctx, span, &event); err != nil {
  span.RecordError(err)
  _ = req.RespondError("400", "pre-hook: " + err.Error(), nil)
  return
 }

 dalSubject := s.subjectPrefix + ".roles.permissions.db.assign"
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

 logger.Info(ctx, "DAL reply received", logger.String("handler", "RolesPermissions.assign"), logger.Int("bytes", len(reply.Data)))

 responseData := s.postAssignHook(ctx, span, &event, reply.Data)
 _ = req.Respond(responseData)
}

func (s *RolesPermissionsHandler) HandleList(req core.Request) {
 ctx := req.Context()
 ctx, span := tracer.StartConsumer(ctx, "RolesPermissions.HandleList")
 defer span.End()
logger.Info(ctx, "request received", logger.String("handler", "RolesPermissions.List"), logger.Int("bytes", len(req.Data())))
 ctx = core.InjectContext(ctx, req.Headers())

 var event RolesPermissionsListRequest
 if err := json.Unmarshal(req.Data(), &event); err != nil {
  span.RecordError(err)
  _ = req.RespondError("400", "invalid JSON: " + err.Error(), nil)
  return
 }

 if event.RolesID == "" {
  err := fmt.Errorf("invalid request: missing roles ID")
  span.RecordError(err)
  _ = req.RespondError("400", err.Error(), nil)
  return
 }

 span.SetAttributes(
  tracer.StringAttr("roles.id", event.RolesID),
  tracer.StringAttr("tenant.id", req.Header(core.HeaderTenantID)),
 )

 if err := s.preListHook(ctx, span, &event); err != nil {
  span.RecordError(err)
  _ = req.RespondError("400", "pre-hook: " + err.Error(), nil)
  return
 }

 dalSubject := s.subjectPrefix + ".roles.permissions.db.list"
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

 logger.Info(ctx, "DAL reply received", logger.String("handler", "RolesPermissions.list"), logger.Int("bytes", len(reply.Data)))

 responseData := s.postListHook(ctx, span, &event, reply.Data)
 _ = req.Respond(responseData)
}

func (s *RolesPermissionsHandler) HandleRemove(req core.Request) {
 ctx := req.Context()
 ctx, span := tracer.StartConsumer(ctx, "RolesPermissions.HandleRemove")
 defer span.End()
logger.Info(ctx, "request received", logger.String("handler", "RolesPermissions.Remove"), logger.Int("bytes", len(req.Data())))
 ctx = core.InjectContext(ctx, req.Headers())

 var event RolesPermissionsRemovedEvent
 if err := json.Unmarshal(req.Data(), &event); err != nil {
  span.RecordError(err)
  _ = req.RespondError("400", "invalid JSON: " + err.Error(), nil)
  return
 }

 if event.RolesID == "" || event.PermissionsID == "" {
  err := fmt.Errorf("invalid data: missing roles or permissions ID")
  span.RecordError(err)
  _ = req.RespondError("400", err.Error(), nil)
  return
 }

 span.SetAttributes(
  tracer.StringAttr("roles.id", event.RolesID),
  tracer.StringAttr("permissions.id", event.PermissionsID),
  tracer.StringAttr("tenant.id", req.Header(core.HeaderTenantID)),
 )

 if err := s.preRemoveHook(ctx, span, &event); err != nil {
  span.RecordError(err)
  _ = req.RespondError("400", "pre-hook: " + err.Error(), nil)
  return
 }

 dalSubject := s.subjectPrefix + ".roles.permissions.db.remove"
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

 logger.Info(ctx, "DAL reply received", logger.String("handler", "RolesPermissions.remove"), logger.Int("bytes", len(reply.Data)))

 responseData := s.postRemoveHook(ctx, span, &event, reply.Data)
 _ = req.Respond(responseData)
}
