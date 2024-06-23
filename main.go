package main

import (
	"net/http"
	"time"
)

func main() {
    mux := http.NewServeMux()
    s := http.Server{
	    Addr:           ":8080",
	    Handler:        mux,
	    ReadTimeout:    10 * time.Second,
	    WriteTimeout:   10 * time.Second,
	    MaxHeaderBytes: 1 << 20,
    }
    s.ListenAndServe()
}

