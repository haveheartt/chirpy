package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/haveheartt/chirpy/database"
)

type apiConfig struct {
    fileserverHits int
}

type ChirpyDB struct {
    db *database.DB
}

type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

var db *database.DB

func main() {
    port := ":8080"
    apiCfg := apiConfig{
        fileserverHits: 0,
    }

    db, err := database.NewDB("database.json")
    if err != nil {
        log.Fatal(err)
    }

    chirpyDB := ChirpyDB{db: db}

    mux := http.NewServeMux()
    mux.Handle("/app/*", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir("./")))))
    mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
    mux.HandleFunc("GET /api/healthz", handlerReadiness)
    mux.HandleFunc("GET /api/reset", apiCfg.handlerReset)
    mux.HandleFunc("POST /api/chirps", chirpyDB.handlerChirpsCreation)
    mux.HandleFunc("GET /api/chirps", chirpyDB.handlerChirpsGet)

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

func (chirpyDB *ChirpyDB) handlerChirpsGet(w http.ResponseWriter, r *http.Request) {
    dbChirps, err := chirpyDB.db.GetChirps()
    if err != nil {
        log.Fatal(err)
    }

	chirps := []Chirp{}
	for _, dbChirp := range dbChirps {
		chirps = append(chirps, Chirp{
			ID:   dbChirp.ID,
			Body: dbChirp.Body,
		})
	}

	sort.Slice(chirps, func(i, j int) bool {
		return chirps[i].ID < chirps[j].ID
	})

	respondWithJSON(w, http.StatusOK, chirps)
}

func (chirpyDB *ChirpyDB) handlerChirpsCreation(w http.ResponseWriter, r *http.Request) {
 	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

    cleaned := chirpsValidate(params.Body, w)

	chirp, err := chirpyDB.db.CreateChirp(cleaned)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp")
		return
	}

	respondWithJSON(w, http.StatusCreated, Chirp{
		ID:   chirp.ID,
		Body: chirp.Body,
	})
}

func chirpsValidate(body string, w http.ResponseWriter) string {
	const maxChirpLength = 140
	if len(body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return ""
	}
 
    return cleanBody(body)
}

func cleanBody(body string) string {
    str := strings.Split(body, " ")

    for i, word := range str {
        if strings.ToLower(word) == "kerfuffle" {
            str[i] = "****"
        } else if strings.ToLower(word) == "sharbert" {
            str[i] = "****"
        } else if strings.ToLower(word) == "fornax" {
            str[i] = "****"
        } else {
            continue
        }
    }

    return strings.Join(str, " ")
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Printf("Responding with 5XX error: %s", msg)
	}
	type errorResponse struct {
		Error string `json:"error"`
	}
	respondWithJSON(w, code, errorResponse{
		Error: msg,
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(code)
	w.Write(dat)
}
