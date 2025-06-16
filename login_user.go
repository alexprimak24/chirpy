package main

import (
	"encoding/json"
	"net/http"

	"github.com/chirpy/internal/auth"
)

func (cfg *apiConfig) loginUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}

	type response struct {
		User
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	
	
	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user", err)
		return
	}
	
	validPassErr := auth.CheckPasswordHash(user.HashedPassword,params.Password)

	if validPassErr != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID: user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email: user.Email,
		},
	})
}
