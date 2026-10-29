package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
	if cfg.platform == "dev" {
		err := cfg.dbQueries.Reset(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "Failed to reset database"}`))
			return
		}
	} else {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`{"error": "Reset is only allowed in dev environment"}`))
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte(fmt.Sprintf("Hits: %d", cfg.fileserverHits.Load())))
}
