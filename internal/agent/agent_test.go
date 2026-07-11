package agent

import (
	"testing"

	"github.com/Tamer-Li/devops-metrics/internal/repository"
)

func TestAgent_pushMetrics(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		url            string
		pollInterval   int
		reportInterval int
		memData        repository.RepoMemStats
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := NewAgentMetric(tt.url, tt.pollInterval, tt.reportInterval, tt.memData)
			a.pushMetrics()
		})
	}
}

func TestAgent_postMetrics(t *testing.T) {
	tests := []struct {
		nameTest string // description of this test case
		// Named input parameters for receiver constructor.
		url            string
		pollInterval   int
		reportInterval int
		memData        repository.RepoMemStats
		// Named input parameters for target function.
		metric string
		name   string
		value  string
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := NewAgentMetric(tt.url, tt.pollInterval, tt.reportInterval, tt.memData)
			got := a.postMetrics(tt.metric, tt.name, tt.value)
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("postMetrics() = %v, want %v", got, tt.want)
			}
		})
	}
}
