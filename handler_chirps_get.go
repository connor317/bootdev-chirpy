package main

import (
	"encoding/json"
	"net/http"
	"sort"

	"github.com/connor317/bootdev-chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetAllChirps(w http.ResponseWriter, r *http.Request) {
	author_id := r.URL.Query().Get("author_id")
	sortDirection := r.URL.Query().Get("sort")
	if sortDirection != "" && sortDirection != "asc" && sortDirection != "desc" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Invalid sort parameter"}`))
		return
	}
	var chirps []database.Chirp
	var err error
	if author_id != "" {
		chirps, err = cfg.dbQueries.GetChirpsByUserID(r.Context(), uuid.MustParse(author_id))
	} else {
		chirps, err = cfg.dbQueries.GetAllChirps(r.Context())
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Failed to retrieve chirps"}`))
		return
	}

	if sortDirection == "desc" {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].CreatedAt.After(chirps[j].CreatedAt)
		})
	} else {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].CreatedAt.Before(chirps[j].CreatedAt)
		})
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	respBody := make([]storedChirp, 0, len(chirps))
	for _, chirp := range chirps {
		respBody = append(respBody, storedChirp{
			Body:      chirp.Body,
			UserId:    chirp.UserID,
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
		})
	}

	dat, err := json.Marshal(respBody)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Something went wrong"}`))
		return
	}
	w.Write(dat)
}

func (cfg *apiConfig) handlerGetChirpByID(w http.ResponseWriter, r *http.Request) {
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

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	respBody := storedChirp{
		Body:      chirp.Body,
		UserId:    chirp.UserID,
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
	}
	dat, err := json.Marshal(respBody)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Something went wrong"}`))
		return
	}
	w.Write(dat)
}
