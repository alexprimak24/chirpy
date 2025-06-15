package main

import (
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
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("POST /api/validate_chirp", handlerChirpsValidate)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

	// Struct that describes a server config
	server := &http.Server{
		Addr: ":8080", // listen to port 8080
		Handler: mux, // user custom ServeMux
	}
	// Start the server
	// The main function blocks until the server is shut down
	server.ListenAndServe()
}