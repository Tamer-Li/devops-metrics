package service

import "runtime"

type MemStatsRuntime struct {
	Alloc         float64
	BuckHashSys   float64
	Frees         float64
	GCCPUFraction float64
	GCSys         float64
	HeapAlloc     float64
	HeapIdle      float64
	HeapInuse     float64
	HeapObjects   float64
	HeapReleased  float64
	HeapSys       float64
	LastGC        float64
	Lookups       float64
	MCacheInuse   float64
	MCacheSys     float64
	MSpanInuse    float64
	MSpanSys      float64
	Mallocs       float64
	NextGC        float64
	NumForcedGC   float64
	NumGC         float64
	OtherSys      float64
	PauseTotalNs  float64
	StackInuse    float64
	StackSys      float64
	Sys           float64
	TotalAlloc    float64
}

func NewMemStatsRuntime() *MemStatsRuntime {
	m := &MemStatsRuntime{}
	m.Update()
	return m
}

func (m *MemStatsRuntime) Update() {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	m.Alloc = float64(mem.Alloc)
	m.BuckHashSys = float64(mem.BuckHashSys)
	m.Frees = float64(mem.Frees)
	m.GCCPUFraction = mem.GCCPUFraction
	m.GCSys = float64(mem.GCSys)
	m.HeapAlloc = float64(mem.HeapAlloc)
	m.HeapIdle = float64(mem.HeapIdle)
	m.HeapInuse = float64(mem.HeapInuse)
	m.HeapObjects = float64(mem.HeapObjects)
	m.HeapReleased = float64(mem.HeapReleased)
	m.HeapSys = float64(mem.HeapSys)
	m.LastGC = float64(mem.LastGC)
	m.Lookups = float64(mem.Lookups)
	m.MCacheInuse = float64(mem.MCacheInuse)
	m.MCacheSys = float64(mem.MCacheSys)
	m.MSpanInuse = float64(mem.MSpanInuse)
	m.MSpanSys = float64(mem.MSpanSys)
	m.Mallocs = float64(mem.Mallocs)
	m.NextGC = float64(mem.NextGC)
	m.NumForcedGC = float64(mem.NumForcedGC)
	m.NumGC = float64(mem.NumGC)
	m.OtherSys = float64(mem.OtherSys)
	m.PauseTotalNs = float64(mem.PauseTotalNs)
	m.StackInuse = float64(mem.StackInuse)
	m.StackSys = float64(mem.StackSys)
	m.Sys = float64(mem.Sys)
	m.TotalAlloc = float64(mem.TotalAlloc)
}

func (m *MemStatsRuntime) GetAll() map[string]float64 {
	return map[string]float64{
		"Alloc":         m.Alloc,
		"BuckHashSys":   m.BuckHashSys,
		"Frees":         m.Frees,
		"GCCPUFraction": m.GCCPUFraction,
		"GCSys":         m.GCSys,
		"HeapAlloc":     m.HeapAlloc,
		"HeapIdle":      m.HeapIdle,
		"HeapInuse":     m.HeapInuse,
		"HeapObjects":   m.HeapObjects,
		"HeapReleased":  m.HeapReleased,
		"HeapSys":       m.HeapSys,
		"LastGC":        m.LastGC,
		"Lookups":       m.Lookups,
		"MCacheInuse":   m.MCacheInuse,
		"MCacheSys":     m.MCacheSys,
		"MSpanInuse":    m.MSpanInuse,
		"MSpanSys":      m.MSpanSys,
		"Mallocs":       m.Mallocs,
		"NextGC":        m.NextGC,
		"NumForcedGC":   m.NumForcedGC,
		"NumGC":         m.NumGC,
		"OtherSys":      m.OtherSys,
		"PauseTotalNs":  m.PauseTotalNs,
		"StackInuse":    m.StackInuse,
		"StackSys":      m.StackSys,
		"Sys":           m.Sys,
		"TotalAlloc":    m.TotalAlloc,
	}
}
