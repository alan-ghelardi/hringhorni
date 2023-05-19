package server

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"sync"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	httpResponseTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_response_total",
		Help: "Total number of HTTP responses by status code",
	}, []string{"code"})

	statusCodes = []int{200, 202, 400, 422, 500, 503}
)

func StartInitialServer() {
	start(1, 2)
}

func StartUnstableServer() {
	start(2, len(statusCodes))
}

func StartStableServer() {
	start(3, 2)
}

func start(version, n int) {
	var waitGroup sync.WaitGroup
	waitGroup.Add(2)

	go func() {
		defer waitGroup.Done()

		http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			i, err := rand.Int(rand.Reader, big.NewInt(int64(n)))
			if err != nil {
				panic(err)
			}

			statusCode := statusCodes[i.Int64()]

			httpResponseTotal.With(prometheus.Labels{
				"code": fmt.Sprint(statusCode),
			}).Inc()

			w.WriteHeader(statusCode)
			if statusCode < 400 {
				fmt.Fprintf(w, "Hello from version %d!", version)
			} else {
				fmt.Fprintf(w, "Version %d: there's been a glitch!", version)
			}
		}))
		if err := http.ListenAndServe(":8080", nil); err != nil { //nolint:golint,gosec
			log.Fatal(err)
		}
	}()

	go func() {
		defer waitGroup.Done()
		http.Handle("/metrics", promhttp.Handler())
		if err := http.ListenAndServe(":2112", nil); err != nil { //nolint:golint,gosec
			log.Fatal(err)
		}
	}()

	waitGroup.Wait()
}
