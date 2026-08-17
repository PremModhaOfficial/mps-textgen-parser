package common

import (
	"fmt"
	"strings"

	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/utils"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

/* ---------------------------------------- Type Definitions -------------------------------------------------- */

// ExporterProtocol represents the OTEL exporter protocol type.
// OpenTelemetry supports multiple transport protocols for exporting telemetry data
// to backend collectors. This type provides type-safe protocol identification.
//
// Supported protocols:
//   - gRPC: High-performance binary protocol, recommended for high-volume telemetry
//   - HTTP: RESTful protocol with protobuf encoding, useful for environments where gRPC is not available
type ExporterProtocol string

/* ---------------------------------------- Constants -------------------------------------------------- */

const (
	// ProtocolGRPC represents the gRPC protocol for OTEL exporters.
	// gRPC offers better performance and lower latency compared to HTTP,
	// making it the preferred choice for high-throughput applications.
	// Uses HTTP/2 with bidirectional streaming and binary protobuf encoding.
	ProtocolGRPC ExporterProtocol = utils.OTELProtocolGRPC

	// ProtocolHTTP represents the HTTP protocol for OTEL exporters.
	// HTTP uses standard REST conventions with protobuf encoding.
	// Useful when gRPC is blocked by firewalls or load balancers.
	ProtocolHTTP ExporterProtocol = utils.OTELProtocolHTTP
)

/* ---------------------------------------- Protocol Resolution -------------------------------------------------- */

// ResolveProtocol normalizes the protocol string and returns the resolved protocol.
// This function provides flexible protocol specification while ensuring type safety.
//
// Protocol resolution rules:
//   - Empty string defaults to gRPC (recommended protocol)
//   - "grpc" -> ProtocolGRPC
//   - "http" or "http/protobuf" -> ProtocolHTTP
//   - Case-insensitive matching is applied
//
// Parameters:
//   - protocol: The protocol string to resolve (can be empty for default)
//
// Returns:
//   - ExporterProtocol: The resolved protocol constant
//   - error: Non-nil if the protocol is not supported
//
// Example:
//
//	protocol, err := ResolveProtocol("grpc")   // Returns ProtocolGRPC, nil
//	protocol, err := ResolveProtocol("")      // Returns ProtocolGRPC, nil (default)
//	protocol, err := ResolveProtocol("HTTP")  // Returns ProtocolHTTP, nil (case-insensitive)
//	protocol, err := ResolveProtocol("tcp")   // Returns "", error
func ResolveProtocol(protocol string) (ExporterProtocol, error) {

	// Normalize to lowercase for case-insensitive matching
	// This allows users to specify "GRPC", "Grpc", or "grpc" interchangeably
	normalizedProtocol := strings.ToLower(protocol)

	// Default to gRPC when no protocol is specified
	// gRPC is the recommended protocol for its performance characteristics
	if normalizedProtocol == "" {

		return ProtocolGRPC, nil
	}

	// Map normalized protocol string to typed constant
	switch normalizedProtocol {

	case utils.OTELProtocolGRPC:
		// gRPC protocol - high performance binary protocol over HTTP/2
		return ProtocolGRPC, nil

	case utils.OTELProtocolHTTP, utils.OTELProtocolHTTPProtobuf:
		// HTTP protocol - REST with protobuf encoding
		// "http/protobuf" is an alias commonly used in OTEL configurations
		return ProtocolHTTP, nil

	default:
		// Unsupported protocol - return error with helpful message
		return "", fmt.Errorf("%w: %s", utils.ErrOTELUnsupportedProtocol, protocol)
	}
}

/* ---------------------------------------- gRPC Configuration -------------------------------------------------- */

// GRPCDialOptions returns gRPC dial options based on the insecure flag.
// This function configures transport credentials for gRPC connections.
//
// Security considerations:
//   - In production, always use secure TLS connections (isInsecure=false)
//   - Insecure mode should only be used for local development or testing
//   - When isInsecure=true, data is transmitted in plaintext
//
// Parameters:
//   - isInsecure: If true, disables TLS for the gRPC connection
//
// Returns:
//   - []grpc.DialOption: Slice of dial options (may be empty if secure mode)
//
// Example:
//
//	// For development with local collector
//	opts := GRPCDialOptions(true)
//
//	// For production with TLS
//	opts := GRPCDialOptions(false)
func GRPCDialOptions(isInsecure bool) []grpc.DialOption {

	// Only return insecure credentials when explicitly requested
	// This prevents accidental plaintext transmission in production
	if isInsecure {

		return []grpc.DialOption{
			// NewCredentials returns credentials that disable transport security
			// WARNING: Only use for local development or testing
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		}
	}

	// Return nil to use default secure transport (TLS)
	// The gRPC client will automatically negotiate TLS
	return nil
}

// GRPCDialOption returns a single gRPC dial option for insecure transport.
// This is a convenience function when only the credentials option is needed.
//
// This function is useful when the OTEL exporter API expects individual options
// rather than a slice of options.
//
// Parameters:
//   - isInsecure: If true, returns insecure transport credentials
//
// Returns:
//   - grpc.DialOption: The dial option, or nil if secure mode
//
// Example:
//
//	if opt := GRPCDialOption(cfg.OTELInsecure); opt != nil {
//	    exporterOptions = append(exporterOptions, otlpgrpc.WithDialOption(opt))
//	}
func GRPCDialOption(isInsecure bool) grpc.DialOption {

	// Return insecure credentials only when explicitly requested
	if isInsecure {

		return grpc.WithTransportCredentials(insecure.NewCredentials())
	}

	// Return nil to indicate no additional option needed
	// The caller should check for nil before using
	return nil
}
