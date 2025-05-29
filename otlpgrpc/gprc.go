// Copyright (c) 2025, The GoKit Authors
// MIT License
// All rights reserved.

// Package otlpgrpc provides gRPC client connection utilities for OpenTelemetry OTLP exporters.
// It handles configuration and creation of properly configured gRPC connections to OTLP collectors,
// including automatic reconnection, keepalive parameters, and backoff strategies.
package otlpgrpc

import (
	"context"
	"crypto/x509"
	"fmt"
	"strings"
	"time"

	"github.com/goxkit/configs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

// NewExporterGRPCClient creates a new gRPC client connection for OpenTelemetry OTLP exporters
// with configurations optimized for telemetry data export. The connection is configured with:
//   - Insecure credentials (for non-TLS connections)
//   - Idle timeout from configuration
//   - Keepalive parameters for maintaining long-lived connections
//   - Exponential backoff strategy for reconnection attempts
//
// Parameters:
//   - cfgs: Application configurations containing OTLP settings
//
// Returns:
//   - *grpc.ClientConn: The configured gRPC client connection
//   - error: Any error encountered during connection setup
func NewExporterGRPCClient(cfgs *configs.Configs) (*grpc.ClientConn, error) {

	conn, err := grpc.NewClient(
		cfgs.OTLPConfigs.Endpoint,
		grpc.WithTransportCredentials(evaluateCredentials(cfgs)),
		grpc.WithPerRPCCredentials(newPerRPCCredentials(cfgs)),
		grpc.WithIdleTimeout(cfgs.OTLPConfigs.ExporterIdleTimeout),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:    cfgs.OTLPConfigs.ExporterKeepAliveTime,
			Timeout: cfgs.OTLPConfigs.ExporterKeepAliveTimeout,
		}),
		grpc.WithConnectParams(grpc.ConnectParams{
			Backoff: backoff.Config{
				BaseDelay:  1 * time.Second,
				Multiplier: 1.6,
				MaxDelay:   15 * time.Second,
			},
			MinConnectTimeout: 0,
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create otel exporter gRPC conn: %w", err)
	}

	return conn, err
}

func evaluateCredentials(cfgs *configs.Configs) credentials.TransportCredentials {
	if !cfgs.OTLPConfigs.ExporterTLSEnabled {
		return insecure.NewCredentials()
	}

	certPool := x509.NewCertPool()
	return credentials.NewClientTLSFromCert(certPool, "")
}

type perRPCCredentials struct {
	tlsEnabled bool
	headers    map[string]string
}

func newPerRPCCredentials(cfgs *configs.Configs) credentials.PerRPCCredentials {
	h := map[string]string{}

	if cfgs.OTLPConfigs.ExporterHeaders != "" {
		keyValue := strings.Split(cfgs.OTLPConfigs.ExporterHeaders, ",")
		for _, kv := range keyValue {
			parts := strings.SplitN(kv, "=", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				if key != "" {
					h[key] = value
				}
			}
		}
	}

	return &perRPCCredentials{
		tlsEnabled: cfgs.OTLPConfigs.ExporterTLSEnabled,
		headers:    h,
	}
}

func (h *perRPCCredentials) GetRequestMetadata(_ context.Context, _ ...string) (map[string]string, error) {
	return h.headers, nil
}

func (h *perRPCCredentials) RequireTransportSecurity() bool {
	return h.tlsEnabled
}
