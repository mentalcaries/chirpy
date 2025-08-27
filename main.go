package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/mentalcaries/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	DBQueries database.Queries
	platform string
	jwtSecret string
}

func main() {
	godotenv.Load()
	const filePathRoot = "."
	const port = "8080"
	mux := http.NewServeMux()
	
	dbURL := os.Getenv("DB_URL")
	env := os.Getenv("PLATFORM")
	jwtSecret := os.Getenv("JWT_SECRET")

	db, err := sql.Open("postgres", dbURL)
	if err != nil{
		log.Fatal("could not connect to database")
	}
	
	dbQueries := database.New(db)
	apiCfg := apiConfig{ DBQueries: *dbQueries, platform: env, jwtSecret: jwtSecret}
	
	mux.Handle("/app/", http.StripPrefix("/app", apiCfg.middlewareMetricsInc(http.FileServer(http.Dir(filePathRoot)))))
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerHitCounter)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	mux.HandleFunc("GET /api/healthz", handlerReadiness)

	
	mux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)
	mux.HandleFunc("POST /api/login", apiCfg.handlerLoginUser)
	mux.HandleFunc("PUT /api/users", apiCfg.handleUpdateUser)
	
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerGetAllChirps)
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerCreateChirp)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerGetChirp)

	mux.HandleFunc("POST /api/refresh", apiCfg.verifyRefreshToken)
	mux.HandleFunc("POST /api/revoke", apiCfg.revokeRefreshToken)

	server := http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}
	log.Printf("Serving files from %s on port: %s\n", filePathRoot, port)
	log.Fatal(server.ListenAndServe())
}
