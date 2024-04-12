package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ammon134/chirpy/internal/database"
)

const (
	filePathRoot = "."
	port         = "8080"
	dbPath       = "database.json"
)

func main() {
	mux := http.NewServeMux()
	corsMux := middlewareCors(mux)
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      corsMux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	db, err := database.NewDB(dbPath)
	if err != nil {
		log.Fatal(err)
	}
	apiConfig := &apiConfig{
		serverHits: 0,
		db:         db,
	}

	mux.Handle("/app/*", apiConfig.middlewareHitInc(http.StripPrefix("/app/", http.FileServer(http.Dir(filePathRoot)))))

	mux.Handle("GET /api/healthz", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(http.StatusText(http.StatusOK)))
	}))
	mux.HandleFunc("GET /admin/metrics", apiConfig.handlerMetrics)
	mux.HandleFunc("GET /api/reset", apiConfig.handlerReset)

	mux.HandleFunc("POST /api/chirps", apiConfig.handlerCreateChirp)
	mux.HandleFunc("GET /api/chirps", apiConfig.handlerGetChirps)
	mux.HandleFunc("GET /api/chirps/{id}", apiConfig.handlerGetChirp)

	fmt.Printf("listening on port %s...\n", port)
	log.Fatal(server.ListenAndServe())
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	html := `
  <html>
    <body>
      <h1>Welcome, Chirpy Admin</h1>
      <p>Chirpy has been visited %d times!</p>
    </body>
  </html>
  `
	fmt.Fprintf(w, html, cfg.serverHits)
}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	cfg.serverHits = 0
	fmt.Fprint(w, "Hits reset to 0")
}

func middlewareCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
