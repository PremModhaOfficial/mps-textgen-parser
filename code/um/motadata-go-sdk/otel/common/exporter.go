package common

import (
	"fmt"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

/* ---------------------------------------- Type Definitions -------------------------------------------------- */

// ExporterProtocol represents the OTEL exporter protocol type.
type ExporterProtocol string

/* ---------------------------------------- Constants -------------------------------------------------- */

const (
	// ProtocolGRPC represents the gRPC protocol for OTEL exporters.
	ProtocolGRPC ExporterProtocol = "grpc"

	// ProtocolHTTP represents the HTTP protocol for OTEL exporters.
	ProtocolHTTP ExporterProtocol = "http"
)

/* ---------------------------------------- Protocol Resolution -------------------------------------------------- */

// ResolveProtocol normalizes the protocol string and returns the resolved protocol.
// Returns ProtocolGRPC as default if protocol is empty.
// Supported values: "grpc", "http", "http/protobuf"
func ResolveProtocol(protocol string) (ExporterProtocol, error) {

	normalizedProtocol := strings.ToLower(protocol)

	if normalizedProtocol == "" {

		return ProtocolGRPC, nil
	}

	switch normalizedProtocol {

	case "grpc":

		return ProtocolGRPC, nil

	case "http", "http/protobuf":

		return ProtocolHTTP, nil

	default:

		return "", fmt.Errorf("unsupported OTEL protocol: %s (use 'grpc' or 'http')", protocol)
	}
}

/* ---------------------------------------- GRPC Configuration -------------------------------------------------- */

// GRPCDialOptions returns gRPC dial options based on the insecure flag.
// When insecure is true, returns options with insecure transport credentials.
func GRPCDialOptions(isInsecure bool) []grpc.DialOption {

	if isInsecure {

		return []grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		}
	}

	return nil
}

// GRPCDialOption returns a single gRPC dial option for insecure transport.
// Returns nil if isInsecure is false.
func GRPCDialOption(isInsecure bool) grpc.DialOption {

	if isInsecure {

		return grpc.WithTransportCredentials(insecure.NewCredentials())
	}

	return nil
}
