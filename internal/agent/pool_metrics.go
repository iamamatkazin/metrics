package agent

import (
	"math/rand/v2"
	"runtime"

	"github.com/iamamatkazin/metrics.git/internal/model"
)

func (a *Agent) poolMetrics(pollCount int) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	a.metrics[model.Gauge]["Alloc"] = float64(memStats.Alloc)
	a.metrics[model.Gauge]["TotalAlloc"] = float64(memStats.TotalAlloc)
	a.metrics[model.Gauge]["BuckHashSys"] = float64(memStats.BuckHashSys)
	a.metrics[model.Gauge]["Frees"] = float64(memStats.Frees)
	a.metrics[model.Gauge]["GCCPUFraction"] = float64(memStats.GCCPUFraction)
	a.metrics[model.Gauge]["GCSys"] = float64(memStats.GCSys)
	a.metrics[model.Gauge]["HeapAlloc"] = float64(memStats.HeapAlloc)
	a.metrics[model.Gauge]["HeapIdle"] = float64(memStats.HeapIdle)
	a.metrics[model.Gauge]["HeapInuse"] = float64(memStats.HeapInuse)
	a.metrics[model.Gauge]["HeapObjects"] = float64(memStats.HeapObjects)
	a.metrics[model.Gauge]["HeapReleased"] = float64(memStats.HeapReleased)
	a.metrics[model.Gauge]["HeapSys"] = float64(memStats.HeapSys)
	a.metrics[model.Gauge]["LastGC"] = float64(memStats.LastGC)
	a.metrics[model.Gauge]["Lookups"] = float64(memStats.Lookups)
	a.metrics[model.Gauge]["MCacheInuse"] = float64(memStats.MCacheInuse)
	a.metrics[model.Gauge]["MCacheSys"] = float64(memStats.MCacheSys)
	a.metrics[model.Gauge]["MSpanInuse"] = float64(memStats.MSpanInuse)
	a.metrics[model.Gauge]["MSpanSys"] = float64(memStats.MSpanSys)
	a.metrics[model.Gauge]["Mallocs"] = float64(memStats.Mallocs)
	a.metrics[model.Gauge]["NextGC"] = float64(memStats.NextGC)
	a.metrics[model.Gauge]["NumForcedGC"] = float64(memStats.NumForcedGC)
	a.metrics[model.Gauge]["NumGC"] = float64(memStats.NumGC)
	a.metrics[model.Gauge]["OtherSys"] = float64(memStats.OtherSys)
	a.metrics[model.Gauge]["PauseTotalNs"] = float64(memStats.PauseTotalNs)
	a.metrics[model.Gauge]["StackInuse"] = float64(memStats.StackInuse)
	a.metrics[model.Gauge]["StackSys"] = float64(memStats.StackSys)
	a.metrics[model.Gauge]["Sys"] = float64(memStats.Sys)
	a.metrics[model.Counter]["PollCount"] = float64(pollCount)
	a.metrics[model.Gauge]["RandomValue"] = rand.Float64()
}
