package agent

import (
	"math/rand/v2"
	"runtime"

	"github.com/iamamatkazin/metrics.git/internal/model"
)

func (a *Agent) poolMetrics(pollCount int) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	a.metrics[model.Gauge]["Alloc"] = memStats.Alloc
	a.metrics[model.Gauge]["BuckHashSys"] = memStats.BuckHashSys
	a.metrics[model.Gauge]["Frees"] = memStats.Frees
	a.metrics[model.Gauge]["GCCPUFraction"] = memStats.GCCPUFraction
	a.metrics[model.Gauge]["GCSys"] = memStats.GCSys
	a.metrics[model.Gauge]["HeapAlloc"] = memStats.HeapAlloc
	a.metrics[model.Gauge]["HeapIdle"] = memStats.HeapIdle
	a.metrics[model.Gauge]["HeapInuse"] = memStats.HeapInuse
	a.metrics[model.Gauge]["HeapObjects"] = memStats.HeapObjects
	a.metrics[model.Gauge]["HeapReleased"] = memStats.HeapReleased
	a.metrics[model.Gauge]["HeapSys"] = memStats.HeapSys
	a.metrics[model.Gauge]["LastGC"] = memStats.LastGC
	a.metrics[model.Gauge]["Lookups"] = memStats.Lookups
	a.metrics[model.Gauge]["MCacheInuse"] = memStats.MCacheInuse
	a.metrics[model.Gauge]["MCacheSys"] = memStats.MCacheSys
	a.metrics[model.Gauge]["MSpanInuse"] = memStats.MSpanInuse
	a.metrics[model.Gauge]["MSpanSys"] = memStats.MSpanSys
	a.metrics[model.Gauge]["Mallocs"] = memStats.Mallocs
	a.metrics[model.Gauge]["NextGC"] = memStats.NextGC
	a.metrics[model.Gauge]["NumForcedGC"] = memStats.NumForcedGC
	a.metrics[model.Gauge]["NumGC"] = memStats.NumGC
	a.metrics[model.Gauge]["OtherSys"] = memStats.OtherSys
	a.metrics[model.Gauge]["PauseTotalNs"] = memStats.PauseTotalNs
	a.metrics[model.Gauge]["StackInuse"] = memStats.StackInuse
	a.metrics[model.Gauge]["StackSys"] = memStats.StackSys
	a.metrics[model.Gauge]["Sys"] = memStats.Sys
	a.metrics[model.Counter]["PollCount"] = pollCount
	a.metrics[model.Gauge]["RandomValue"] = rand.Float64()
}
