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
	client         *http.Client
	pollInterval   int
	reportInterval int
	pollCount      int64
	randomValue    float64
	memData        repository.RepoMemStats
}

func NewAgentMetric(
	url string,
	pollInterval int,
	reportInterval int,
	memData repository.RepoMemStats,
) *Agent {
	return &Agent{
		url:            url,
		client:         &http.Client{},
		pollInterval:   pollInterval,
		reportInterval: reportInterval,
		pollCount:      0,
		randomValue:    0,
		memData:        memData,
	}
}

func (a *Agent) postMetrics(metric, name, value string) bool {
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

func (a *Agent) updateMetric() {
	a.memData.Update()
	a.randomValue = rand.Float64()
	a.pollCount++
}

func (a *Agent) pushMetrics() {
	for k, v := range a.memData.GetAll() {
		valStr := strconv.FormatFloat(v, 'f', -1, 64)
		err := a.postMetrics("gauge", k, valStr)
		if !err {
			log.Printf("Failed send metrics gauge %s = %s", k, valStr)
		}
	}
	randGauge := strconv.FormatFloat(a.randomValue, 'f', -1, 64)
	counter := strconv.FormatInt(a.pollCount, 10)
	err := a.postMetrics("gauge", "RandomValue", randGauge)
	if !err {
		log.Printf("Failed send metrics gauge RandomValue = %s", randGauge)
	}

	err = a.postMetrics("counter", "PollCount", counter)
	if !err {
		log.Printf("Failed send metrics counter PollCount = %s", counter)
	}
}

func (a *Agent) Run() {
	pollTicker := time.NewTicker(time.Duration(a.pollInterval) * time.Second)
	reportTicker := time.NewTicker(time.Duration(a.reportInterval) * time.Second)
	defer pollTicker.Stop()
	defer reportTicker.Stop()

	log.Println("Agent run")

	for {
		select {
		case <-pollTicker.C:
			log.Println("Agent update")
			a.updateMetric()
		case <-reportTicker.C:
			log.Println("Agent push metrics")
			a.pushMetrics()
		}
	}
}
