package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func handlerChirpsValidate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type returnVals struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}
	FORBIDDEN_WORDS := [3]string{"kerfuffle", "sharbert", "fornax"}

	words := strings.Fields(params.Body)

	for _, badWord := range FORBIDDEN_WORDS {
		for i := range words {
			if strings.ToLower(words[i]) == badWord {
				words[i] = "****"
			}
		}
	}

	cleaned_sentence := strings.Join(words, " ")

	respondWithJSON(w, http.StatusOK, returnVals{
		CleanedBody: cleaned_sentence,
	})
}
