package agent

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/Tamer-Li/devops-metrics/internal/repository"
)

type Agent struct {
	url            string
	pollInterval   int64
	reportInterval int64
	pollCount      int64
	randomValue    float64
	memData        repository.RepoMemStats
	data           map[string]float64
}

func NewAgentMetric(
	url string,
	pollInterval int64,
	reportInterval int64,
	memData repository.RepoMemStats,
) *Agent {
	return &Agent{
		url:            url,
		pollInterval:   pollInterval,
		reportInterval: reportInterval,
		pollCount:      0,
		randomValue:    0,
		memData:        memData,
	}
}

func (a Agent) postMetrics(metric, name, value string) bool {
	client := &http.Client{}

	url := fmt.Sprintf("%s/update/%s/%s/%s", a.url, metric, name, value)

	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		log.Printf("Failed to create request: %v", err)
		return false
	}

	req.Header.Set("Content-Type", "text/plain")

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to send metric %s/%s: %v", metric, name, err)
		return false
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Non-OK status for %s/%s: %d", metric, name, resp.StatusCode)
		return false
	}

	return true
}

func (a Agent) updateMetric() {
	for {
		time.Sleep(time.Duration(a.pollInterval) * time.Second)
		a.memData.Update()
		for k, v := range a.memData.GetAll() {
			a.data[k] = v
		}

		a.randomValue = rand.Float64()
		a.pollCount++
	}
}

func (a Agent) PushMetrics() {
	for {
		time.Sleep(time.Duration(a.pollCount) * time.Second)
		for k, v := range a.data {
			valStr := strconv.FormatFloat(v, 'f', -1, 64)
			a.postMetrics("gauge", k, valStr)
		}
		randGauge := strconv.FormatFloat(a.randomValue, 'f', -1, 64)
		counter := strconv.FormatInt(a.pollCount, 10)
		a.postMetrics("gauge", "RandomValue", randGauge)
		a.postMetrics("counter", "PollCount", counter)
	}
}

func (a Agent) Run() {
	go a.PushMetrics()
	go a.updateMetric()
}
