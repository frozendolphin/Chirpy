package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) getHits(w http.ResponseWriter, req *http.Request){
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	hits := fmt.Sprintf(`<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`, cfg.fileserverHits.Load())
	w.Write([]byte(hits))
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        cfg.fileserverHits.Add(1) 
        next.ServeHTTP(w, r)      
    })
}