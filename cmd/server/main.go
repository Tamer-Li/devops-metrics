package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/Tamer-Li/devops-metrics/internal/arc"
	"github.com/Tamer-Li/devops-metrics/internal/config"
	"github.com/Tamer-Li/devops-metrics/internal/files"
	"github.com/Tamer-Li/devops-metrics/internal/handler"
	"github.com/Tamer-Li/devops-metrics/internal/logger"
	"github.com/Tamer-Li/devops-metrics/internal/storage"
)

var settings struct {
	address         string
	storeInterval   int
	fileStoragePath string
	restore         bool
}

func init() {
	flag.StringVar(&settings.address, "a", "localhost:8080", "address endpoint http-server")
	flag.IntVar(&settings.storeInterval, "i", 300, "Interval save storage")
	flag.StringVar(&settings.fileStoragePath, "f", "./storage.json", "path to file storage")
	flag.BoolVar(&settings.restore, "r", false, "download save file")
}

func main() {
	flag.Parse()

	cfg := config.NewConfig()
	if cfg != nil {
		settings.address = cfg.Address
		settings.storeInterval = cfg.StoreInterval
		settings.fileStoragePath = cfg.FileStoragePath
		settings.restore = cfg.Restore
	}

	logHome := logger.NewHomeLogger()

	memStorage := storage.NewMemStorage()

	fileStorage := files.NewFilesStorage(
		memStorage,
		settings.fileStoragePath,
		settings.storeInterval,
	)

	if settings.restore {
		fileStorage.Load()
	}

	if settings.storeInterval > 0 {
		fileStorage.StartIntervalSave()
	}

	apiRouter := handler.NewHandler(memStorage, fileStorage).Router()

	log.Printf("Starting server on %s", settings.address)
	err := http.ListenAndServe(settings.address, logHome.WithLogging(arc.GZIPHandle(apiRouter)))
	if err != nil {
		panic(err)
	}
}
