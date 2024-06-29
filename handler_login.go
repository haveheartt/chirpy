package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
    type parameters struct {
		Email string `json:"email"`
        Password string `json:"password"`
    }

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

   
	user, err := cfg.DB.LoginUser(params.Email, params.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't create user")
		return
	}

    token, err := handleJWT(user.ID, cfg)
    if err != nil {
        log.Fatalf("error jwt handle: %v", err)
    }

	respondWithJSON(w, http.StatusOK, User{
		ID:   user.ID,
		Email: user.Email,
	    Token: token,
    })
}

func handleJWT(id int, cfg *apiConfig) (string, error){
    intID := int(id)
    strID := strconv.Itoa(intID)

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
        Issuer: "chirpy",
        IssuedAt: jwt.NewNumericDate(time.Now().Local()),
        ExpiresAt: &jwt.NumericDate{time.Now().Add(time.Duration(10) * time.Second)},
        Subject: strID,
    })

    jwt, err := token.SignedString([]byte(cfg.jwtSecret))
    if err != nil {
        log.Fatalf("jwt token error: %v", err)
    }

    return jwt, nil
}
