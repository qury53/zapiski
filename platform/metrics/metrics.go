package prometheus_metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

type registryConfig struct {
	goCollector      bool
	processCollector bool
	extraCollectors  []prometheus.Collector
}

type Option func(*registryConfig)

// WithGoCollector включает/выключает метрики Go runtime. По дефолту включены
func WithGoCollector(enabled bool) Option {
	return func(rc *registryConfig) { rc.goCollector = enabled }
}

// WithProcessCollector включает/выключает метрики процесса. По дефолту включены
func WithProcessCollector(enabled bool) Option {
	return func(rc *registryConfig) { rc.processCollector = enabled }
}

// WithCollectors добавляет коллекторы
func WithCollectors(collectors ...prometheus.Collector) Option {
	return func(rc *registryConfig) { rc.extraCollectors = append(rc.extraCollectors, collectors...) }
}

func NewRegistry(opts ...Option) *prometheus.Registry {
	config := &registryConfig{
		goCollector:      true,
		processCollector: true,
	}

	for _, opt := range opts {
		opt(config)
	}

	reg := prometheus.NewRegistry()

	if config.goCollector {
		reg.MustRegister(collectors.NewGoCollector())
	}
	if config.processCollector {
		reg.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	}
	if len(config.extraCollectors) > 0 {
		reg.MustRegister(config.extraCollectors...)
	}

	return reg
}
