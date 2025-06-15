package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	// atomic.Int32 allows to safely increment and 
	// read an integer value across multiple goroutines
	fileserverHits atomic.Int32
}

func main() {
	mux := http.NewServeMux()

	// created an instance of apiConfig
	apiCfg := &apiConfig{}

	// Serve static file users /app/ by stripping prefix
	// before handling it to the file server
	fileServer := http.FileServer(http.Dir("."))
	mux.Handle("/app/", http.StripPrefix("/app", apiCfg.middlewareMetricsInc(fileServer)))

	// Register a different handler for the root path
	mux.HandleFunc("/healthz", handlerReadiness)
	mux.HandleFunc("/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("/reset", apiCfg.resetMetricsCount)

	// Struct that describes a server config
	server := &http.Server{
		Addr: ":8080", // listen to port 8080
		Handler: mux, // user custom ServeMux
	}
	// Start the server
	// The main function blocks until the server is shut down
	server.ListenAndServe()
}


func handlerReadiness(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(http.StatusText(http.StatusOK)))
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		// safely increment the counter
		cfg.fileserverHits.Add(1)
		// call the next handler
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	count := cfg.fileserverHits.Load()
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hits: %d", count)
}

func (cfg *apiConfig) resetMetricsCount(w http.ResponseWriter, r *http.Request){
	cfg.fileserverHits.Store(0)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}