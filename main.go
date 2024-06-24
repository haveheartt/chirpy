package main

import (
	"log"
	"net/http"
	"time"
)

type apiConfig struct {
    fileserverHits int
}

func main() {
    port := ":8080"
    apiCfg := apiConfig{
        fileserverHits: 0,
    }
   
    mux := http.NewServeMux()
    mux.Handle("/app/*", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir("./")))))
    mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
    mux.HandleFunc("GET /api/healthz", handlerReadiness)
    mux.HandleFunc("/api/reset", apiCfg.handlerReset)

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

