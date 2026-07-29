package files

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	models "github.com/Tamer-Li/devops-metrics/internal/model"
	"github.com/Tamer-Li/devops-metrics/internal/repository"
)

type FilesStorage struct {
	memStorage    repository.MetricsStorage
	filePath      string
	storeInterval int
	stopChan      chan struct{}
}

func NewFilesStorage(memStorage repository.MetricsStorage, filePath string, storeInterval int) *FilesStorage {
	return &FilesStorage{
		memStorage:    memStorage,
		filePath:      filePath,
		storeInterval: storeInterval,
		stopChan:      make(chan struct{}),
	}
}

func (fs *FilesStorage) Save() error {
	gauges := fs.memStorage.Gauges()
	counters := fs.memStorage.Counters()

	var metrics []models.Metrics

	for name, value := range gauges {
		val := value
		metrics = append(metrics, models.Metrics{
			ID:    name,
			MType: models.Gauge,
			Value: &val,
			Delta: nil,
		})
	}

	for name, value := range counters {
		delta := value
		metrics = append(metrics, models.Metrics{
			ID:    name,
			MType: models.Counter,
			Value: nil,
			Delta: &delta,
		})
	}

	file, err := os.Create(fs.filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(metrics); err != nil {
		return fmt.Errorf("failed to encode metrics to JSON: %w", err)
	}

	return nil
}

func (fs *FilesStorage) Load() error {

	if _, err := os.Stat(fs.filePath); os.IsNotExist(err) {
		return nil
	}

	file, err := os.Open(fs.filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var metrics []models.Metrics
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&metrics); err != nil {
		return fmt.Errorf("failed to decode JSON: %w", err)
	}

	for _, metric := range metrics {
		switch metric.MType {
		case models.Gauge:
			if metric.Value != nil {
				fs.memStorage.GaugeSet(metric.ID, *metric.Value)
			}
		case models.Counter:
			if metric.Delta != nil {
				fs.memStorage.CounterSet(metric.ID, *metric.Delta)
			}
		}
	}

	return nil
}

func (fs *FilesStorage) StartIntervalSave() {
	ticker := time.NewTicker(time.Duration(fs.storeInterval) * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				if err := fs.Save(); err != nil {
					fmt.Printf("Failed to save storage: %v\n", err)
				}
			case <-fs.stopChan:
				ticker.Stop()
				return
			}
		}
	}()
}
