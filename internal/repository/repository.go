package repository

type MetricsStorage interface {
	GaugeSet(name string, value float64)
	CounterSet(name string, value int64)
	Gauge() map[string]float64
	Counter() map[string]int64
}

type RepoMemStats interface {
	Update()
	GetAll() map[string]float64
}
