package main

import (
	"database/sql"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/connor317/bootdev-chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries      *database.Queries
	platform       string
	jwtSecret      string
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")
	jwtSecret := os.Getenv("JWT_SECRET")
	db, _ := sql.Open("postgres", dbURL)
	dbQueries := database.New(db)
	mux := http.NewServeMux()
	handler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
	apiCfg := &apiConfig{
		dbQueries: dbQueries,
		platform:  platform,
		jwtSecret: jwtSecret,
	}

	mux.Handle("/app/", apiCfg.middlewareMetricsInc(handler))
	mux.Handle("/app/assets/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app/assets/", http.FileServer(http.Dir("assets")))))
	mux.HandleFunc("GET /api/healthz", handlerHealthz)

	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)

	mux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)
	mux.HandleFunc("POST /api/login", apiCfg.handlerLoginUser)

	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerCreateChirp)
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerGetAllChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerGetChirpByID)
	var srv http.Server
	srv.Addr = ":8080"
	srv.Handler = mux

	srv.ListenAndServe()
}
