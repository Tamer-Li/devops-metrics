package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func main() {

	memStorage := NewMemStorage()

	apiHandler := NewHandler(memStorage)

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

// -----------------------------------------------------------------
// Блок Handler

type MetricsStorage interface {
	GaugeSet(name string, value float64)
	CounterSet(name string, value int64)
	Gauge() map[string]float64
	Counter() map[string]int64
}

type Handler struct {
	MS MetricsStorage
}

func NewHandler(ms MetricsStorage) *Handler {
	return &Handler{
		MS: ms,
	}
}

func (h Handler) MetricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	if r.Method != http.MethodPost {
		log.Printf("Rejected: method %s not allowed", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	urlPath := strings.Split(r.URL.Path, "/")
	if len(urlPath) != 5 {
		log.Printf("Bad URL: %s (expected 5 parts)", r.URL.Path)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	metricType := urlPath[2]
	name := urlPath[3]
	value := urlPath[4]

	if name == "" {
		log.Printf("Empty metric name in URL: %s", r.URL.Path)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	switch metricType {
	case "gauge":
		val, err := strconv.ParseFloat(value, 64)
		if err != nil {
			log.Printf("Invalid gauge value %q for metric %q: %v", value, name, err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Printf("Gauge updated: %s = %f", name, val)
		h.MS.GaugeSet(name, val)

	case "counter":
		val, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			log.Printf("Invalid counter value %q for metric %q: %v", value, name, err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Printf("Counter updated: %s += %d", name, val)
		h.MS.CounterSet(name, val)

	default:
		log.Printf("Unknown metric type %q in request %s", metricType, r.URL.Path)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h Handler) DataMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	if r.Method != http.MethodGet {
		log.Printf("Rejected: method %s not allowed on /metrics", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	body := "====== Metrics ======\n"

	body += "\n### Gauge\n"

	gauges := h.MS.Gauge()

	if len(gauges) == 0 {
		body += "\tNot found\n"
	} else {
		for name, value := range gauges {
			body += fmt.Sprintf("\t%s = %f\n", name, value)
		}
	}

	body += "\n### Counter\n"

	counters := h.MS.Counter()

	if len(counters) == 0 {
		body += "\tNot found\n"
	} else {
		for name, value := range counters {
			body += fmt.Sprintf("\t%s = %d\n", name, value)
		}
	}

	log.Printf("Metrics:\nGauges = %v\nCounters = %v", gauges, counters)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(body))
}

// -----------------------------------------------------------------
