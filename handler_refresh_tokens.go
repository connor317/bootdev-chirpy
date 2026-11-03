package main

import (
	"net/http"
	"time"

	"encoding/json"

	"github.com/connor317/bootdev-chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefreshToken(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error": "Authorization header is missing"}`))
		return
	}

	entry, err := cfg.dbQueries.GetRefreshTokenByToken(r.Context(), refreshToken)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error": "Invalid refresh token"}`))
		return
	}

	if entry.RevokedAt.Valid || entry.ExpiresAt.Before(time.Now()) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error": "Refresh token has expired"}`))
		return
	}

	token, err := auth.MakeJWT(entry.UserID, cfg.jwtSecret, time.Hour)
	type RespBody struct {
		Token string `json:"token"`
	}

	respBody := RespBody{
		Token: token,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	dat, err := json.Marshal(respBody)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Failed to marshal response"}`))
		return
	}
	w.Write(dat)
}

func (cfg *apiConfig) handlerRevokeRefreshToken(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error": "Authorization header is missing"}`))
		return
	}

	err = cfg.dbQueries.RevokeRefreshToken(r.Context(), refreshToken)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error": "Invalid refresh token"}`))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
