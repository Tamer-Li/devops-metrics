package handler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

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

			if tt.name == "counter increment" {
				ms.CounterSet("req", 5)
			}

			srv := NewHandler(ms)
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
			for k, v := range tt.gauges {
				ms.GaugeSet(k, v)
			}
			for k, v := range tt.counters {
				ms.CounterSet(k, v)
			}

			srv := NewHandler(ms)
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
			for k, v := range tt.gauges {
				ms.GaugeSet(k, v)
			}
			for k, v := range tt.counters {
				ms.CounterSet(k, v)
			}

			srv := NewHandler(ms)
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
