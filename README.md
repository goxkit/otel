# Goxkit OpenTelemetry

<p align="center">
  <a href="https://github.com/goxkit/otel/blob/main/LICENSE">
    <img src="https://img.shields.io/badge/License-MIT-blue.svg" alt="License">
  </a>
  <a href="https://pkg.go.dev/github.com/goxkit/otel">
    <img src="https://godoc.org/github.com/goxkit/otel?status.svg" alt="Go Doc">
  </a>
  <a href="https://goreportcard.com/report/github.com/goxkit/otel">
    <img src="https://goreportcard.com/badge/github.com/goxkit/otel" alt="Go Report Card">
  </a>
</p>

The `otel` package provides shared utilities for OpenTelemetry integration across the Goxkit ecosystem. It contains common functionality used by the tracing, metrics, and logging packages for working with the OpenTelemetry Protocol (OTLP).

## Features

- **OTLP gRPC Integration**: Utilities for creating optimized gRPC connections to OpenTelemetry collectors
- **Configuration Integration**: Seamless integration with the Goxkit configs package
- **Connection Resilience**: Built-in reconnection strategies, keepalive mechanisms, and backoff policies
- **Common Foundation**: Shared components for use by the specialized observability packages

## Overview

This package serves as a foundation for the observability tools in Goxkit. Rather than implementing complete functionality by itself, it provides shared utilities that are used by the more specialized packages:

- `github.com/goxkit/tracing`: For distributed tracing
- `github.com/goxkit/metrics`: For application and system metrics
- `github.com/goxkit/logging`: For structured logging

The separation allows for more focused packages while sharing common code for OpenTelemetry integration.

## Usage

### OTLP gRPC Connection

```go
package main

import (
	"github.com/goxkit/configs"
	"github.com/goxkit/otel/otlp_grpc"
)

func main() {
	// Get your application configs
	cfgs := &configs.Configs{
		OTLPConfigs: &configs.OTLPConfigs{
			Endpoint:               "localhost:4317",
			ExporterIdleTimeout:    30 * time.Second,
			ExporterKeepAliveTime:  5 * time.Second,
			ExporterKeepAliveTimeout: 1 * time.Second,
		},
	}
	
	// Create a gRPC connection to the OTLP collector
	conn, err := otlpgrpc.NewExporterGRPCClient(cfgs)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	
	// Use the connection with OpenTelemetry exporters
	// ...
}
```

## Using with ConfigsBuilder

The recommended approach is to use this package indirectly through the `configs_builder` package, which handles proper initialization of all observability components:

```go
package main

import (
	"github.com/goxkit/configs_builder"
)

func main() {
	// Create configurations with all observability components enabled
	cfgs, err := configsBuilder.NewConfigsBuilder().
		Otlp().    // Enables OpenTelemetry for tracing, metrics, and logging
		Build()
	if err != nil {
		panic(err)
	}
	
	// All OpenTelemetry components are now configured and ready to use
	// cfgs.Logger - configured logger
	// cfgs.TracerProvider - configured tracer provider
	// cfgs.MeterProvider - configured meter provider
}
```

## Configuration Options

The following configuration options are used by this package:

| Setting | Environment Variable | Description |
|---------|---------------------|-------------|
| Endpoint | `OTEL_EXPORTER_OTLP_ENDPOINT` | OTLP collector endpoint (default: `localhost:4317`) |
| ExporterIdleTimeout | `OTEL_EXPORTER_IDLE_TIMEOUT` | Maximum idle time before connection is closed |
| ExporterKeepAliveTime | `OTEL_EXPORTER_KEEPALIVE_TIME` | Interval between keepalive pings |
| ExporterKeepAliveTimeout | `OTEL_EXPORTER_KEEPALIVE_TIMEOUT` | Time to wait for keepalive ack |

## License

MIT