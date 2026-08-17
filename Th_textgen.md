text gen component for concept Th {
file name :
  generated
file path :
extension :
  go
(node)->void {

```
append {package main\n} ;
append {\n} ;
append {import (\n} ;
append {\t"context"\n} ;
append {\t"encoding/json"\n} ;
append {\t"fmt"\n} ;
append {\t"time"\n} ;
append {\n} ;
append {\t"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/core"\n} ;
append {\t"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/transport/nats"\n} ;
append {\t"go.opentelemetry.io/otel/attribute"\n} ;
append {\t"go.opentelemetry.io/otel/trace"\n} ;
append {)\n} ;
append {\n} ;

append {// DAL interface\n} ;
append {type DataAccessLayer interface {\n} ;
append {\tGetModelByID(ctx context.Context, id string) (*Model, error)\n} ;
append {}\n} ;
append {\n} ;

append {// Domain Model\n} ;
append {type Model struct {\n} ;
append {\tID        string    `json:"id" db:"id"`\n} ;
append {\tName      string    `json:"name" db:"name"`\n} ;
append {\tCreatedAt time.Time `json:"created_at" db:"created_at"`\n} ;
append {}\n} ;
append {\n} ;

append {// Event Payload\n} ;
append {type ModelCreatedEvent struct {\n} ;
append {\tModel     Model     `json:"model"`\n} ;
append {\tTimestamp time.Time `json:"timestamp"`\n} ;
append {}\n} ;
append {\n} ;

append {// Service Interface\n} ;
append {type ModelService interface {\n} ;
append {\tCreate(ctx context.Context, model Model) error\n} ;
append {}\n} ;
append {\n} ;

append {// Service Implementation\n} ;
append {type modelService struct {\n} ;
append {\tpublisher *nats.Publisher\n} ;
append {\ttracer    trace.Tracer\n} ;
append {\tdal       DataAccessLayer\n} ;
append {}\n} ;
append {\n} ;

append {func NewModelService(pub *nats.Publisher, t trace.Tracer, dal DataAccessLayer) ModelService {\n} ;
append {\treturn &modelService{\n} ;
append {\t\tpublisher: pub,\n} ;
append {\t\ttracer:    t,\n} ;
append {\t\tdal:       dal,\n} ;
append {\t}\n} ;
append {}\n} ;
append {\n} ;

append {func (s *modelService) Create(ctx context.Context, model Model) error {\n} ;
append {\tctx, span := s.tracer.Start(ctx, "ModelService.Create")\n} ;
append {\tdefer span.End()\n} ;
append {\tspan.SetAttributes(\n} ;
append {\t\tattribute.String("model.id", model.ID),\n} ;
append {\t)\n} ;
append {\n} ;
append {\tevent := ModelCreatedEvent{\n} ;
append {\t\tModel:     model,\n} ;
append {\t\tTimestamp: time.Now().UTC(),\n} ;
append {\t}\n} ;
append {\n} ;
append {\tpayload, err := json.Marshal(event)\n} ;
append {\tif err != nil {\n} ;
append {\t\tspan.RecordError(err)\n} ;
append {\t\treturn fmt.Errorf("failed to marshal event: %w", err)\n} ;
append {\t}\n} ;
append {\n} ;
append {\tmsg := core.NewMessage(payload)\n} ;
append {\tmsg.Subject = "model.db.create"\n} ;
append {\tif err := s.publisher.Publish(ctx, msg.Subject, msg); err != nil {\n} ;
append {\t\tspan.RecordError(err)\n} ;
append {\t\treturn fmt.Errorf("failed to publish event: %w", err)\n} ;
append {\t}\n} ;
append {\treturn nil\n} ;
append {}\n} ;
append {\n} ;

append {var _ = context.Background\n} ;
append {var _ = fmt.Sprintf\n} ;
```

}
}
