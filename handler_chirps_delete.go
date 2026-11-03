package main

import (
	"net/http"

	"github.com/connor317/bootdev-chirpy/internal/auth"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error": "Authorization header is missing"}`))
		return
	}
	userId, err := auth.ValidateJWT(bearerToken, cfg.jwtSecret)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error": "Invalid authorization token"}`))
		return
	}

	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Invalid chirp ID"}`))
		return
	}

	chirp, err := cfg.dbQueries.GetChirpByID(r.Context(), chirpID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "Chirp not found"}`))
		return
	}

	if chirp.UserID != userId {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`{"error": "You are not authorized to delete this chirp"}`))
		return
	}

	err = cfg.dbQueries.DeleteChirpByID(r.Context(), chirpID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Failed to delete chirp"}`))
		return
	}

	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte(`{"message": "Chirp deleted successfully"}`))
}
