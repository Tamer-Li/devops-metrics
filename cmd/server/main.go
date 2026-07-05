package main

import (
	"log"
	"net/http"

	"github.com/Tamer-Li/devops-metrics/internal/handler"
	"github.com/Tamer-Li/devops-metrics/internal/storage"
)

func main() {

	memStorage := storage.NewMemStorage()

	apiRouter := handler.NewHandler(memStorage).Router()

	log.Println("Starting server on :8080")
	err := http.ListenAndServe(`:8080`, apiRouter)
	if err != nil {
		log.Fatal(err)
	}
}
