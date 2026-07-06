package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Tamer-Li/devops-metrics/internal/repository"
	"github.com/go-chi/chi/v5"
)

func TestHandler_updateMetric(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}

	tests := []struct {
		name   string
		method string
		url    string
		want   want
	}{
		{
			name:   "Test 1",
			method: http.MethodGet,
			url:    "/",
			want: want{
				code:        200,
				response:    `{"status":"ok"}`,
				contentType: "text/html",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.url, nil)
			w := httptest.NewRecorder()
			r := NewHandler()
		})
	}
}

func TestHandler_metricValue(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		ms repository.MetricsStorage
		// Named input parameters for target function.
		rw http.ResponseWriter
		r  *http.Request
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHandler(tt.ms)
			h.metricValue(tt.rw, tt.r)
		})
	}
}

func TestHandler_metricsHandler(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		ms repository.MetricsStorage
		// Named input parameters for target function.
		rw http.ResponseWriter
		r  *http.Request
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHandler(tt.ms)
			h.metricsHandler(tt.rw, tt.r)
		})
	}
}

func TestHandler_Router(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}

	tests := []struct {
		name string
		url  string
		ms   repository.MetricsStorage
		want chi.Router
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHandler(tt.ms)
			got := h.Router()
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("Router() = %v, want %v", got, tt.want)
			}
		})
	}
}
