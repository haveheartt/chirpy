package main

import (
	"net/http"
	"strconv"
)

func (cfg *apiConfig) handlerChirpFind(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	chirp, err := cfg.DB.FindChirp(id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't find chirp")
		return
	}

	respondWithJSON(w, http.StatusCreated, Chirp{
		ID:   chirp.ID,
		Body: chirp.Body,
	})
}
