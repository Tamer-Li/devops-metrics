package main

import (
	"github.com/Tamer-Li/devops-metrics/internal/agent"
	"github.com/Tamer-Li/devops-metrics/internal/service"
)

func main() {
	pollInterval := 2
	reportInterval := 10
	url := "http://localhost:8080"
	memStats := service.NewMemStatsRuntime()

	agentMy := agent.NewAgentMetric(
		url,
		pollInterval,
		reportInterval,
		memStats,
	)

	agentMy.Run()
}
