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

package main

import (
	"context"
	"net"
	"os"
	"time"

	"github.com/cloudwego/biz-demo/gomall/app/product/biz/dal"
	"github.com/cloudwego/biz-demo/gomall/app/product/biz/service/rag"
	"github.com/cloudwego/biz-demo/gomall/app/product/conf"
	"github.com/cloudwego/biz-demo/gomall/app/product/infra/mq"
	"github.com/cloudwego/biz-demo/gomall/common/mtl"
	"github.com/cloudwego/biz-demo/gomall/common/serversuite"
	"github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/product/productcatalogservice"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/server"
	"github.com/joho/godotenv"
	"github.com/kitex-contrib/obs-opentelemetry/provider"
	"gopkg.in/natefinch/lumberjack.v2"
)

var serviceName = conf.GetConf().Kitex.Service

func main() {
	_ = godotenv.Load()
	mtl.InitLog(&lumberjack.Logger{
		Filename:   conf.GetConf().Kitex.LogFileName,
		MaxSize:    conf.GetConf().Kitex.LogMaxSize,
		MaxBackups: conf.GetConf().Kitex.LogMaxBackups,
		MaxAge:     conf.GetConf().Kitex.LogMaxAge,
	})
	mtl.InitTracing(serviceName)
	mtl.InitMetric(serviceName, conf.GetConf().Kitex.MetricsPort, conf.GetConf().Registry.RegistryAddress[0])
	dal.Init()

	if conf.GetConf().RocketMQ.Enabled {
		klog.Info("Initializing RocketMQ...")
		done := make(chan error, 1)
		go func() {
			done <- mq.InitRocketMQ()
		}()

		select {
		case err := <-done:
			if err != nil {
				klog.Warnf("Failed to initialize RocketMQ: %v", err)
			} else {
				klog.Info("RocketMQ initialized successfully")

				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()

				consumer := mq.NewStockConsumer()
				go func() {
					if err := mq.StartConsumer(ctx, consumer.HandleMessage); err != nil {
						klog.Errorf("Failed to start RocketMQ consumer: %v", err)
					}
				}()
			}
		case <-time.After(30 * time.Second):
			klog.Warn("RocketMQ initialization timeout (30s), continuing without RocketMQ")
		}
	} else {
		klog.Info("RocketMQ is disabled, skipping initialization")
	}

	huggingfaceToken := os.Getenv("SILICONFLOW_API_TOKEN")
	if huggingfaceToken != "" {
		klog.Info("Initializing RAG service with SiliconFlow API token...")
		rag.InitRAGService(huggingfaceToken)

		klog.Info("Indexing products for RAG search...")
		if err := rag.GetRAGService().IndexProducts(context.Background()); err != nil {
			klog.Errorf("Failed to index products: %v", err)
		} else {
			klog.Info("RAG product indexing completed successfully!")
		}
	} else {
		klog.Warn("HUGGINGFACE_API_TOKEN not set, RAG search will fall back to SQL search")
	}

	opts := kitexInit()

	svr := productcatalogservice.NewServer(new(ProductCatalogServiceImpl), opts...)
	if err := svr.Run(); err != nil {
		klog.Error(err.Error())
	}

	mq.Shutdown()
}

func kitexInit() (opts []server.Option) {
	// address
	address := conf.GetConf().Kitex.Address
	addr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		panic(err)
	}
	opts = append(opts, server.WithServiceAddr(addr))

	_ = provider.NewOpenTelemetryProvider(
		provider.WithSdkTracerProvider(mtl.TracerProvider),
		provider.WithEnableMetrics(false),
	)

	opts = append(opts, server.WithSuite(serversuite.CommonServerSuite{CurrentServiceName: serviceName, RegistryAddr: conf.GetConf().Registry.RegistryAddress[0]}))
	return
}
