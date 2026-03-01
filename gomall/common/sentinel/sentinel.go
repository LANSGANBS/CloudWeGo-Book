package sentinel

import (
	"context"
	"log"
	"time"

	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/circuitbreaker"
	"github.com/alibaba/sentinel-golang/core/config"
	"github.com/alibaba/sentinel-golang/core/flow"
	"github.com/alibaba/sentinel-golang/core/hotspot"
	"github.com/alibaba/sentinel-golang/core/isolation"
	"github.com/alibaba/sentinel-golang/core/system"
	"github.com/alibaba/sentinel-golang/logging"
	"github.com/cloudwego/kitex/pkg/klog"
)

type Config struct {
	Enabled             bool    `yaml:"enabled"`
	AppName             string  `yaml:"app_name"`
	LogDir              string  `yaml:"log_dir"`
	CPUThreshold        float64 `yaml:"cpu_threshold"`
	MemoryThreshold     float64 `yaml:"memory_threshold"`
	LoadThreshold       float64 `yaml:"load_threshold"`
	QPSThreshold        int64   `yaml:"qps_threshold"`
	ConcurrentThreshold int     `yaml:"concurrent_threshold"`
}

func InitSentinel(cfg *Config) error {
	if !cfg.Enabled {
		klog.Info("Sentinel is disabled")
		return nil
	}

	logging.ResetGlobalLoggerLevel(logging.InfoLevel)

	conf := config.NewDefaultConfig()
	conf.Sentinel.App.Name = cfg.AppName
	conf.Sentinel.Log.Dir = cfg.LogDir
	conf.Sentinel.Log.UsePid = true

	err := sentinel.InitWithConfig(conf)
	if err != nil {
		return err
	}

	initFlowRules(cfg)
	initCircuitBreakerRules()
	initSystemRules(cfg)
	initIsolationRules(cfg)
	initHotspotRules()

	klog.Info("Sentinel initialized successfully")
	return nil
}

func initFlowRules(cfg *Config) {
	_, err := flow.LoadRules([]*flow.Rule{
		{
			Resource:               "api_list_products",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              float64(cfg.QPSThreshold),
			StatIntervalInMs:       1000,
		},
		{
			Resource:               "api_get_product",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              float64(cfg.QPSThreshold * 2),
			StatIntervalInMs:       1000,
		},
		{
			Resource:               "api_deduct_stock",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              float64(cfg.QPSThreshold / 2),
			StatIntervalInMs:       1000,
		},
		{
			Resource:               "api_checkout",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Throttling,
			Threshold:              float64(cfg.QPSThreshold / 4),
			StatIntervalInMs:       1000,
		},
	})
	if err != nil {
		klog.Errorf("Failed to load flow rules: %v", err)
	}
}

func initCircuitBreakerRules() {
	_, err := circuitbreaker.LoadRules([]*circuitbreaker.Rule{
		{
			Resource:         "rpc_product",
			Strategy:         circuitbreaker.ErrorRatio,
			RetryTimeoutMs:   3000,
			MinRequestAmount: 10,
			StatIntervalMs:   5000,
			Threshold:        0.5,
		},
		{
			Resource:         "rpc_cart",
			Strategy:         circuitbreaker.ErrorRatio,
			RetryTimeoutMs:   3000,
			MinRequestAmount: 10,
			StatIntervalMs:   5000,
			Threshold:        0.5,
		},
		{
			Resource:         "rpc_order",
			Strategy:         circuitbreaker.ErrorRatio,
			RetryTimeoutMs:   5000,
			MinRequestAmount: 5,
			StatIntervalMs:   5000,
			Threshold:        0.3,
		},
		{
			Resource:         "rpc_payment",
			Strategy:         circuitbreaker.ErrorRatio,
			RetryTimeoutMs:   10000,
			MinRequestAmount: 5,
			StatIntervalMs:   5000,
			Threshold:        0.2,
		},
		{
			Resource:         "redis_stock",
			Strategy:         circuitbreaker.ErrorRatio,
			RetryTimeoutMs:   2000,
			MinRequestAmount: 20,
			StatIntervalMs:   5000,
			Threshold:        0.3,
		},
	})
	if err != nil {
		klog.Errorf("Failed to load circuit breaker rules: %v", err)
	}
}

func initSystemRules(cfg *Config) {
	_, err := system.LoadRules([]*system.Rule{
		{
			ID:           "system_load_rule",
			MetricType:   system.Load,
			TriggerCount: cfg.LoadThreshold,
			Strategy:     system.BBR,
		},
		{
			ID:           "system_concurrency_rule",
			MetricType:   system.Concurrency,
			TriggerCount: float64(cfg.ConcurrentThreshold),
			Strategy:     system.BBR,
		},
		{
			ID:           "system_inbound_qps_rule",
			MetricType:   system.InboundQPS,
			TriggerCount: float64(cfg.QPSThreshold),
			Strategy:     system.BBR,
		},
	})
	if err != nil {
		klog.Errorf("Failed to load system rules: %v", err)
	}
}

func initIsolationRules(cfg *Config) {
	_, err := isolation.LoadRules([]*isolation.Rule{
		{
			Resource:   "db_write",
			MetricType: isolation.Concurrency,
			Threshold:  uint32(cfg.ConcurrentThreshold / 2),
		},
		{
			Resource:   "db_read",
			MetricType: isolation.Concurrency,
			Threshold:  uint32(cfg.ConcurrentThreshold),
		},
		{
			Resource:   "redis_write",
			MetricType: isolation.Concurrency,
			Threshold:  uint32(cfg.ConcurrentThreshold),
		},
		{
			Resource:   "redis_read",
			MetricType: isolation.Concurrency,
			Threshold:  uint32(cfg.ConcurrentThreshold * 2),
		},
	})
	if err != nil {
		klog.Errorf("Failed to load isolation rules: %v", err)
	}
}

func initHotspotRules() {
	_, err := hotspot.LoadRules([]*hotspot.Rule{
		{
			Resource:        "hot_product",
			MetricType:      hotspot.QPS,
			ControlBehavior: hotspot.Reject,
			ParamIndex:      0,
			Threshold:       100,
			DurationInSec:   1,
		},
		{
			Resource:        "hot_user",
			MetricType:      hotspot.QPS,
			ControlBehavior: hotspot.Throttling,
			ParamIndex:      0,
			Threshold:       50,
			DurationInSec:   1,
		},
	})
	if err != nil {
		klog.Errorf("Failed to load hotspot rules: %v", err)
	}
}

type CircuitBreakerListener struct{}

func (l *CircuitBreakerListener) OnTransformToClosed(prev circuitbreaker.State, rule circuitbreaker.Rule) {
	klog.Infof("Circuit breaker [%s] transformed to closed", rule.ResourceName())
}

func (l *CircuitBreakerListener) OnTransformToOpen(prev circuitbreaker.State, rule circuitbreaker.Rule, snapshot interface{}) {
	klog.Warnf("Circuit breaker [%s] transformed to open, snapshot: %v", rule.ResourceName(), snapshot)
}

func (l *CircuitBreakerListener) OnTransformToHalfOpen(prev circuitbreaker.State, rule circuitbreaker.Rule) {
	klog.Infof("Circuit breaker [%s] transformed to half-open", rule.ResourceName())
}

func RegisterCircuitBreakerListener() {
	circuitbreaker.RegisterStateChangeListeners(&CircuitBreakerListener{})
}

func Entry(resource string) (interface{}, error) {
	entry, blockErr := sentinel.Entry(resource)
	if blockErr != nil {
		return nil, blockErr
	}
	return entry, nil
}

func WithCircuitBreaker(resource string, fallback func() error) func() error {
	return func() error {
		entry, blockErr := sentinel.Entry(resource)
		if blockErr != nil {
			klog.Warnf("Circuit breaker [%s] is open, executing fallback", resource)
			return fallback()
		}
		defer entry.Exit()

		return nil
	}
}

func WrapWithSentinel(resource string, fn func(ctx context.Context) error) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		entry, blockErr := sentinel.Entry(resource)
		if blockErr != nil {
			klog.Warnf("Request rejected by sentinel [%s]: %v", resource, blockErr)
			return blockErr
		}
		defer entry.Exit()

		return fn(ctx)
	}
}

func WrapWithSentinelAndFallback(resource string, fn func(ctx context.Context) error, fallback func(ctx context.Context) error) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		entry, blockErr := sentinel.Entry(resource)
		if blockErr != nil {
			klog.Warnf("Request rejected by sentinel [%s], executing fallback", resource)
			return fallback(ctx)
		}
		defer entry.Exit()

		return fn(ctx)
	}
}

type AdaptiveLimiter struct {
	resource         string
	minThreshold     float64
	maxThreshold     float64
	currentThreshold float64
	adjustInterval   time.Duration
	lastAdjustTime   time.Time
}

func NewAdaptiveLimiter(resource string, minThreshold, maxThreshold int64) *AdaptiveLimiter {
	return &AdaptiveLimiter{
		resource:         resource,
		minThreshold:     float64(minThreshold),
		maxThreshold:     float64(maxThreshold),
		currentThreshold: float64(minThreshold),
		adjustInterval:   10 * time.Second,
		lastAdjustTime:   time.Now(),
	}
}

func (l *AdaptiveLimiter) Adjust(cpuUsage float64, avgRT float64, errorRate float64) {
	now := time.Now()
	if now.Sub(l.lastAdjustTime) < l.adjustInterval {
		return
	}
	l.lastAdjustTime = now

	adjustFactor := 1.0

	if cpuUsage > 0.8 {
		adjustFactor *= 0.7
	} else if cpuUsage > 0.6 {
		adjustFactor *= 0.85
	} else if cpuUsage < 0.3 {
		adjustFactor *= 1.2
	}

	if avgRT > 500 {
		adjustFactor *= 0.8
	} else if avgRT > 200 {
		adjustFactor *= 0.9
	} else if avgRT < 50 {
		adjustFactor *= 1.1
	}

	if errorRate > 0.1 {
		adjustFactor *= 0.5
	} else if errorRate > 0.05 {
		adjustFactor *= 0.7
	}

	newThreshold := l.currentThreshold * adjustFactor

	if newThreshold < l.minThreshold {
		newThreshold = l.minThreshold
	} else if newThreshold > l.maxThreshold {
		newThreshold = l.maxThreshold
	}

	if newThreshold != l.currentThreshold {
		l.updateFlowRule(newThreshold)
		l.currentThreshold = newThreshold
		klog.Infof("Adaptive limiter [%s] adjusted threshold from %.0f to %.0f",
			l.resource, l.currentThreshold, newThreshold)
	}
}

func (l *AdaptiveLimiter) updateFlowRule(threshold float64) {
	_, err := flow.LoadRules([]*flow.Rule{
		{
			Resource:               l.resource,
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              threshold,
			StatIntervalInMs:       1000,
		},
	})
	if err != nil {
		log.Printf("Failed to update flow rule: %v", err)
	}
}

func (l *AdaptiveLimiter) GetCurrentThreshold() float64 {
	return l.currentThreshold
}

type SelfHealer struct {
	healthCheckInterval time.Duration
	healthChecks        map[string]func() bool
	recoveryActions     map[string]func() error
}

func NewSelfHealer() *SelfHealer {
	return &SelfHealer{
		healthCheckInterval: 30 * time.Second,
		healthChecks:        make(map[string]func() bool),
		recoveryActions:     make(map[string]func() error),
	}
}

func (h *SelfHealer) RegisterHealthCheck(name string, check func() bool, recovery func() error) {
	h.healthChecks[name] = check
	h.recoveryActions[name] = recovery
}

func (h *SelfHealer) Start(ctx context.Context) {
	ticker := time.NewTicker(h.healthCheckInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			h.performHealthChecks()
		}
	}
}

func (h *SelfHealer) performHealthChecks() {
	for name, check := range h.healthChecks {
		if !check() {
			klog.Warnf("Health check failed for [%s], attempting recovery", name)
			if recovery, ok := h.recoveryActions[name]; ok {
				if err := recovery(); err != nil {
					klog.Errorf("Recovery failed for [%s]: %v", name, err)
				} else {
					klog.Infof("Recovery successful for [%s]", name)
				}
			}
		}
	}
}
