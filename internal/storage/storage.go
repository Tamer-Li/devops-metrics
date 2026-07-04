package storage

type MemStorage struct {
	gauge   map[string]float64
	counter map[string]int64
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		gauge:   make(map[string]float64),
		counter: make(map[string]int64),
	}
}

func (ms *MemStorage) Gauge() map[string]float64 {
	return ms.gauge
}

func (ms *MemStorage) Counter() map[string]int64 {
	return ms.counter
}

func (ms *MemStorage) GaugeSet(name string, value float64) {
	ms.gauge[name] = value
}

func (ms *MemStorage) CounterSet(name string, value int64) {
	ms.counter[name] += value
}
