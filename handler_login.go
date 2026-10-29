package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/connor317/bootdev-chirpy/internal/auth"
	"github.com/google/uuid"
)

type UserWithToken struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Token     string    `json:"token"`
}

func (cfg *apiConfig) handlerLoginUser(w http.ResponseWriter, r *http.Request) {
	type loginUserRequest struct {
		Email     string         `json:"email"`
		Password  string         `json:"password"`
		ExpiresIn *time.Duration `json:"expires_in,omitempty"`
	}
	var req loginUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Invalid request body"}`))
		return
	}

	if req.ExpiresIn == nil || *req.ExpiresIn <= 0 || *req.ExpiresIn > time.Hour {
		defaultExpiresIn := time.Hour
		req.ExpiresIn = &defaultExpiresIn
	}

	user, err := cfg.dbQueries.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Failed to retrieve user password"}`))
		return
	}

	valid, _ := auth.CheckPasswordHash(req.Password, user.Password)

	if !valid {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error": "Invalid email or password"}`))
		return
	}

	token, err := auth.MakeJWT(user.ID, cfg.jwtSecret, *req.ExpiresIn)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Failed to generate JWT"}`))
		return
	}

	respBody := UserWithToken{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
		Token:     token,
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
