package main

import (
	"fmt"
	"net/http"
)

var gauge float64

type MemStorage struct {
	counters []int64
}

func mainPage(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Main Page"))
}

func gaugeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		rData := r.URL.Path
		fmt.Println(rData)
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc(`/`, mainPage)
	mux.HandleFunc(`/update/gauge`, gaugeHandler)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
