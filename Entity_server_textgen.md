```
text gen component for concept Entity {
file name :
  ${node.name.toLowerCase()}
file path :
extension :
  go
(node)->void {

// === Top-level variables ===
string name = node.name;
string nameLower = node.name.toLowerCase();
string pkField = node.primaryKeyField().name;

// === Package + Imports ===
append {package main\n} ;
append {\n} ;
append {import (\n} ;
append {\t"context"\n} ;
append {\t"encoding/json"\n} ;
append {\t"fmt"\n} ;
append {\t"log"\n} ;
append {\t"time"\n} ;
append {\n} ;
append {\t"github.com/nats-io/nats.go"\n} ;
append {\n} ;
append {\t"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events"\n} ;
append {\t"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/core"\n} ;
append {\t"dev.azure.com/Motadata/NextGen/motadata-go-sdk/otel/tracer"\n} ;
append {)\n} ;
append {\n} ;

// === Domain Struct ===
append {type } ${name} { struct {\n} ;
foreach field in node.fields {
append {\t} ${field.name.capitalize()} { } ${field.goType()} { `json:"} ${field.jsonName()} {" db:"} ${field.dbName()} {"`\n} ;
}
append {}\n} ;
append {\n} ;

// === Event / Request Types ===
foreach op in node.operations {

if (op.entityOperation == EntityOperation:create) {
append {type } ${name} {CreatedEvent struct {\n} ;
append {\t} ${name} { } ${name} { `json:"} ${nameLower} {"`\n} ;
append {\tTimestamp time.Time `json:"timestamp"`\n} ;
append {}\n} ;
append {\n} ;
}

if (op.entityOperation == EntityOperation:update) {
append {type } ${name} {UpdatedEvent struct {\n} ;
append {\t} ${name} { } ${name} { `json:"} ${nameLower} {"`\n} ;
append {\tTimestamp time.Time `json:"timestamp"`\n} ;
append {}\n} ;
append {\n} ;
}

if (op.entityOperation == EntityOperation:delete) {
append {type } ${name} {DeletedEvent struct {\n} ;
append {\t} ${name} {ID string `json:"} ${nameLower} {_id"`\n} ;
append {\tTimestamp time.Time `json:"timestamp"`\n} ;
append {}\n} ;
append {\n} ;
}

if (op.entityOperation == EntityOperation:list) {
append {type } ${name} {ListRequest struct {\n} ;
append {\tLimit     int       `json:"limit"`\n} ;
append {\tOffset    int       `json:"offset"`\n} ;
append {\tTimestamp time.Time `json:"timestamp"`\n} ;
append {}\n} ;
append {\n} ;
}

if (op.entityOperation == EntityOperation:get) {
append {type } ${name} {GetRequest struct {\n} ;
append {\t} ${name} {ID string `json:"} ${nameLower} {_id"`\n} ;
append {\tTimestamp time.Time `json:"timestamp"`\n} ;
append {}\n} ;
append {\n} ;
}

}

// === Handler Struct + Constructor ===
append {type } ${name} {Handler struct {\n} ;
append {\tpublisher     *events.Publisher\n} ;
append {\tsubjectPrefix string\n} ;
append {}\n} ;
append {\n} ;
append {func New} ${name} {Handler(pub *events.Publisher, subjectPrefix string) *} ${name} {Handler {\n} ;
append {\treturn &} ${name} {Handler{\n} ;
append {\t\tpublisher:     pub,\n} ;
append {\t\tsubjectPrefix: subjectPrefix,\n} ;
append {\t}\n} ;
append {}\n} ;
append {\n} ;

// === Handler Methods ===
foreach op in node.operations {
string opName = op.capitalizedName();
string opKind = op.entityOperation.name;

append {func (s *} ${name} {Handler) Handle} ${opName} {(req core.Request) {\n} ;
append {\tctx := req.Context()\n} ;
append {\tctx, span := tracer.StartConsumer(ctx, "} ${name} {.Handle} ${opName} {")\n} ;
append {\tdefer span.End()\n} ;
append {\tctx = core.InjectContext(ctx, req.Headers())\n} ;
append {\n} ;

// --- Unmarshal into correct event type ---
if (op.entityOperation == EntityOperation:create) {
append {\tvar event } ${name} {CreatedEvent\n} ;
}
if (op.entityOperation == EntityOperation:update) {
append {\tvar event } ${name} {UpdatedEvent\n} ;
}
if (op.entityOperation == EntityOperation:delete) {
append {\tvar event } ${name} {DeletedEvent\n} ;
}
if (op.entityOperation == EntityOperation:list) {
append {\tvar event } ${name} {ListRequest\n} ;
}
if (op.entityOperation == EntityOperation:get) {
append {\tvar event } ${name} {GetRequest\n} ;
}

append {\tif err := json.Unmarshal(req.Data(), &event); err != nil {\n} ;
append {\t\tspan.RecordError(err)\n} ;
append {\t\t_ = req.RespondError("400", "invalid JSON: " + err.Error(), nil)\n} ;
append {\t\treturn\n} ;
append {\t}\n} ;
append {\n} ;

// --- Validation ---
if (op.entityOperation == EntityOperation:create) {
int valIdx = 0;
foreach field in node.fields {
if (!(field.hasAnnotation(FieldAnnotation:primaryKey)) && !(field.hasAnnotation(FieldAnnotation:auto)) && !(field.hasAnnotation(FieldAnnotation:hidden)) && !(field.hasAnnotation(FieldAnnotation:nullable))) {
if (valIdx == 0) {
append {\tif event.} ${name} {.} ${field.name.capitalize()} { == ""} ;
}
if (valIdx > 0) {
append { || event.} ${name} {.} ${field.name.capitalize()} { == ""} ;
}
valIdx = valIdx + 1;
}
}
append { {\n} ;
append {\t\terr := fmt.Errorf("invalid } ${nameLower} { data: missing required fields")\n} ;
append {\t\tspan.RecordError(err)\n} ;
append {\t\t_ = req.RespondError("400", err.Error(), nil)\n} ;
append {\t\treturn\n} ;
append {\t}\n} ;
}

if (op.entityOperation == EntityOperation:update) {
append {\tif event.} ${name} {.} ${pkField} { == "" {\n} ;
append {\t\terr := fmt.Errorf("invalid } ${nameLower} { data: missing ID")\n} ;
append {\t\tspan.RecordError(err)\n} ;
append {\t\t_ = req.RespondError("400", err.Error(), nil)\n} ;
append {\t\treturn\n} ;
append {\t}\n} ;
}

if (op.entityOperation == EntityOperation:delete) {
append {\tif event.} ${name} {ID == "" {\n} ;
append {\t\terr := fmt.Errorf("invalid request: missing } ${nameLower} { ID")\n} ;
append {\t\tspan.RecordError(err)\n} ;
append {\t\t_ = req.RespondError("400", err.Error(), nil)\n} ;
append {\t\treturn\n} ;
append {\t}\n} ;
}

if (op.entityOperation == EntityOperation:get) {
append {\tif event.} ${name} {ID == "" {\n} ;
append {\t\terr := fmt.Errorf("invalid request: missing } ${nameLower} { ID")\n} ;
append {\t\tspan.RecordError(err)\n} ;
append {\t\t_ = req.RespondError("400", err.Error(), nil)\n} ;
append {\t\treturn\n} ;
append {\t}\n} ;
}

if (op.entityOperation == EntityOperation:list) {
append {\tif event.Limit < 0 || event.Offset < 0 {\n} ;
append {\t\terr := fmt.Errorf("invalid pagination parameters")\n} ;
append {\t\tspan.RecordError(err)\n} ;
append {\t\t_ = req.RespondError("400", err.Error(), nil)\n} ;
append {\t\treturn\n} ;
append {\t}\n} ;
}

// --- Span Attributes ---
append {\n} ;
append {\tspan.SetAttributes(\n} ;
if (op.entityOperation == EntityOperation:create) {
append {\t\ttracer.StringAttr("} ${nameLower} {.id", event.} ${name} {.} ${pkField} {),\n} ;
}
if (op.entityOperation == EntityOperation:update) {
append {\t\ttracer.StringAttr("} ${nameLower} {.id", event.} ${name} {.} ${pkField} {),\n} ;
}
if (op.entityOperation == EntityOperation:delete) {
append {\t\ttracer.StringAttr("} ${nameLower} {.id", event.} ${name} {ID),\n} ;
}
if (op.entityOperation == EntityOperation:get) {
append {\t\ttracer.StringAttr("} ${nameLower} {.id", event.} ${name} {ID),\n} ;
}
append {\t\ttracer.StringAttr("tenant.id", req.Header(core.HeaderTenantID)),\n} ;
append {\t)\n} ;
append {\n} ;

// --- Pre-hook (only if operation is in node.prehooks) ---
foreach prehook in node.prehooks {
if (prehook.entityOperation == op.entityOperation) {
append {\tif err := s.pre} ${opName} {Hook(ctx, span, &event); err != nil {\n} ;
append {\t\tspan.RecordError(err)\n} ;
append {\t\t_ = req.RespondError("400", "pre-hook: " + err.Error(), nil)\n} ;
append {\t\treturn\n} ;
append {\t}\n} ;
append {\n} ;
}
}

// --- DAL forwarding via Request-Reply ---
append {\tdalSubject := s.subjectPrefix + ".} ${nameLower} {.db.} ${opKind} {"\n} ;
append {\toutMsg := &nats.Msg{Data: req.Data()}\n} ;
append {\toutMsg.Header = core.ExtractHeaders(ctx, nil)\n} ;
append {\toutMsg.Header.Set("X-Business-Validated", "true")\n} ;
append {\n} ;
append {\tdalCtx, dalCancel := context.WithTimeout(ctx, 10*time.Second)\n} ;
append {\tdefer dalCancel()\n} ;
append {\n} ;
append {\treply, err := s.publisher.Request(dalCtx, dalSubject, outMsg)\n} ;
append {\tif err != nil {\n} ;
append {\t\tspan.RecordError(err)\n} ;
append {\t\t_ = req.RespondError("500", "DAL request error: " + err.Error(), nil)\n} ;
append {\t\treturn\n} ;
append {\t}\n} ;
append {\n} ;
append {\tlog.Printf("} ${name} {.} ${opKind} { DAL reply: %d bytes", len(reply.Data))\n} ;
append {\n} ;

// --- Post-hook or direct respond ---
boolean hasPostHook = false;
foreach posthook in node.posthooks {
if (posthook.entityOperation == op.entityOperation) {
hasPostHook = true;
}
}
if (hasPostHook) {
append {\tresponseData := s.post} ${opName} {Hook(ctx, span, &event, reply.Data)\n} ;
append {\t_ = req.Respond(responseData)\n} ;
}
if (!(hasPostHook)) {
append {\t_ = req.Respond(reply.Data)\n} ;
}
append {}\n} ;
append {\n} ;

}

// === Pre-Hook Stubs (only for operations in node.prehooks) ===
foreach op in node.prehooks {
string hookName = op.capitalizedName();

if (op.entityOperation == EntityOperation:create) {
append {func (s *} ${name} {Handler) pre} ${hookName} {Hook(ctx context.Context, span tracer.Span, event *} ${name} {CreatedEvent) error {\n} ;
append {\treturn nil\n} ;
append {}\n} ;
append {\n} ;
}

if (op.entityOperation == EntityOperation:update) {
append {func (s *} ${name} {Handler) pre} ${hookName} {Hook(ctx context.Context, span tracer.Span, event *} ${name} {UpdatedEvent) error {\n} ;
append {\treturn nil\n} ;
append {}\n} ;
append {\n} ;
}

if (op.entityOperation == EntityOperation:delete) {
append {func (s *} ${name} {Handler) pre} ${hookName} {Hook(ctx context.Context, span tracer.Span, event *} ${name} {DeletedEvent) error {\n} ;
append {\treturn nil\n} ;
append {}\n} ;
append {\n} ;
}

if (op.entityOperation == EntityOperation:get) {
append {func (s *} ${name} {Handler) pre} ${hookName} {Hook(ctx context.Context, span tracer.Span, event *} ${name} {GetRequest) error {\n} ;
append {\treturn nil\n} ;
append {}\n} ;
append {\n} ;
}

if (op.entityOperation == EntityOperation:list) {
append {func (s *} ${name} {Handler) pre} ${hookName} {Hook(ctx context.Context, span tracer.Span, event *} ${name} {ListRequest) error {\n} ;
append {\treturn nil\n} ;
append {}\n} ;
append {\n} ;
}

}

// === Post-Hook Stubs (only for operations in node.posthooks) ===
foreach op in node.posthooks {
string hookName = op.capitalizedName();

if (op.entityOperation == EntityOperation:create) {
append {func (s *} ${name} {Handler) post} ${hookName} {Hook(ctx context.Context, span tracer.Span, event *} ${name} {CreatedEvent, data []byte) []byte {\n} ;
append {\treturn data\n} ;
append {}\n} ;
append {\n} ;
}

if (op.entityOperation == EntityOperation:update) {
append {func (s *} ${name} {Handler) post} ${hookName} {Hook(ctx context.Context, span tracer.Span, event *} ${name} {UpdatedEvent, data []byte) []byte {\n} ;
append {\treturn data\n} ;
append {}\n} ;
append {\n} ;
}

if (op.entityOperation == EntityOperation:delete) {
append {func (s *} ${name} {Handler) post} ${hookName} {Hook(ctx context.Context, span tracer.Span, event *} ${name} {DeletedEvent, data []byte) []byte {\n} ;
append {\treturn data\n} ;
append {}\n} ;
append {\n} ;
}

if (op.entityOperation == EntityOperation:get) {
append {func (s *} ${name} {Handler) post} ${hookName} {Hook(ctx context.Context, span tracer.Span, event *} ${name} {GetRequest, data []byte) []byte {\n} ;
append {\treturn data\n} ;
append {}\n} ;
append {\n} ;
}

if (op.entityOperation == EntityOperation:list) {
append {func (s *} ${name} {Handler) post} ${hookName} {Hook(ctx context.Context, span tracer.Span, event *} ${name} {ListRequest, data []byte) []byte {\n} ;
append {\treturn data\n} ;
append {}\n} ;
append {\n} ;
}

}

}
}
```
