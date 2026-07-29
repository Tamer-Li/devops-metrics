package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/Tamer-Li/devops-metrics/internal/arc"
	"github.com/Tamer-Li/devops-metrics/internal/config"
	"github.com/Tamer-Li/devops-metrics/internal/handler"
	"github.com/Tamer-Li/devops-metrics/internal/logger"
	"github.com/Tamer-Li/devops-metrics/internal/storage"
)

var settings struct {
	address string
}

func init() {
	flag.StringVar(&settings.address, "a", "localhost:8080", "address endpoint http-server")
}

func main() {
	flag.Parse()

	cfg := config.NewConfig()
	if cfg != nil {
		settings.address = cfg.Address
	}

	logHome := logger.NewHomeLogger()

	memStorage := storage.NewMemStorage()

	apiRouter := handler.NewHandler(memStorage).Router()

	log.Printf("Starting server on %s", settings.address)
	err := http.ListenAndServe(settings.address, logHome.WithLogging(arc.GZIPHandle(apiRouter)))
	if err != nil {
		panic(err)
	}
}
