package main

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) getChirp(w http.ResponseWriter, r *http.Request){
	chirpIdStr := r.PathValue("chirpID")

	// convert from string to UUID
	chirpId, err := uuid.Parse(chirpIdStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID format", err)
	}


	dbChirp, err := cfg.db.ReadChirp(r.Context(),chirpId)

	if errors.Is(err, sql.ErrNoRows) {
		respondWithError(w, http.StatusNotFound, "Chirp not found", err)
		return
	} else if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't read chirp", err)
		return
	}

	chirpEntry := Chirp{
		ID: dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		UserID:    dbChirp.UserID,
		Body:      dbChirp.Body,
	}

	respondWithJSON(w, http.StatusOK, chirpEntry)
}