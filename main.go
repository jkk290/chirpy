package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/jkk290/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Printf("Error connecting to postgres database: %s", err)
	}
	defer db.Close()

	const filepathRoot = "."
	const port = "8080"

	apiCfg := &apiConfig{
		dbQueries: database.New(db),
	}

	mux := http.NewServeMux()
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("POST /api/validate_chirp", validateChirp)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerFileserverHitsCount)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerFileserverHitsReset)

	srv := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}
