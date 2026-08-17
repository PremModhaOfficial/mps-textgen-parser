```
text gen component for concept Relation {
file name :
  ${node.from.name}${node.to.name}Handler
file path :
extension :
  go
(node)->void {

// === Top-level variables ===
string fromName = node.from.name;
string toName = node.to.name;
string fromLower = node.from.name.toLowerCase();
string toLower = node.to.name.toLowerCase();

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

// === Event / Request Types ===
foreach op in node.operations {

if (op.relationOperation == RelationOperation:assign) {
append {type } ${fromName} ${toName} {AssignedEvent struct {\n} ;
append {\t} ${fromName} {ID string `json:"} ${fromLower} {_id"`\n} ;
append {\t} ${toName} {ID string `json:"} ${toLower} {_id"`\n} ;
append {\tTimestamp time.Time `json:"timestamp"`\n} ;
append {}\n} ;
append {\n} ;
}

if (op.relationOperation == RelationOperation:remove) {
append {type } ${fromName} ${toName} {RemovedEvent struct {\n} ;
append {\t} ${fromName} {ID string `json:"} ${fromLower} {_id"`\n} ;
append {\t} ${toName} {ID string `json:"} ${toLower} {_id"`\n} ;
append {\tTimestamp time.Time `json:"timestamp"`\n} ;
append {}\n} ;
append {\n} ;
}

if (op.relationOperation == RelationOperation:list) {
append {type } ${fromName} ${toName} {ListRequest struct {\n} ;
append {\t} ${fromName} {ID string `json:"} ${fromLower} {_id"`\n} ;
append {\tLimit     int       `json:"limit"`\n} ;
append {\tOffset    int       `json:"offset"`\n} ;
append {\tTimestamp time.Time `json:"timestamp"`\n} ;
append {}\n} ;
append {\n} ;
}

}

// === Handler Struct + Constructor ===
append {type } ${fromName} ${toName} {Handler struct {\n} ;
append {\tpublisher     *events.Publisher\n} ;
append {\tsubjectPrefix string\n} ;
append {}\n} ;
append {\n} ;
append {func New} ${fromName} ${toName} {Handler(pub *events.Publisher, subjectPrefix string) *} ${fromName} ${toName} {Handler {\n} ;
append {\treturn &} ${fromName} ${toName} {Handler{\n} ;
append {\t\tpublisher:     pub,\n} ;
append {\t\tsubjectPrefix: subjectPrefix,\n} ;
append {\t}\n} ;
append {}\n} ;
append {\n} ;

// === Handler Methods ===
foreach op in node.operations {
string opName = op.capitalizedName();
string opKind = op.relationOperation.name;

append {func (s *} ${fromName} ${toName} {Handler) Handle} ${opName} {(req core.Request) {\n} ;
append {\tctx := req.Context()\n} ;
append {\tctx, span := tracer.StartConsumer(ctx, "} ${fromName} ${toName} {.Handle} ${opName} {")\n} ;
append {\tdefer span.End()\n} ;
append {\tctx = core.InjectContext(ctx, req.Headers())\n} ;
append {\n} ;

// --- Unmarshal ---
if (op.relationOperation == RelationOperation:assign) {
append {\tvar event } ${fromName} ${toName} {AssignedEvent\n} ;
}
if (op.relationOperation == RelationOperation:remove) {
append {\tvar event } ${fromName} ${toName} {RemovedEvent\n} ;
}
if (op.relationOperation == RelationOperation:list) {
append {\tvar event } ${fromName} ${toName} {ListRequest\n} ;
}

append {\tif err := json.Unmarshal(req.Data(), &event); err != nil {\n} ;
append {\t\tspan.RecordError(err)\n} ;
append {\t\t_ = req.RespondError("400", "invalid JSON: " + err.Error(), nil)\n} ;
append {\t\treturn\n} ;
append {\t}\n} ;
append {\n} ;

// --- Validation ---
if (op.relationOperation == RelationOperation:assign || op.relationOperation == RelationOperation:remove) {
append {\tif event.} ${fromName} {ID == "" || event.} ${toName} {ID == "" {\n} ;
append {\t\terr := fmt.Errorf("invalid data: missing } ${fromLower} { or } ${toLower} { ID")\n} ;
append {\t\tspan.RecordError(err)\n} ;
append {\t\t_ = req.RespondError("400", err.Error(), nil)\n} ;
append {\t\treturn\n} ;
append {\t}\n} ;
}
if (op.relationOperation == RelationOperation:list) {
append {\tif event.} ${fromName} {ID == "" {\n} ;
append {\t\terr := fmt.Errorf("invalid request: missing } ${fromLower} { ID")\n} ;
append {\t\tspan.RecordError(err)\n} ;
append {\t\t_ = req.RespondError("400", err.Error(), nil)\n} ;
append {\t\treturn\n} ;
append {\t}\n} ;
}

// --- Span Attributes ---
append {\n} ;
append {\tspan.SetAttributes(\n} ;
if (op.relationOperation == RelationOperation:assign || op.relationOperation == RelationOperation:remove) {
append {\t\ttracer.StringAttr("} ${fromLower} {.id", event.} ${fromName} {ID),\n} ;
append {\t\ttracer.StringAttr("} ${toLower} {.id", event.} ${toName} {ID),\n} ;
}
if (op.relationOperation == RelationOperation:list) {
append {\t\ttracer.StringAttr("} ${fromLower} {.id", event.} ${fromName} {ID),\n} ;
}
append {\t\ttracer.StringAttr("tenant.id", req.Header(core.HeaderTenantID)),\n} ;
append {\t)\n} ;
append {\n} ;

// --- Pre-hook ---
append {\tif err := s.pre} ${opName} {Hook(ctx, span, &event); err != nil {\n} ;
append {\t\tspan.RecordError(err)\n} ;
append {\t\t_ = req.RespondError("400", "pre-hook: " + err.Error(), nil)\n} ;
append {\t\treturn\n} ;
append {\t}\n} ;
append {\n} ;

// --- DAL forwarding via Request-Reply ---
append {\tdalSubject := s.subjectPrefix + ".} ${fromLower} {.} ${toLower} {.db.} ${opKind} {"\n} ;
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
append {\tlog.Printf("} ${fromName} ${toName} {.} ${opKind} { DAL reply: %d bytes", len(reply.Data))\n} ;
append {\n} ;

// --- Post-hook ---
append {\tresponseData := s.post} ${opName} {Hook(ctx, span, &event, reply.Data)\n} ;
append {\t_ = req.Respond(responseData)\n} ;
append {}\n} ;
append {\n} ;

}

// === Pre/Post Hook Stubs ===
foreach op in node.operations {
string hookName = op.capitalizedName();

if (op.relationOperation == RelationOperation:assign) {
append {func (s *} ${fromName} ${toName} {Handler) pre} ${hookName} {Hook(ctx context.Context, span tracer.Span, event *} ${fromName} ${toName} {AssignedEvent) error {\n} ;
append {\treturn nil\n} ;
append {}\n} ;
append {\n} ;
append {func (s *} ${fromName} ${toName} {Handler) post} ${hookName} {Hook(ctx context.Context, span tracer.Span, event *} ${fromName} ${toName} {AssignedEvent, data []byte) []byte {\n} ;
append {\treturn data\n} ;
append {}\n} ;
append {\n} ;
}

if (op.relationOperation == RelationOperation:remove) {
append {func (s *} ${fromName} ${toName} {Handler) pre} ${hookName} {Hook(ctx context.Context, span tracer.Span, event *} ${fromName} ${toName} {RemovedEvent) error {\n} ;
append {\treturn nil\n} ;
append {}\n} ;
append {\n} ;
append {func (s *} ${fromName} ${toName} {Handler) post} ${hookName} {Hook(ctx context.Context, span tracer.Span, event *} ${fromName} ${toName} {RemovedEvent, data []byte) []byte {\n} ;
append {\treturn data\n} ;
append {}\n} ;
append {\n} ;
}

if (op.relationOperation == RelationOperation:list) {
append {func (s *} ${fromName} ${toName} {Handler) pre} ${hookName} {Hook(ctx context.Context, span tracer.Span, event *} ${fromName} ${toName} {ListRequest) error {\n} ;
append {\treturn nil\n} ;
append {}\n} ;
append {\n} ;
append {func (s *} ${fromName} ${toName} {Handler) post} ${hookName} {Hook(ctx context.Context, span tracer.Span, event *} ${fromName} ${toName} {ListRequest, data []byte) []byte {\n} ;
append {\treturn data\n} ;
append {}\n} ;
append {\n} ;
}

}

}
}
```
