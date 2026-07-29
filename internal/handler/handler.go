package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	models "github.com/Tamer-Li/devops-metrics/internal/model"
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

	router.Post("/value", h.metricValueJSON)
	router.Post("/update", h.updateMetricJSON)

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
		body += fmt.Sprintf("<tr><td>gauge</td><td>%s</td><td>%s</td></tr>", name, strconv.FormatFloat(value, 'f', -1, 64))
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
	case models.Gauge:
		value, ok := h.MS.Gauge(nameMetric)
		if !ok {
			rw.WriteHeader(http.StatusNotFound)
			return
		}
		rw.Header().Set("Content-Type", "text/plain")
		rw.WriteHeader(http.StatusOK)

		body := strconv.FormatFloat(value, 'f', -1, 64)

		rw.Write([]byte(body))
	case models.Counter:
		value, ok := h.MS.Counter(nameMetric)
		if !ok {
			rw.WriteHeader(http.StatusNotFound)
			return
		}
		rw.Header().Set("Content-Type", "text/plain")
		rw.WriteHeader(http.StatusOK)

		body := strconv.FormatInt(value, 10)

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
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	switch typeMetric {
	case models.Gauge:
		val, err := strconv.ParseFloat(value, 64)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		h.MS.GaugeSet(name, val)

	case models.Counter:
		val, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		h.MS.CounterSet(name, val)

	default:
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	rw.WriteHeader(http.StatusOK)
}

func (h *Handler) metricValueJSON(rw http.ResponseWriter, r *http.Request) {

	if r.Header.Get("Content-Type") != "application/json" {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Content-Type is Failed"))
		return
	}
	var metric models.Metrics
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Request body is Failed"))
		return
	}

	err = json.Unmarshal(buf.Bytes(), &metric)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Request body model is Failed"))
		return
	}

	switch metric.MType {
	case models.Gauge:
		value, ok := h.MS.Gauge(metric.ID)
		if !ok {
			rw.WriteHeader(http.StatusNotFound)
			rw.Write([]byte("No found"))
			return
		}
		metric.Value = &value
		resp, err := json.Marshal(metric)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("Server err json serialized"))
			return
		}

		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
		rw.Write(resp)

	case models.Counter:
		value, ok := h.MS.Counter(metric.ID)
		if !ok {
			rw.WriteHeader(http.StatusNotFound)
			rw.Write([]byte("No found"))
			return
		}
		metric.Delta = &value
		resp, err := json.Marshal(metric)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("Server err json serialized"))
			return
		}

		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
		rw.Write(resp)
	default:
		rw.WriteHeader(http.StatusNotFound)
	}
}

func (h *Handler) updateMetricJSON(rw http.ResponseWriter, r *http.Request) {

	if r.Header.Get("Content-Type") != "application/json" {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Content-Type is Failed"))
		return
	}

	var metric models.Metrics
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Request body is Failed"))
		return
	}

	err = json.Unmarshal(buf.Bytes(), &metric)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Request body model is Failed"))
		return
	}

	if metric.ID == "" {
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	switch metric.MType {
	case models.Gauge:
		h.MS.GaugeSet(metric.ID, *metric.Value)

	case models.Counter:
		h.MS.CounterSet(metric.ID, *metric.Delta)

	default:
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	rw.WriteHeader(http.StatusOK)
}
