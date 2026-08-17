# User Management NATS Bridge - POC Context

**Working Directory**: `/home/prem-modha/projects/dsl/src/code/um/`

## POC Scope

This is a template/POC code to validate the architecture before DSL generation.

### Files to Create

```
um/
├── main.go              # NATS bridge with extensible generics
├── user.go              # User entity template
├── role.go              # Role entity template
├── init.sql             # SQL schema template
├── go.mod
└── go.sum
```

## Architecture Pattern

### NATS Bridge (main.go)
- Generic handler registration system
- Subject mapping: `um.{entity}.{operation}` → `{tenant}.{entity}.{operation}`
- Request-reply forwarding to DAL
- Header injection (tenant ID, trace context, correlation ID)
- OTEL span creation with semantic attributes

### Entity Templates (user.go, role.go)
- Struct definition with JSON/DB tags
- Event types (CreatedEvent, UpdatedEvent, DeletedEvent, GetRequest, ListRequest)
- Handler methods with validation
- NATS message publishing to DAL subjects

### Database Schema (init.sql)
- User table: id, username, email, password, created_at, updated_at
- Role table: id, name, permissions (JSON array), created_at, updated_at
- UserRole junction table: user_id, role_id (for RBAC)

## Integration Points

### Go SDK (`motadata-go-sdk`)
- `core.Message` - message structure
- `core.Publisher` - publish to NATS
- `core.Subscriber` - subscribe to topics
- `core.Headers` - metadata headers

### OpenTelemetry
- `trace.Tracer` - create spans
- `attribute` - add semantic attributes
- `span.RecordError()` - error tracking

## Subject Pattern

**UM Layer**:
```
um.user.create
um.user.get
um.user.update
um.user.delete
um.user.list

um.role.create
um.role.get
um.role.update
um.role.delete
um.role.list
```

**DAL Layer**:
```
{tenantID}.user.create
{tenantID}.user.get
{tenantID}.user.update
{tenantID}.user.delete
{tenantID}.user.list

{tenantID}.role.create
{tenantID}.role.get
{tenantID}.role.update
{tenantID}.role.delete
{tenantID}.role.list
```

## Header Handling

**Injected by bridge**:
- `X-Tenant-ID` - from env or request
- `X-Trace-ID` - generated or propagated
- `X-Span-ID` - from span context
- `X-Correlation-ID` - request correlation

**Added by entity handlers**:
- `X-Business-Validated` - set to "true" after validation

## Next Phase: DSL Generation

Once POC is validated, this code will serve as template for:
- Entity TextGen templates
- Relation TextGen templates
- SqlSchema TextGen templates
