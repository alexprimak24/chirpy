package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/chirpy/internal/auth"
	"github.com/chirpy/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID       `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Password string `json:"-"`
}

func (cfg *apiConfig) createUser(w http.ResponseWriter, r *http.Request){
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

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password", err)
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, response{
		User: User{
			ID: user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email: user.Email,
		},
	})

}