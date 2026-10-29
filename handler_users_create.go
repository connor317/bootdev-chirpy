package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/connor317/bootdev-chirpy/internal/auth"
	"github.com/connor317/bootdev-chirpy/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type createUserRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req createUserRequest
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

	user, err := cfg.dbQueries.CreateUser(r.Context(), database.CreateUserParams{
		Email:    req.Email,
		Password: hashedPassword,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("%s", err)
		w.Write([]byte(`{"error": "Failed to create user"}`))
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	formatted_user := User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}
	respBody, _ := json.Marshal(formatted_user)
	w.Write(respBody)
}
