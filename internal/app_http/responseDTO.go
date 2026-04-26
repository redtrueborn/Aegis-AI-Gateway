package apphttp

import (
	"aegis-ai-gateway/internal/domain"

	"github.com/google/uuid"
)

type ResponseDTO struct {
	ID      uuid.UUID        `json:"id"`
	Status  domain.Status          `json:"status"`
	Provider domain.ProvideName        `json:"provider"`
	Model   string          `json:"model"`
	Output  map[string]any `json:"output"`
	
}

