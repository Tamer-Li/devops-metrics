package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Tamer-Li/devops-metrics/internal/files"
	models "github.com/Tamer-Li/devops-metrics/internal/model"
	"github.com/Tamer-Li/devops-metrics/internal/storage"
	"github.com/stretchr/testify/assert"
)

func TestHandler_updateMetric(t *testing.T) {
	tests := []struct {
		name        string
		url         string
		wantStatus  int
		wantGauge   map[string]float64
		wantCounter map[string]int64
	}{
		{
			name:       "valid gauge",
			url:        "/update/gauge/cpu/0.85",
			wantStatus: 200,
			wantGauge:  map[string]float64{"cpu": 0.85},
		},
		{
			name:        "valid counter",
			url:         "/update/counter/requests/10",
			wantStatus:  200,
			wantCounter: map[string]int64{"requests": 10},
		},
		{
			name:       "invalid gauge value (letters)",
			url:        "/update/gauge/bad/abc",
			wantStatus: 400,
		},
		{
			name:       "invalid counter value (float)",
			url:        "/update/counter/bad/3.14",
			wantStatus: 400,
		},
		{
			name:       "empty metric name",
			url:        "/update/gauge//5",
			wantStatus: 404,
		},
		{
			name:       "unknown metric type",
			url:        "/update/histogram/x/1",
			wantStatus: 400,
		},
		{
			name:        "counter increment",
			url:         "/update/counter/req/3",
			wantStatus:  200,
			wantCounter: map[string]int64{"req": 8},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := storage.NewMemStorage()
			repoFS := files.NewFilesStorage(
				ms,
				"./config.json",
				300,
			)

			if tt.name == "counter increment" {
				ms.CounterSet("req", 5)
			}

			srv := NewHandler(ms, repoFS)
			req := httptest.NewRequest(http.MethodPost, tt.url, nil)
			w := httptest.NewRecorder()
			srv.Router().ServeHTTP(w, req)

			res := w.Result()
			defer res.Body.Close()
			assert.Equal(t, tt.wantStatus, res.StatusCode)

			for k, v := range tt.wantGauge {
				val, ok := ms.Gauge(k)
				assert.True(t, ok, "gauge %s should exist", k)
				assert.Equal(t, v, val)
			}
			for k, v := range tt.wantCounter {
				val, ok := ms.Counter(k)
				assert.True(t, ok, "counter %s should exist", k)
				assert.Equal(t, v, val)
			}
		})
	}
}

func TestHandler_metricValue(t *testing.T) {
	tests := []struct {
		name       string
		gauges     map[string]float64
		counters   map[string]int64
		url        string
		wantStatus int
		wantBody   string
		wantCT     string
	}{
		{
			name:       "gauge exists",
			gauges:     map[string]float64{"temp": 36.6},
			counters:   map[string]int64{},
			url:        "/value/gauge/temp",
			wantStatus: 200,
			wantBody:   "36.6",
			wantCT:     "text/plain",
		},
		{
			name:       "gauge not found",
			gauges:     map[string]float64{},
			url:        "/value/gauge/unknown",
			wantStatus: 404,
		},
		{
			name:       "counter exists",
			counters:   map[string]int64{"hits": 42},
			url:        "/value/counter/hits",
			wantStatus: 200,
			wantBody:   "42",
			wantCT:     "text/plain",
		},
		{
			name:       "counter not found",
			counters:   map[string]int64{},
			url:        "/value/counter/missing",
			wantStatus: 404,
		},
		{
			name:       "invalid type",
			url:        "/value/histogram/something",
			wantStatus: 404,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := storage.NewMemStorage()
			repoFS := files.NewFilesStorage(
				ms,
				"./config.json",
				300,
			)
			for k, v := range tt.gauges {
				ms.GaugeSet(k, v)
			}
			for k, v := range tt.counters {
				ms.CounterSet(k, v)
			}

			srv := NewHandler(ms, repoFS)
			req := httptest.NewRequest(http.MethodGet, tt.url, nil)
			w := httptest.NewRecorder()
			srv.Router().ServeHTTP(w, req)

			res := w.Result()
			defer res.Body.Close()
			assert.Equal(t, tt.wantStatus, res.StatusCode)

			if tt.wantStatus == 200 {
				assert.Equal(t, tt.wantCT, res.Header.Get("Content-Type"))
				body, _ := io.ReadAll(res.Body)
				assert.Equal(t, tt.wantBody, string(body))
			}
		})
	}
}

func TestHandler_metricsHandler(t *testing.T) {
	tests := []struct {
		name        string
		gauges      map[string]float64
		counters    map[string]int64
		wantSubstrs []string
	}{
		{
			name:     "empty storage",
			gauges:   map[string]float64{},
			counters: map[string]int64{},
			wantSubstrs: []string{
				"<table border=\"1\">",
				"<th>Type</th>",
				"<th>Name</th>",
				"<th>Value</th>",
				"</table>",
			},
		},
		{
			name:     "with gauge and counter",
			gauges:   map[string]float64{"temp": 36.6},
			counters: map[string]int64{"hits": 42},
			wantSubstrs: []string{
				"<td>gauge</td>",
				"<td>temp</td>",
				"<td>36.6</td>",
				"<td>counter</td>",
				"<td>hits</td>",
				"<td>42</td>",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := storage.NewMemStorage()
			repoFS := files.NewFilesStorage(
				ms,
				"./config.json",
				300,
			)
			for k, v := range tt.gauges {
				ms.GaugeSet(k, v)
			}
			for k, v := range tt.counters {
				ms.CounterSet(k, v)
			}

			srv := NewHandler(ms, repoFS)
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			w := httptest.NewRecorder()
			srv.Router().ServeHTTP(w, req)

			res := w.Result()
			defer res.Body.Close()
			assert.Equal(t, 200, res.StatusCode)
			assert.Equal(t, "text/html", res.Header.Get("Content-Type"))

			body, _ := io.ReadAll(res.Body)
			html := string(body)
			for _, sub := range tt.wantSubstrs {
				assert.Contains(t, html, sub)
			}
		})
	}
}

func TestHandler_metricValueJSON(t *testing.T) {
	tests := []struct {
		name        string
		gauges      map[string]float64
		counters    map[string]int64
		requestBody interface{}
		wantStatus  int
		wantCT      string
		wantBody    interface{}
		wantSubstr  string
	}{
		{
			name:   "gauge exists",
			gauges: map[string]float64{"temp": 36.6},
			requestBody: models.Metrics{
				ID:    "temp",
				MType: models.Gauge,
			},
			wantStatus: 200,
			wantCT:     "application/json",
			wantBody: models.Metrics{
				ID:    "temp",
				MType: models.Gauge,
				Value: floatPtr(36.6),
			},
		},
		{
			name:     "counter exists",
			counters: map[string]int64{"hits": 42},
			requestBody: models.Metrics{
				ID:    "hits",
				MType: models.Counter,
			},
			wantStatus: 200,
			wantCT:     "application/json",
			wantBody: models.Metrics{
				ID:    "hits",
				MType: models.Counter,
				Delta: int64Ptr(42),
			},
		},
		{
			name: "gauge not found",
			requestBody: models.Metrics{
				ID:    "unknown",
				MType: models.Gauge,
			},
			wantStatus: 404,
			wantSubstr: "No found",
		},
		{
			name: "counter not found",
			requestBody: models.Metrics{
				ID:    "missing",
				MType: models.Counter,
			},
			wantStatus: 404,
			wantSubstr: "No found",
		},
		{
			name:        "invalid json",
			requestBody: []byte(`{"id": "test", "type": "gauge"`),
			wantStatus:  400,
			wantSubstr:  "Request body model is Failed",
		},
		{
			name:        "empty body",
			requestBody: []byte{},
			wantStatus:  400,
			wantSubstr:  "Request body model is Failed",
		},
		{
			name: "invalid metric type",
			requestBody: models.Metrics{
				ID:    "test",
				MType: "histogram",
			},
			wantStatus: 404,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := storage.NewMemStorage()
			repoFS := files.NewFilesStorage(
				ms,
				"./config.json",
				300,
			)

			for k, v := range tt.gauges {
				ms.GaugeSet(k, v)
			}
			for k, v := range tt.counters {
				ms.CounterSet(k, v)
			}

			var body io.Reader
			if tt.requestBody != nil {
				if _, ok := tt.requestBody.([]byte); ok {
					body = bytes.NewReader(tt.requestBody.([]byte))
				} else {
					jsonBody, _ := json.Marshal(tt.requestBody)
					body = bytes.NewReader(jsonBody)
				}
			} else {
				body = bytes.NewReader([]byte{})
			}

			srv := NewHandler(ms, repoFS)
			req := httptest.NewRequest(http.MethodPost, "/value", body)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			srv.Router().ServeHTTP(w, req)

			res := w.Result()
			defer res.Body.Close()
			assert.Equal(t, tt.wantStatus, res.StatusCode)

			if tt.wantStatus == 200 {
				assert.Equal(t, tt.wantCT, res.Header.Get("Content-Type"))
				respBody, _ := io.ReadAll(res.Body)
				var resp models.Metrics
				err := json.Unmarshal(respBody, &resp)
				assert.NoError(t, err)

				expected, _ := tt.wantBody.(models.Metrics)
				assert.Equal(t, expected.ID, resp.ID)
				assert.Equal(t, expected.MType, resp.MType)
				if expected.Value != nil {
					assert.Equal(t, *expected.Value, *resp.Value)
				}
				if expected.Delta != nil {
					assert.Equal(t, *expected.Delta, *resp.Delta)
				}
			}

			if tt.wantSubstr != "" {
				respBody, _ := io.ReadAll(res.Body)
				assert.Contains(t, string(respBody), tt.wantSubstr)
			}
		})
	}
}

func floatPtr(f float64) *float64 {
	return &f
}

func int64Ptr(i int64) *int64 {
	return &i
}

func TestHandler_updateMetricJSON(t *testing.T) {
	tests := []struct {
		name        string
		initialData map[string]interface{}
		requestBody interface{}
		wantStatus  int
		wantGauge   map[string]float64
		wantCounter map[string]int64
		wantSubstr  string
	}{
		{
			name: "update gauge",
			requestBody: models.Metrics{
				ID:    "cpu",
				MType: models.Gauge,
				Value: floatPtr(0.85),
			},
			wantStatus: 200,
			wantGauge:  map[string]float64{"cpu": 0.85},
		},
		{
			name: "update counter",
			requestBody: models.Metrics{
				ID:    "requests",
				MType: models.Counter,
				Delta: int64Ptr(10),
			},
			wantStatus:  200,
			wantCounter: map[string]int64{"requests": 10},
		},
		{
			name: "update existing gauge",
			initialData: map[string]interface{}{
				"gauge": map[string]float64{"cpu": 0.5},
			},
			requestBody: models.Metrics{
				ID:    "cpu",
				MType: models.Gauge,
				Value: floatPtr(0.85),
			},
			wantStatus: 200,
			wantGauge:  map[string]float64{"cpu": 0.85},
		},
		{
			name: "update existing counter",
			initialData: map[string]interface{}{
				"counter": map[string]int64{"requests": 5},
			},
			requestBody: models.Metrics{
				ID:    "requests",
				MType: models.Counter,
				Delta: int64Ptr(10),
			},
			wantStatus:  200,
			wantCounter: map[string]int64{"requests": 15},
		},
		{
			name:        "invalid json",
			requestBody: []byte(`{"id": "cpu", "type": "gauge", "value": "invalid"`),
			wantStatus:  400,
			wantSubstr:  "Request body model is Failed",
		},
		{
			name:        "empty body",
			requestBody: []byte{},
			wantStatus:  400,
			wantSubstr:  "Request body model is Failed",
		},
		{
			name: "empty metric id",
			requestBody: models.Metrics{
				ID:    "",
				MType: models.Gauge,
				Value: floatPtr(0.85),
			},
			wantStatus: 404,
		},
		{
			name: "invalid metric type",
			requestBody: models.Metrics{
				ID:    "test",
				MType: "histogram",
				Value: floatPtr(0.85),
			},
			wantStatus: 400,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := storage.NewMemStorage()
			repoFS := files.NewFilesStorage(
				ms,
				"./config.json",
				300,
			)

			if tt.initialData != nil {
				if gauges, ok := tt.initialData["gauge"].(map[string]float64); ok {
					for k, v := range gauges {
						ms.GaugeSet(k, v)
					}
				}
				if counters, ok := tt.initialData["counter"].(map[string]int64); ok {
					for k, v := range counters {
						ms.CounterSet(k, v)
					}
				}
			}

			var body io.Reader
			if tt.requestBody != nil {
				if _, ok := tt.requestBody.([]byte); ok {
					body = bytes.NewReader(tt.requestBody.([]byte))
				} else {
					jsonBody, _ := json.Marshal(tt.requestBody)
					body = bytes.NewReader(jsonBody)
				}
			} else {
				body = bytes.NewReader([]byte{})
			}

			srv := NewHandler(ms, repoFS)
			req := httptest.NewRequest(http.MethodPost, "/update", body)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			srv.Router().ServeHTTP(w, req)

			res := w.Result()
			defer res.Body.Close()
			assert.Equal(t, tt.wantStatus, res.StatusCode)

			if tt.wantSubstr != "" {
				respBody, _ := io.ReadAll(res.Body)
				assert.Contains(t, string(respBody), tt.wantSubstr)
			}

			for k, v := range tt.wantGauge {
				val, ok := ms.Gauge(k)
				assert.True(t, ok, "gauge %s should exist", k)
				assert.Equal(t, v, val)
			}
			for k, v := range tt.wantCounter {
				val, ok := ms.Counter(k)
				assert.True(t, ok, "counter %s should exist", k)
				assert.Equal(t, v, val)
			}
		})
	}
}
