package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/chirpy/internal/database"
	"github.com/google/uuid"
)

type Chirp struct {
	ID        uuid.UUID       `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body string `json:"body"`
	UserID uuid.UUID  `json:"user_id"`
}

func (cfg *apiConfig) createChirp(w http.ResponseWriter, r *http.Request){
	type parameters struct {
		Body string `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	type response struct {
		Chirp
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	cleaned_sentence, err := validateChirp(params.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body: cleaned_sentence,
		UserID: params.UserID,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp", err)
		return
	}


	respondWithJSON(w, http.StatusCreated, Chirp{
			ID: chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt:chirp.UpdatedAt,
			Body:chirp.Body,
			UserID: chirp.UserID,
		},
	)
}

func validateChirp(body string) (string, error) {
	const maxChirpLength = 140
	if len(body) > maxChirpLength {
		return "", errors.New("Chirp is too long")

	}
	FORBIDDEN_WORDS := [3]string{"kerfuffle", "sharbert", "fornax"}

	words := strings.Fields(body)

	for _, badWord := range FORBIDDEN_WORDS {
		for i := range words {
			if strings.ToLower(words[i]) == badWord {
				words[i] = "****"
			}
		}
	}

	cleaned_sentence := strings.Join(words, " ")

	return cleaned_sentence, nil
}