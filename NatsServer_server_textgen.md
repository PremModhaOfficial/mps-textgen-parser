```
text gen component for concept NatsServer {
file name :
  main
file path :
extension :
  go
(node)->void {

// === Top-level variables ===
string svcName = node.name;
string prefix = node.subjectPrefix;
string clientID = node.clientID;
string natsUrl = node.defaultNatsUrl;

// === Package + Imports ===
append {package main\n} ;
append {\n} ;
append {import (\n} ;
append {\t"context"\n} ;
append {\t"log"\n} ;
append {\t"os"\n} ;
append {\t"os/signal"\n} ;
append {\t"syscall"\n} ;
append {\t"time"\n} ;
append {\n} ;
append {\t"github.com/nats-io/nats.go"\n} ;
append {\t"github.com/nats-io/nats.go/micro"\n} ;
append {\n} ;
append {\t"dev.azure.com/Motadata/NextGen/motadata-go-sdk/config"\n} ;
append {\t"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events"\n} ;
append {\t"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/auth"\n} ;
append {\t"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/core"\n} ;
append {\t"dev.azure.com/Motadata/NextGen/motadata-go-sdk/otel"\n} ;
append {)\n} ;
append {\n} ;

// === micro.Request -> core.Request adapter ===
append {type microRequestAdapter struct {\n} ;
append {\tmr micro.Request\n} ;
append {}\n} ;
append {\n} ;
append {func (a *microRequestAdapter) Context() context.Context { return context.Background() }\n} ;
append {func (a *microRequestAdapter) Subject() string          { return a.mr.Subject() }\n} ;
append {func (a *microRequestAdapter) Reply() string            { return "" }\n} ;
append {func (a *microRequestAdapter) Data() []byte             { return a.mr.Data() }\n} ;
append {func (a *microRequestAdapter) Headers() nats.Header     { return nats.Header(a.mr.Headers()) }\n} ;
append {func (a *microRequestAdapter) Header(key string) string { return a.mr.Headers().Get(key) }\n} ;
append {\n} ;
append {func (a *microRequestAdapter) Respond(data []byte, opts ...core.RespondOption) error {\n} ;
append {\treturn a.mr.Respond(data)\n} ;
append {}\n} ;
append {func (a *microRequestAdapter) RespondJSON(v any, opts ...core.RespondOption) error {\n} ;
append {\treturn a.mr.RespondJSON(v)\n} ;
append {}\n} ;
append {func (a *microRequestAdapter) RespondError(code, desc string, data []byte, opts ...core.RespondOption) error {\n} ;
append {\treturn a.mr.Error(code, desc, data)\n} ;
append {}\n} ;
append {\n} ;
append {func (a *microRequestAdapter) Ack() error                         { return nil }\n} ;
append {func (a *microRequestAdapter) Nak() error                         { return nil }\n} ;
append {func (a *microRequestAdapter) NakWithDelay(_ time.Duration) error { return nil }\n} ;
append {func (a *microRequestAdapter) Term() error                        { return nil }\n} ;
append {func (a *microRequestAdapter) InProgress() error                  { return nil }\n} ;
append {\n} ;
append {func adaptRequest(mr micro.Request) core.Request {\n} ;
append {\treturn &microRequestAdapter{mr: mr}\n} ;
append {}\n} ;
append {\n} ;

// === main() ===
append {const prefix = "} ${prefix} {"\n} ;
append {\n} ;
append {func main() {\n} ;
append {\tlog.Println("Starting } ${svcName} { (Business Logic Layer)...")\n} ;
append {\n} ;
append {\totel.InitFromEnv()\n} ;
append {\n} ;
append {\tnatsURL := os.Getenv("NATS_URL")\n} ;
append {\tif natsURL == "" {\n} ;
append {\t\tnatsURL = "} ${natsUrl} {"\n} ;
append {\t}\n} ;
append {\n} ;
append {\tcfg := config.DefaultEventsConfig()\n} ;
append {\tcfg.Servers = []string{natsURL}\n} ;
append {\n} ;
append {\tcreds := &auth.Credentials{Type: auth.TypeNone}\n} ;
append {\tconn := events.NewConnection("} ${clientID} {", cfg, creds)\n} ;
append {\n} ;
append {\tctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)\n} ;
append {\tdefer cancel()\n} ;
append {\n} ;
append {\tif err := conn.Connect(ctx); err != nil {\n} ;
append {\t\tlog.Fatalf("failed to connect to NATS: %v", err)\n} ;
append {\t}\n} ;
append {\tdefer conn.Close(context.Background())\n} ;
append {\n} ;
append {\tpublisher := events.NewPublisher(conn)\n} ;
append {\n} ;

// --- Instantiate entity handlers ---
foreach entity in node.entities {
string eName = entity.name;
string eVar = entity.name.toLowerCaseFirst();
append {\t} ${eVar} {Handler := New} ${eName} {Handler(publisher, prefix)\n} ;
}

// --- Instantiate relation handlers ---
foreach relation in node.relations {
string rFrom = relation.from.name;
string rTo = relation.to.name;
string rVar = relation.from.name.toLowerCaseFirst();
append {\t} ${rVar} ${rTo} {Handler := New} ${rFrom} ${rTo} {Handler(publisher, prefix)\n} ;
}

append {\n} ;

// --- Create micro service ---
append {\tsrv, err := micro.AddService(conn.Conn(), micro.Config{\n} ;
append {\t\tName:    "} ${svcName} {",\n} ;
append {\t\tVersion: "0.1.0",\n} ;
append {\t})\n} ;
append {\tif err != nil {\n} ;
append {\t\tlog.Fatalf("failed to create micro service: %v", err)\n} ;
append {\t}\n} ;
append {\n} ;
append {\troot := srv.AddGroup(prefix)\n} ;
append {\n} ;

// --- Entity endpoints ---
foreach entity in node.entities {
string eName = entity.name;
string eVar = entity.name.toLowerCaseFirst();
string eLower = entity.name.toLowerCase();
append {\t// } ${eName} { endpoints\n} ;
append {\t} ${eVar} {Group := root.AddGroup("} ${eLower} {")\n} ;
foreach op in entity.operations {
string epName = op.capitalizedName();
string epKind = op.entityOperation.name;
append {\t} ${eVar} {Group.AddEndpoint("} ${epKind} {", micro.HandlerFunc(func(mr micro.Request) { } ${eVar} {Handler.Handle} ${epName} {(adaptRequest(mr)) }))\n} ;
}
append {\n} ;
}

// --- Relation endpoints ---
foreach relation in node.relations {
string rFrom = relation.from.name;
string rTo = relation.to.name;
string rVar = relation.from.name.toLowerCaseFirst();
string rFromLower = relation.from.name.toLowerCase();
string rToLower = relation.to.name.toLowerCase();
append {\t// } ${rFrom} { -> } ${rTo} { endpoints\n} ;
append {\t} ${rVar} ${rTo} {Group := root.AddGroup("} ${rFromLower} {").AddGroup("} ${rToLower} {")\n} ;
foreach op in relation.operations {
string rpName = op.capitalizedName();
string rpKind = op.relationOperation.name;
append {\t} ${rVar} ${rTo} {Group.AddEndpoint("} ${rpKind} {", micro.HandlerFunc(func(mr micro.Request) { } ${rVar} ${rTo} {Handler.Handle} ${rpName} {(adaptRequest(mr)) }))\n} ;
}
append {\n} ;
}

// --- Listen + graceful shutdown ---
append {\tlog.Println("} ${svcName} { listening on all subjects")\n} ;
append {\n} ;
append {\tsigCh := make(chan os.Signal, 1)\n} ;
append {\tsignal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)\n} ;
append {\t<-sigCh\n} ;
append {\n} ;
append {\tlog.Println("} ${svcName} { shutting down...")\n} ;
append {\tif err := srv.Stop(); err != nil {\n} ;
append {\t\tlog.Printf("error stopping service: %v", err)\n} ;
append {\t}\n} ;
append {\tlog.Println("} ${svcName} { stopped.")\n} ;
append {}\n} ;

}
}
```
