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

	"github.com/cloudwego/biz-demo/gomall/app/frontend/utils"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/route"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

var TracerProvider *tracesdk.TracerProvider

func InitTracing() route.CtxCallback {
	otelEndpoint := os.Getenv("OTEL_EXPORTER_OTLP_TRACES_ENDPOINT")
	if otelEndpoint == "" {
		hlog.Infof("OTEL_EXPORTER_OTLP_TRACES_ENDPOINT not set, skipping tracing initialization")
		TracerProvider = tracesdk.NewTracerProvider()
		otel.SetTracerProvider(TracerProvider)
		return route.CtxCallback(func(ctx context.Context) {})
	}

	otelEndpoint = strings.TrimPrefix(otelEndpoint, "http://")
	otelEndpoint = strings.TrimPrefix(otelEndpoint, "https://")

	hlog.Infof("Initializing tracing with endpoint: %s", otelEndpoint)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	exporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(otelEndpoint),
	)
	if err != nil {
		hlog.Warnf("Failed to initialize OTLP trace exporter: %v, tracing will be disabled", err)
		TracerProvider = tracesdk.NewTracerProvider()
		otel.SetTracerProvider(TracerProvider)
		return route.CtxCallback(func(ctx context.Context) {})
	}

	processor := tracesdk.NewBatchSpanProcessor(exporter)
	res, err := resource.New(context.Background(), resource.WithAttributes(semconv.ServiceNameKey.String(utils.ServiceName)))
	if err != nil {
		res = resource.Default()
	}
	TracerProvider = tracesdk.NewTracerProvider(tracesdk.WithSpanProcessor(processor), tracesdk.WithResource(res))
	otel.SetTracerProvider(TracerProvider)

	hlog.Infof("Tracing initialized successfully")

	return route.CtxCallback(func(ctx context.Context) {
		exporter.Shutdown(ctx) //nolint:errcheck
	})
}
