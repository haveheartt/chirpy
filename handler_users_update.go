package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
    Issuer   string `json:"iss"`
    Subject  string `json:"sub"`
    ExpiresAt int64  `json:"exp"`
    IssuedAt  int64  `json:"iat"`
    jwt.RegisteredClaims
}

func (cfg *apiConfig) handlerUsersUpdate(w http.ResponseWriter, r *http.Request) {
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

    authHeader := r.Header.Get("Authorization")
    bearerToken := strings.Split(authHeader, "Bearer ")
    token := bearerToken[1]

    claims := &jwt.RegisteredClaims{}
    jwtToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
        return []byte(cfg.jwtSecret), nil
    })
    if err != nil {
        log.Fatalf("error parse with claims: %v", err)
    }  

    id := ""
    if claims, ok := jwtToken.Claims.(*jwt.RegisteredClaims); ok && jwtToken.Valid {
        id = claims.Subject
    } 

	user, err := cfg.DB.UpdateUser(id, params.Email, params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't update user")
		return
	}

	respondWithJSON(w, http.StatusOK, User{
		ID:   user.ID,
		Email: user.Email,
	})
}
