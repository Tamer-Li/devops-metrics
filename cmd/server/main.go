package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type MemStorage struct {
	gauge   map[string]float64
	counter map[string][]int64
}

var memStorage = MemStorage{
	gauge:   make(map[string]float64),
	counter: make(map[string][]int64),
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc(`/`, mainPage)
	mux.HandleFunc(`/update/`, gaugeHandler)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}

func mainPage(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Main Page"))
	fmt.Println(memStorage)
}

func gaugeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		rData := r.URL.Path
		fmt.Println(rData)
		return
	}

	urlPath := strings.Split(r.URL.Path, "/")
	if len(urlPath) != 5 {
		fmt.Printf("Error more path: %s\n", r.URL.Path)
		return
	}

	fmt.Println(urlPath)

	if urlPath[2] == "gauge" {
		nameGauge := urlPath[3]
		countGauge := urlPath[4]

		value, err := strconv.ParseFloat(countGauge, 64)
		if err != nil {
			fmt.Printf("Error type: %s\n", countGauge)
			return
		}
		memStorage.gauge[nameGauge] = value
		return
	}

	if urlPath[2] == "counter" {
		nameCounter := urlPath[3]
		countCounter := urlPath[4]

		valueCounter, err := strconv.ParseInt(countCounter, 10, 64)
		if err != nil {
			fmt.Printf("Error type: %s\n", countCounter)
			return
		}
		if memStorage.counter[nameCounter] == nil {
			memStorage.counter[nameCounter] = make([]int64, 0, 8)
		}
		memStorage.counter[nameCounter] = append(memStorage.counter[nameCounter], valueCounter)
		return
	}
	fmt.Println("NULL")
}
