package main

import (
    "encoding/json"
	"net/http"
	"strings"
)

func (cfg *apiConfig) handlerUsersUpgrade(w http.ResponseWriter, r *http.Request) {	
    type UserData struct {
        UserID int `json:"user_id"`
    }

    type parameters struct {
		Event string `json:"event"`
		Data  UserData `json:"data"`
	}

	keyHeader := r.Header.Get("Authorization")
	if keyHeader == "" {
	    respondWithError(w, http.StatusUnauthorized, "unauthorized")
	}
	splitKey := strings.Split(keyHeader, " ")
	if len(keyHeader) < 2 || splitKey[0] != "ApiKey" {
        respondWithError(w, http.StatusBadRequest, "malformed authorization header")
	}

    decoder := json.NewDecoder(r.Body)
	params := parameters{}
    err := decoder.Decode(&params)
	if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

    if params.Event == "user.upgraded" {
        err = cfg.DB.UpdateUserMembership(params.Data.UserID)
	    if err != nil {
		    respondWithError(w, http.StatusNotFound, "Couldn't find user")
		    return
	    }
	    respondWithJSON(w, http.StatusNoContent, nil)
    }
    respondWithError(w, http.StatusNoContent, "")
}
