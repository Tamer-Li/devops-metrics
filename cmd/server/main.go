package main

import (
	"log"
	"net/http"

	"github.com/Tamer-Li/devops-metrics/internal/handler"
)

func main() {

	memStorage := NewMemStorage()

	apiHandler := handler.NewHandler(memStorage)

	mux := http.NewServeMux()
	mux.HandleFunc(`/metrics`, apiHandler.DataMetrics)
	mux.HandleFunc(`/update/`, apiHandler.MetricsHandler)

	log.Println("Starting server on :8080")
	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		log.Fatal(err)
	}
}

// -----------------------------------------------------------------
// Блок MemStorage

type MemStorage struct {
	gauge   map[string]float64
	counter map[string]int64
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		gauge:   make(map[string]float64),
		counter: make(map[string]int64),
	}
}

func (ms *MemStorage) Gauge() map[string]float64 {
	return ms.gauge
}

func (ms *MemStorage) Counter() map[string]int64 {
	return ms.counter
}

func (ms *MemStorage) GaugeSet(name string, value float64) {
	ms.gauge[name] = value
}

func (ms *MemStorage) CounterSet(name string, value int64) {
	ms.counter[name] += value
}

// -----------------------------------------------------------------
