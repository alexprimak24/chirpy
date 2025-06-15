package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"github.com/chirpy/internal/database"
	"github.com/joho/godotenv"

	// Underscore tells Go that you are importing it for its side effects not because you need to use it
	_ "github.com/lib/pq"
)

type apiConfig struct {
	// atomic.Int32 allows to safely increment and 
	// read an integer value across multiple goroutines
	fileserverHits atomic.Int32
	db *database.Queries
	platform string
}

func main() {
	godotenv.Load()
	const port = "8080"

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}

	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatal("PLATFORM must be set")
	}

	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening database: %s", err)
	}

	dbQueries := database.New(dbConn)

	// created an instance of apiConfig
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db: dbQueries,
		platform: platform,
	}
	
	mux := http.NewServeMux()


	// Serve static file users /app/ by stripping prefix
	// before handling it to the file server
	fileServer := http.FileServer(http.Dir("."))
	mux.Handle("/app/", http.StripPrefix("/app", apiCfg.middlewareMetricsInc(fileServer)))

	// Register a different handler for the root path
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("POST /api/validate_chirp", handlerChirpsValidate)
	mux.HandleFunc("POST /api/users", apiCfg.createUser)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

	// Struct that describes a server config
	server := &http.Server{
		Addr: ":" + port, // listen to port 8080
		Handler: mux, // user custom ServeMux
	}
	// Start the server
	// The main function blocks until the server is shut down
	log.Printf("Serving on port: %s\n", port)
	server.ListenAndServe()
}