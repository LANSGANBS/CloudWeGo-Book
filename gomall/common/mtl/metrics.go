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
	"net"
	"net/http"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/registry"
	"github.com/cloudwego/kitex/server"
	consul "github.com/kitex-contrib/registry-consul"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var Registry *prometheus.Registry

func InitMetric(serviceName string, metricsPort string, registryAddr string) {
	klog.Infof("InitMetric: serviceName=%s, metricsPort=%s, registryAddr=%s", serviceName, metricsPort, registryAddr)

	Registry = prometheus.NewRegistry()
	Registry.MustRegister(collectors.NewGoCollector())
	Registry.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))

	go func() {
		klog.Infof("InitMetric: Creating consul register...")
		r, err := consul.NewConsulRegister(registryAddr)
		if err != nil {
			klog.Warnf("Failed to create consul register for metrics: %v", err)
			return
		}

		klog.Infof("InitMetric: Consul register created, resolving address...")
		addr, err := net.ResolveTCPAddr("tcp", metricsPort)
		if err != nil {
			klog.Warnf("Failed to resolve metrics address: %v", err)
			return
		}

		registryInfo := &registry.Info{
			ServiceName: "prometheus",
			Addr:        addr,
			Weight:      1,
			Tags:        map[string]string{"service": serviceName},
		}

		klog.Infof("InitMetric: Registering to consul...")
		if err := r.Register(registryInfo); err != nil {
			klog.Warnf("Failed to register prometheus to consul: %v", err)
		} else {
			klog.Infof("InitMetric: Registered to consul successfully")
		}

		server.RegisterShutdownHook(func() {
			r.Deregister(registryInfo) //nolint:errcheck
		})
	}()

	klog.Infof("InitMetric: Starting metrics HTTP server on %s", metricsPort)
	http.Handle("/metrics", promhttp.HandlerFor(Registry, promhttp.HandlerOpts{}))
	go http.ListenAndServe(metricsPort, nil) //nolint:errcheck
	klog.Infof("InitMetric: Completed")
}
