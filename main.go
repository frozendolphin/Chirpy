package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/frozendolphin/Chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db *database.Queries
	platform string
	secret string
}

func main() {

	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	userPlatform := os.Getenv("PLATFORM")
	jwtSecret := os.Getenv("SECRET")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("couldn't open connection with database: %v", err)
	}

	dbQueries := database.New(db)

	mux := http.NewServeMux()

	apicfg := apiConfig {
		db: dbQueries,
		platform: userPlatform,
		secret: jwtSecret,
	}

	mux.Handle("/app/", apicfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	mux.HandleFunc("GET /api/healthz", healthz)
	mux.HandleFunc("GET /admin/metrics", apicfg.getHits)
	mux.HandleFunc("POST /admin/reset", apicfg.resetHits)
	mux.HandleFunc("POST /api/chirps", apicfg.createChirps)
	mux.HandleFunc("POST /api/users", apicfg.createUsers)
	mux.HandleFunc("GET /api/chirps", apicfg.getAllChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apicfg.getAChirp)
	mux.HandleFunc("POST /api/login", apicfg.loginUser)
	mux.HandleFunc("POST /api/refresh", apicfg.newRefresh)
	mux.HandleFunc("POST /api/revoke", apicfg.revokeRefresh)
	
	server_struct := http.Server {
		Handler: mux,
		Addr: ":8080",
	}

	err = server_struct.ListenAndServe()
	if err != nil {
		log.Fatalf("err occured: %v", err)
	}
}