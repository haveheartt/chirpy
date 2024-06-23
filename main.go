package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type apiConfig struct {
    fileserverHits int
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits++
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handlerWriteHits(w http.ResponseWriter, r *http.Request) {
    s := fmt.Sprintf("Hits: %v", cfg.fileserverHits)
    w.Write([]byte(s))
}

func (cfg *apiConfig) handlerResetHits(w http.ResponseWriter, r *http.Request) {
       cfg.fileserverHits = 0
}

func main() {
    port := ":8080"
    apiCfg := apiConfig{}
   
    mux := http.NewServeMux()
    mux.Handle("/app/*", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir("./")))))
    mux.HandleFunc("/healthz", handlerReadiness)
    mux.HandleFunc("/metrics", apiCfg.handlerWriteHits)
    mux.HandleFunc("/reset", apiCfg.handlerResetHits)

    s := http.Server{
	    Addr:           port,
	    Handler:        mux,
	    ReadTimeout:    10 * time.Second,
	    WriteTimeout:   10 * time.Second,
	    MaxHeaderBytes: 1 << 20,
    }

    log.Printf("🚀 Server started at: http://localhost%s\n", port)
    log.Fatal(s.ListenAndServe())
}

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
    w.Header().Add("Content-Type", "text/plain; charset=utf-8")    
    w.WriteHeader(200)
    w.Write([]byte("OK"))
}

