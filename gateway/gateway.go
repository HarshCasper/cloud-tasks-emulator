// Package gateway implements a REST gateway for the Cloud Tasks emulator gRPC service.
// It uses grpc-gateway to translate RESTful HTTP/JSON API calls into gRPC.
package gateway

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	// Import the generated gateway package (standalone mode)
	cloudtasksgw "github.com/aertje/cloud-tasks-emulator/gateway/cloudtasksgw/google/cloud/tasks/v2"
)

// Options encapsulates options for the REST gateway.
type Options struct {
	// GatewayAddress is the address for the REST gateway to listen on (e.g., "localhost:8080")
	GatewayAddress string
	// GRPCAddress is the address of the gRPC server to proxy to (e.g., "localhost:8123")
	GRPCAddress string
}

// Gateway implements the REST gateway server for Cloud Tasks emulator.
type Gateway struct {
	opts   Options
	server *http.Server
}

// New returns a new Gateway with the given options.
func New(opts Options) *Gateway {
	return &Gateway{opts: opts}
}

// Run starts the REST gateway server. It blocks until the server is stopped.
func (gw *Gateway) Run(ctx context.Context) error {
	// Create the gRPC-Gateway mux with JSON marshaling options
	mux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{}),
	)

	// Set up gRPC dial options (insecure for local development)
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	// Register the CloudTasks handler from the gRPC endpoint
	err := cloudtasksgw.RegisterCloudTasksHandlerFromEndpoint(ctx, mux, gw.opts.GRPCAddress, opts)
	if err != nil {
		return fmt.Errorf("failed to register CloudTasks handler: %w", err)
	}

	// Create the HTTP server
	gw.server = &http.Server{
		Addr:    gw.opts.GatewayAddress,
		Handler: mux,
	}

	log.Printf("REST gateway listening at %s (proxying to gRPC at %s)\n", gw.opts.GatewayAddress, gw.opts.GRPCAddress)

	// Start serving
	if err := gw.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("REST gateway server error: %w", err)
	}

	return nil
}

// Shutdown gracefully shuts down the gateway server.
func (gw *Gateway) Shutdown(ctx context.Context) error {
	if gw.server != nil {
		return gw.server.Shutdown(ctx)
	}
	return nil
}
