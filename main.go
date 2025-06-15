package main

import "net/http"

func main() {
	mux := http.NewServeMux()

	// use .Handle to register the custom handler for '/'
	mux.Handle("/", http.FileServer(http.Dir(".")))

	// Struct that describes a server config
	server := &http.Server{
		Addr: ":8080", // listen to port 8080
		Handler: mux, // user custom ServeMux
	}
	// Start the server
	// The main function blocks until the server is shut down
	server.ListenAndServe()
}
