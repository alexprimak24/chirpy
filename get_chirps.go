package main

import (
	"net/http"

)

func (cfg *apiConfig) getChirps(w http.ResponseWriter, r *http.Request){
	dbChirps, err := cfg.db.ReadChirps(r.Context())

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't read chirp", err)
		return
	}

	chirps := []Chirp{}

	for _, dbChirp := range dbChirps {
		chirps = append(chirps, Chirp{
			ID: dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			UserID:    dbChirp.UserID,
			Body:      dbChirp.Body,
		})
	}
	
	respondWithJSON(w, http.StatusOK, chirps)
}