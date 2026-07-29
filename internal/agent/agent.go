package agent

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/Tamer-Li/devops-metrics/internal/arc"
	models "github.com/Tamer-Li/devops-metrics/internal/model"
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

func (a *Agent) postMetricJSON(metricType, name string, value float64, delta int64) bool {
	var metric models.Metrics
	metric.ID = name
	metric.MType = metricType

	switch metricType {
	case models.Gauge:
		metric.Value = &value
		metric.Delta = nil
	case models.Counter:
		metric.Delta = &delta
		metric.Value = nil
	default:
		log.Printf("Unsupported metric type: %s", metricType)
		return false
	}

	jsonData, err := json.Marshal(metric)
	if err != nil {
		log.Printf("Failed to marshal metric %s/%s: %v", metricType, name, err)
		return false
	}

	url := fmt.Sprintf("%s/update", a.url)

	resp, err := a.compress(url, metric, jsonData)
	if err != nil {
		log.Printf("Failed to send metric %s/%s: %v", metricType, name, err)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Non-OK status for %s/%s: %d", metricType, name, resp.StatusCode)
		return false
	}

	log.Printf("Successfully sent metric %s/%s", metricType, name)
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

	err = a.postMetricJSON("gauge", "RandomValue", a.randomValue, 0)
	if !err {
		log.Printf("Failed send JSON metrics gauge RandomValue = %s", randGauge)
	}

	err = a.postMetricJSON("counter", "PollCount", 0, a.pollCount)
	if !err {
		log.Printf("Failed send JSON metrics counter PollCount = %s", counter)
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

// func (a Agent) noCompress(url string, metric models.Metrics, data []byte) (*http.Response, error) {
// 	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
// 	if err != nil {
// 		return nil, fmt.Errorf("Failed to create request:%v\nError: %v", metric, err)
// 	}

// 	req.Header.Set("Content-Type", "application/json")

// 	return a.client.Do(req)
// }

func (a Agent) compress(url string, metric models.Metrics, data []byte) (*http.Response, error) {
	compressedData, err := arc.Compress(data)
	if err != nil {
		return nil, fmt.Errorf("failed to compress data: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, compressedData)
	if err != nil {
		return nil, fmt.Errorf("Failed to create request:%v\nError: %v", metric, err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")

	return a.client.Do(req)
}
