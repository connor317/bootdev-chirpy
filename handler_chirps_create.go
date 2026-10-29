package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/connor317/bootdev-chirpy/internal/auth"
	"github.com/connor317/bootdev-chirpy/internal/database"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type storedChirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserId    uuid.UUID `json:"user_id"`
}

type chirpRequest struct {
	Body   string    `json:"body"`
	UserId uuid.UUID `json:"user_id"`
	Token  string    `json:"token"`
}

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	type errorResponse struct {
		Error string `json:"error"`
	}

	type successResponse struct {
		CleanedBody string `json:"cleaned_body"`
	}
	var c chirpRequest

	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		respBody := errorResponse{Error: "Something went wrong"}
		dat, err := json.Marshal(respBody)

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "Something went wrong"}`))
			return
		}
		w.Write(dat)
		return
	}

	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		respBody := errorResponse{Error: "Authorization header is missing"}
		dat, err := json.Marshal(respBody)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "Something went wrong"}`))
			return
		}
		w.Write(dat)
		return
	}

	userUUID, err := auth.ValidateJWT(bearerToken, cfg.jwtSecret)
	log.Printf("userUUID: %v, c.UserId: %v", userUUID, c.UserId)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		respBody := errorResponse{Error: "Invalid or expired token"}
		dat, err := json.Marshal(respBody)

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "Something went wrong"}`))
			return
		}
		w.Write(dat)
		return
	}

	if len(c.Body) > 140 {
		w.WriteHeader(http.StatusBadRequest)
		respBody := errorResponse{Error: "Chirp is too long"}
		dat, err := json.Marshal(respBody)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "Something went wrong"}`))
			return
		}
		w.Write(dat)
		return
	}
	words := strings.Fields(c.Body)
	for i, word := range words {
		if strings.ToLower(word) == "kerfuffle" || strings.ToLower(word) == "sharbert" || strings.ToLower(word) == "fornax" {
			words[i] = "****"
		}
	}

	respChirp, err := cfg.dbQueries.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   strings.Join(words, " "),
		UserID: userUUID,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Failed to create chirp"}`))
		return
	}

	w.WriteHeader(http.StatusCreated)
	respBody := storedChirp{
		Body:      respChirp.Body,
		UserId:    respChirp.UserID,
		ID:        respChirp.ID,
		CreatedAt: respChirp.CreatedAt,
		UpdatedAt: respChirp.UpdatedAt,
	}
	dat, err := json.Marshal(respBody)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Something went wrong"}`))
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(dat)
}
