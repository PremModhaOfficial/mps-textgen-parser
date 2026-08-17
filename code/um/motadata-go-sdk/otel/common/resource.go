package common

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

/* ---------------------------------------- Type Definitions -------------------------------------------------- */

// ServiceInfo contains service identification information for OTEL resources.
type ServiceInfo struct {
	Name string

	Version string

	Environment string
}

/* ---------------------------------------- Resource Factory Functions -------------------------------------------------- */

// NewOTELResource creates a new OTEL resource with service identification attributes.
// This is used consistently across logger, metrics, and tracer providers.
func NewOTELResource(serviceInfo ServiceInfo) (*resource.Resource, error) {

	backgroundContext := context.Background()

	otelResource, resourceError := resource.New(backgroundContext,
		resource.WithAttributes(
			semconv.ServiceName(serviceInfo.Name),
			semconv.ServiceVersion(serviceInfo.Version),
			semconv.DeploymentEnvironment(serviceInfo.Environment),
		),
	)

	if resourceError != nil {

		return nil, fmt.Errorf("create OTEL resource: %w", resourceError)
	}

	return otelResource, nil
}

// NewOTELResourceFromConfig creates an OTEL resource from any config that has service info.
// This is a convenience function for configs that embed ServiceInfo fields.
func NewOTELResourceFromConfig(serviceName, serviceVersion, environment string) (*resource.Resource, error) {

	return NewOTELResource(ServiceInfo{
		Name:        serviceName,
		Version:     serviceVersion,
		Environment: environment,
	})
}
