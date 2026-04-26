package server

import (
	"encoding/json"
	"net/http"
)



type HealthResponse struct{
	Status string `json:"status"`
	Service string `json:"service"`
}

func HeathHandler(w http.ResponseWriter, r *http.Request) {
	
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
     response := HealthResponse{
     Status: "ok",
     Service: "aegis-ai-gateway",
     }	
	_ = json.NewEncoder(w).Encode(response)
	
}