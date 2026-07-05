package handler

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Tamer-Li/devops-metrics/internal/repository"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	MS repository.MetricsStorage
}

func NewHandler(ms repository.MetricsStorage) *Handler {
	return &Handler{
		MS: ms,
	}
}

func (h *Handler) Router() chi.Router {
	router := chi.NewRouter()

	router.Get("/", h.metricsHandler)
	router.Get("/value/{typeMetric}/{nameMetric}", h.metricValue)
	router.Post("/update/{typeMetric}/{nameMetric}/{valueMetric}", h.updateMetric)

	return router
}

func (h *Handler) metricsHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "text/html")
	rw.WriteHeader(http.StatusOK)

	body := `
<html>
<head><title>Metrics</title></head>
<body>
<h1>Metrics</h1>
<table border="1">
<tr><th>Type</th><th>Name</th><th>Value</th></tr>	
	`
	for name, value := range h.MS.Gauges() {
		body += fmt.Sprintf("<tr><td>gauge</td><td>%s</td><td>%f</td></tr>", name, value)
	}

	for name, value := range h.MS.Counters() {
		body += fmt.Sprintf("<tr><td>counter</td><td>%s</td><td>%d</td></tr>", name, value)
	}

	body += "</table></body></html>"

	rw.Write([]byte(body))
}

func (h *Handler) metricValue(rw http.ResponseWriter, r *http.Request) {
	typeMetric := chi.URLParam(r, "typeMetric")
	nameMetric := chi.URLParam(r, "nameMetric")

	switch typeMetric {
	case "gauge":
		value, ok := h.MS.Gauge(nameMetric)
		if !ok {
			rw.WriteHeader(http.StatusNotFound)
			return
		}
		rw.Header().Set("Content-Type", "text/plain")
		rw.WriteHeader(http.StatusOK)

		body := strconv.FormatFloat(value, 'f', 0, 64)

		rw.Write([]byte(body))
	case "counter":
		value, ok := h.MS.Counter(nameMetric)
		if !ok {
			rw.WriteHeader(http.StatusNotFound)
			return
		}
		rw.Header().Set("Content-Type", "text/plain")
		rw.WriteHeader(http.StatusOK)

		body := strconv.FormatInt(value, 64)

		rw.Write([]byte(body))
	default:
		rw.WriteHeader(http.StatusNotFound)
	}
}

func (h *Handler) updateMetric(rw http.ResponseWriter, r *http.Request) {
	typeMetric := chi.URLParam(r, "typeMetric")
	name := chi.URLParam(r, "nameMetric")
	value := chi.URLParam(r, "valueMetric")

	if name == "" {
		log.Println("Empty metric name")
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	switch typeMetric {
	case "gauge":
		val, err := strconv.ParseFloat(value, 64)
		if err != nil {
			log.Printf("Invalid gauge value %q for metric %q: %v", value, name, err)
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Printf("Gauge updated: %s = %f", name, val)
		h.MS.GaugeSet(name, val)

	case "counter":
		val, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			log.Printf("Invalid counter value %q for metric %q: %v", value, name, err)
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Printf("Counter updated: %s += %d", name, val)
		h.MS.CounterSet(name, val)

	default:
		log.Printf("Unknown metric type %q in request %s", typeMetric, r.URL.Path)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	rw.WriteHeader(http.StatusOK)
}
