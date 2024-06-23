package main

import (
	"log"
	"net/http"
	"time"
)

func main() {
    mux := http.NewServeMux()

    port := ":8080"

    s := http.Server{
	    Addr:           port,
	    Handler:        mux,
	    ReadTimeout:    10 * time.Second,
	    WriteTimeout:   10 * time.Second,
	    MaxHeaderBytes: 1 << 20,
    }

    mux.Handle("/app", http.StripPrefix("/app", http.FileServer(http.Dir("./"))))
    mux.Handle("/app/assets", http.StripPrefix("/app/assets", http.FileServer(http.Dir("./assets"))))
    mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Add("Content-Type", "text/plain; charset=utf-8")    
        w.WriteHeader(200)
        w.Write([]byte("OK"))
    })

    log.Printf("🚀 Server started at: http://localhost%s\n", port)
    log.Fatal(s.ListenAndServe())
}

