package main

import (
	"log"
	"net/http"
	"time"
)

type apiConfig struct {
    fileserverHits int
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
    

    return next
}

func main() {
    port := ":8080"
    apiCfg := apiConfig{}
   
    mux := http.NewServeMux()
    mux.Handle("/app/*", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir("./")))))
    mux.HandleFunc("/healthz", handlerReadiness)

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
