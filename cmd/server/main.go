package main

import (
	"log"
	"net/http"

	"github.com/Tamer-Li/devops-metrics/internal/handler"
	"github.com/Tamer-Li/devops-metrics/internal/storage"
)

func main() {

	memStorage := storage.NewMemStorage()

	apiHandler := handler.NewHandler(memStorage)

	mux := http.NewServeMux()
	mux.HandleFunc(`/metrics`, apiHandler.DataMetrics)
	mux.HandleFunc(`/update/`, apiHandler.MetricsHandler)

	log.Println("Starting server on :8080")
	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		log.Fatal(err)
	}
}
