package nats

import (
	"context"

	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/auth"
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/config"
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/core"
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/tenant"
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/transport"
)

// TransportName is the name of this transport
const TransportName = "nats"

// Transport implements the NATS transport
type Transport struct{}

// NewTransport creates a new NATS transport
func NewTransport() *Transport {
	return &Transport{}
}

// Name returns the transport name
func (t *Transport) Name() string {
	return TransportName
}

// Connect creates a new NATS connection
func (t *Transport) Connect(ctx context.Context, cfg *config.Config, creds *auth.Credentials) (core.Connection, error) {
	conn := NewConnection("", cfg, creds)
	if err := conn.Connect(ctx); err != nil {
		return nil, err
	}
	return conn, nil
}

// Publisher creates a publisher for a connection
func (t *Transport) Publisher(conn core.Connection) (core.Publisher, error) {
	natsConn, ok := conn.(*Connection)
	if !ok {
		return nil, core.NewError("transport", "publisher", core.ErrNotConnected)
	}
	return NewPublisher(natsConn), nil
}

// Subscriber creates a subscriber for a connection
func (t *Transport) Subscriber(conn core.Connection) (core.Subscriber, error) {
	natsConn, ok := conn.(*Connection)
	if !ok {
		return nil, core.NewError("transport", "subscriber", core.ErrNotConnected)
	}
	return NewSubscriber(natsConn), nil
}

// Register registers the NATS transport with the default registry
func Register() {
	transport.Register(NewTransport())
}

// TenantConnection wraps Connection to implement tenant.Connection
type TenantConnection struct {
	*Connection
	publisher  *Publisher
	subscriber *Subscriber
}

// NewTenantConnection creates a tenant connection
func NewTenantConnection(tenantID string, cfg *config.Config, creds *auth.Credentials) *TenantConnection {
	conn := NewConnection(tenantID, cfg, creds)
	return &TenantConnection{
		Connection: conn,
	}
}

// Publisher returns the publisher
func (tc *TenantConnection) Publisher() core.Publisher {
	if tc.publisher == nil {
		tc.publisher = NewPublisher(tc.Connection)
	}
	return tc.publisher
}

// Subscriber returns the subscriber
func (tc *TenantConnection) Subscriber() core.Subscriber {
	if tc.subscriber == nil {
		tc.subscriber = NewSubscriber(tc.Connection)
	}
	return tc.subscriber
}

// Close closes the tenant connection
func (tenantConnection *TenantConnection) Close(ctx context.Context) error {

	var errs []error

	if tenantConnection.publisher != nil {

		if err := tenantConnection.publisher.Close(ctx); err != nil {

			errs = append(errs, err)
		}
	}

	if tenantConnection.subscriber != nil {

		if err := tenantConnection.subscriber.Close(ctx); err != nil {

			errs = append(errs, err)
		}
	}

	if err := tenantConnection.Connection.Close(ctx); err != nil {

		errs = append(errs, err)
	}

	if len(errs) == 0 {

		return nil
	}

	if len(errs) == 1 {

		return errs[0]
	}

	return &core.MultiError{Errors: errs}
}

// Factory implements tenant.ConnectionFactory for NATS
type Factory struct{}

// NewFactory creates a new NATS connection factory
func NewFactory() *Factory {
	return &Factory{}
}

// Create creates a new tenant connection
func (f *Factory) Create(ctx context.Context, tenantID string, cfg *config.Config, creds *auth.Credentials) (tenant.Connection, error) {
	return NewTenantConnection(tenantID, cfg, creds), nil
}
