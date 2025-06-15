package main

import "net/http"

func main() {
	mux := http.NewServeMux()

	// Serve static file users /app/ by stripping prefix
	// before handling it to the file server
	fileServer := http.FileServer(http.Dir("."))
	mux.Handle("/app/", http.StripPrefix("/app", fileServer))

	// Register a different handler for the root path
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	// Struct that describes a server config
	server := &http.Server{
		Addr: ":8080", // listen to port 8080
		Handler: mux, // user custom ServeMux
	}
	// Start the server
	// The main function blocks until the server is shut down
	server.ListenAndServe()
}
