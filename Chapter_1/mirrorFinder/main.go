package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/0oPhantom0o/Rest_Api_Book/Chapter_1/mirrors"
)

type response struct {
	FastestURL string        `json:"fastest_url"`
	Latency    time.Duration `json:"latency"`
}

func findFastest(urls []string) response {
	urlChan := make(chan string)
	latencyChan := make(chan time.Duration)
	for _, url := range urls {
		mirrorURL := url
		go func() {
			start := time.Now()
			_, err := http.Get(mirrorURL + "/README")
			latency := time.Now().Sub(start) / time.Millisecond
			if err == nil {
				urlChan <- mirrorURL
				latencyChan <- latency
			}
		}()

	}
	return response{<-urlChan, <-latencyChan}

}
func main() {
	fmt.Println("hi")
	http.HandleFunc("/fastest-mirror",
		func(w http.ResponseWriter, r *http.Request) {
			response := findFastest(mirrors.MirrorList)
			respJSON, _ := json.Marshal(response)
			w.Header().Set("Content-Type", "application/json")
			w.Write(respJSON)

		})
	port := "8000"
	server := &http.Server{
		Addr:           port,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	fmt.Println("Startubg server ib oirt %sn", port)
	log.Fatal(server.ListenAndServe())
}
