package repository

type MetricsStorage interface {
	GaugeSet(name string, value float64)
	CounterSet(name string, value int64)
	Gauges() map[string]float64
	Counters() map[string]int64
	Gauge(name string) (float64, bool)
	Counter(name string) (int64, bool)
}

type RepoMemStats interface {
	Update()
	GetAll() map[string]float64
}
