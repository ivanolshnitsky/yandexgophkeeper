package server

import (
	"encoding/json"
	"net/http"
)

type PingResponse struct {
	Status string `json:"status"`
}

func (a *App) handlePing(w http.ResponseWriter, r *http.Request) {
	resp := PingResponse{Status: "ok"}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}
