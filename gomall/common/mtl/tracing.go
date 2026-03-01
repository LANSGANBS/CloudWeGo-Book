// Copyright 2024 CloudWeGo Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mtl

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/server"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

var TracerProvider *tracesdk.TracerProvider

func InitTracing(serviceName string) {
	otelEndpoint := os.Getenv("OTEL_EXPORTER_OTLP_TRACES_ENDPOINT")
	if otelEndpoint == "" {
		klog.Infof("OTEL_EXPORTER_OTLP_TRACES_ENDPOINT not set, skipping tracing initialization")
		TracerProvider = tracesdk.NewTracerProvider()
		otel.SetTracerProvider(TracerProvider)
		return
	}

	otelEndpoint = strings.TrimPrefix(otelEndpoint, "http://")
	otelEndpoint = strings.TrimPrefix(otelEndpoint, "https://")

	klog.Infof("Initializing tracing with endpoint: %s", otelEndpoint)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client := otlptracegrpc.NewClient(
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(otelEndpoint),
	)

	exporter, err := otlptrace.New(ctx, client)
	if err != nil {
		klog.Warnf("Failed to initialize OTLP trace exporter: %v, tracing will be disabled", err)
		TracerProvider = tracesdk.NewTracerProvider()
		otel.SetTracerProvider(TracerProvider)
		return
	}

	server.RegisterShutdownHook(func() {
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()
		_ = exporter.Shutdown(shutdownCtx)
	})

	processor := tracesdk.NewBatchSpanProcessor(exporter)
	res, err := resource.New(context.Background(), resource.WithAttributes(semconv.ServiceNameKey.String(serviceName)))
	if err != nil {
		res = resource.Default()
	}
	TracerProvider = tracesdk.NewTracerProvider(tracesdk.WithSpanProcessor(processor), tracesdk.WithResource(res))
	otel.SetTracerProvider(TracerProvider)

	klog.Infof("Tracing initialized successfully")
}
