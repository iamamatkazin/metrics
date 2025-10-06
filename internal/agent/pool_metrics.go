package agent

import (
	"math/rand/v2"
	"runtime"

	"github.com/iamamatkazin/metrics.git/internal/model"
)

func (a *Agent) poolMetrics(pollCount int) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	a.metrics[models.Gauge]["Alloc"] = memStats.Alloc
	a.metrics[models.Gauge]["BuckHashSys"] = memStats.BuckHashSys
	a.metrics[models.Gauge]["Frees"] = memStats.Frees
	a.metrics[models.Gauge]["GCCPUFraction"] = memStats.GCCPUFraction
	a.metrics[models.Gauge]["GCSys"] = memStats.GCSys
	a.metrics[models.Gauge]["HeapAlloc"] = memStats.HeapAlloc
	a.metrics[models.Gauge]["HeapIdle"] = memStats.HeapIdle
	a.metrics[models.Gauge]["HeapInuse"] = memStats.HeapInuse
	a.metrics[models.Gauge]["HeapObjects"] = memStats.HeapObjects
	a.metrics[models.Gauge]["HeapReleased"] = memStats.HeapReleased
	a.metrics[models.Gauge]["HeapSys"] = memStats.HeapSys
	a.metrics[models.Gauge]["LastGC"] = memStats.LastGC
	a.metrics[models.Gauge]["Lookups"] = memStats.Lookups
	a.metrics[models.Gauge]["MCacheInuse"] = memStats.MCacheInuse
	a.metrics[models.Gauge]["MCacheSys"] = memStats.MCacheSys
	a.metrics[models.Gauge]["MSpanInuse"] = memStats.MSpanInuse
	a.metrics[models.Gauge]["MSpanSys"] = memStats.MSpanSys
	a.metrics[models.Gauge]["Mallocs"] = memStats.Mallocs
	a.metrics[models.Gauge]["NextGC"] = memStats.NextGC
	a.metrics[models.Gauge]["NumForcedGC"] = memStats.NumForcedGC
	a.metrics[models.Gauge]["NumGC"] = memStats.NumGC
	a.metrics[models.Gauge]["OtherSys"] = memStats.OtherSys
	a.metrics[models.Gauge]["PauseTotalNs"] = memStats.PauseTotalNs
	a.metrics[models.Gauge]["StackInuse"] = memStats.StackInuse
	a.metrics[models.Gauge]["StackSys"] = memStats.StackSys
	a.metrics[models.Gauge]["Sys"] = memStats.Sys
	a.metrics[models.Counter]["PollCount"] = pollCount
	a.metrics[models.Gauge]["RandomValue"] = rand.Float64()
}
