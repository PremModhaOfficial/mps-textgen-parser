// Package common provides shared utilities for OpenTelemetry components.
//
// This package contains common functionality used across logger, tracer, and metrics
// packages including:
//   - OTEL resource creation with service identification attributes
//   - Protocol resolution for OTLP exporters (gRPC/HTTP)
//   - gRPC dial option configuration for secure/insecure connections
//   - Shutdown error aggregation for graceful component shutdown
//
// The package follows the SDK coding standards and utilizes centralized error
// definitions from the utils package for consistent error handling.
//
// Usage:
//
//	// Create an OTEL resource with service info
//	resource, err := common.NewOTELResource(common.ServiceInfo{
//	    Name:        "my-service",
//	    Version:     "1.0.0",
//	    Environment: "production",
//	})
//
//	// Or use the convenience function
//	resource, err := common.NewOTELResourceFromConfig("my-service", "1.0.0", "production")
package common

import (
	"context"
	"fmt"

	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/utils"

	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

/* ---------------------------------------- Type Definitions -------------------------------------------------- */

// ServiceInfo contains service identification information for OTEL resources.
// This struct encapsulates the standard service attributes required by OpenTelemetry
// for proper service identification in observability backends.
//
// Fields:
//   - Name: The logical name of the service (e.g., "payment-service")
//   - Version: The semantic version of the service (e.g., "1.2.3")
//   - Environment: The deployment environment (e.g., "production", "staging", "development")
type ServiceInfo struct {
	// Name is the logical name of the service as it appears in observability tools
	Name string

	// Version is the semantic version of the service for tracking deployments
	Version string

	// Environment identifies the deployment environment for filtering and routing
	Environment string
}

/* ---------------------------------------- Resource Factory Functions -------------------------------------------------- */

// NewOTELResource creates a new OTEL resource with service identification attributes.
// This is used consistently across logger, metrics, and tracer providers to ensure
// uniform service identification in all telemetry data.
//
// The function creates a resource with the following semantic convention attributes:
//   - service.name: Identifies the service in observability backends
//   - service.version: Tracks the deployed version for correlation
//   - deployment.environment: Distinguishes between environments (prod/staging/dev)
//
// Parameters:
//   - serviceInfo: ServiceInfo struct containing name, version, and environment
//
// Returns:
//   - *resource.Resource: The configured OTEL resource
//   - error: Non-nil if resource creation fails
//
// Example:
//
//	resource, err := NewOTELResource(ServiceInfo{
//	    Name:        "api-gateway",
//	    Version:     "2.1.0",
//	    Environment: "production",
//	})
func NewOTELResource(serviceInfo ServiceInfo) (*resource.Resource, error) {

	// Use background context for resource creation as this is typically
	// called during application startup where no request context exists
	backgroundContext := context.Background()

	// Create the OTEL resource with standard semantic convention attributes
	// These attributes are automatically propagated to all telemetry data
	otelResource, resourceError := resource.New(backgroundContext,
		resource.WithAttributes(
			// service.name: Primary identifier for the service in observability tools
			semconv.ServiceName(serviceInfo.Name),

			// service.version: Enables version-based filtering and deployment tracking
			semconv.ServiceVersion(serviceInfo.Version),

			// deployment.environment: Critical for separating prod/staging/dev data
			semconv.DeploymentEnvironment(serviceInfo.Environment),
		),
	)

	// Wrap the error with context for better debugging
	if resourceError != nil {

		return nil, fmt.Errorf("%w: %v", utils.ErrOTELResourceCreationFailed, resourceError)
	}

	return otelResource, nil
}

// NewOTELResourceFromConfig creates an OTEL resource from individual config fields.
// This is a convenience function for configs that have service info as separate fields
// rather than embedded as a ServiceInfo struct.
//
// This function is commonly used by provider constructors that receive configuration
// with flattened service identification fields.
//
// Parameters:
//   - serviceName: The logical name of the service
//   - serviceVersion: The semantic version of the service
//   - environment: The deployment environment identifier
//
// Returns:
//   - *resource.Resource: The configured OTEL resource
//   - error: Non-nil if resource creation fails
//
// Example:
//
//	resource, err := NewOTELResourceFromConfig("user-service", "1.0.0", "staging")
func NewOTELResourceFromConfig(serviceName, serviceVersion, environment string) (*resource.Resource, error) {

	// Delegate to NewOTELResource with a constructed ServiceInfo
	// This maintains a single code path for resource creation
	return NewOTELResource(ServiceInfo{
		Name:        serviceName,
		Version:     serviceVersion,
		Environment: environment,
	})
}
