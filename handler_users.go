package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/connor317/bootdev-chirpy/internal/auth"
	"github.com/connor317/bootdev-chirpy/internal/database"
)

func (cfg *apiConfig) handlerUpdateUsernameAndPassword(w http.ResponseWriter, r *http.Request) {
	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error": "Authorization header is missing"}`))
		return
	}
	userID, err := auth.ValidateJWT(bearerToken, cfg.jwtSecret)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error": "Invalid authorization token"}`))
		return
	}

	type updateUserRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req updateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Invalid request body"}`))
		return
	}

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Failed to hash password"}`))
		return
	}

	updatedUser, err := cfg.dbQueries.UpdateUserEmailPassword(r.Context(), database.UpdateUserEmailPasswordParams{
		ID:       userID,
		Email:    req.Email,
		Password: hashedPassword,
	})

	log.Printf("Updated user: %v", updatedUser)

	response := struct {
		ID          string `json:"id"`
		Email       string `json:"email"`
		CreatedAt   string `json:"created_at"`
		UpdatedAt   string `json:"updated_at"`
		IsChirpyRed bool   `json:"is_chirpy_red"`
	}{
		ID:          updatedUser.ID.String(),
		Email:       updatedUser.Email,
		CreatedAt:   updatedUser.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   updatedUser.UpdatedAt.Format("2006-01-02 15:04:05"),
		IsChirpyRed: updatedUser.IsChirpyRed,
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Failed to update user"}`))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	dat, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Failed to marshal response"}`))
		return
	}
	w.Write(dat)
}
