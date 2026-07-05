package main

import (
	"flag"
	"fmt"

	"github.com/Tamer-Li/devops-metrics/internal/agent"
	"github.com/Tamer-Li/devops-metrics/internal/service"
)

var settings struct {
	address        string
	reportInterval int
	pollInterval   int
}

func init() {
	flag.StringVar(&settings.address, "a", "localhost:8080", "address endpoint http-server")
	flag.IntVar(&settings.reportInterval, "r", 10, "report interval in seconds")
	flag.IntVar(&settings.pollInterval, "p", 2, "poll interval in seconds")
}

func main() {
	flag.Parse()
	url := fmt.Sprintf("http://%s", settings.address)
	memStats := service.NewMemStatsRuntime()

	agentMy := agent.NewAgentMetric(
		url,
		settings.pollInterval,
		settings.reportInterval,
		memStats,
	)

	agentMy.Run()
}
