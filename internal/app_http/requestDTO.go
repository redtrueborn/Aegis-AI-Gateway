package apphttp

import "aegis-ai-gateway/internal/domain"


type RequestDTO struct {
	Model string `json:"model"`
	Message domain.Message `json:"message"`
	ResponseFormat string `json:"response_format"`
	Metadata map[string]string `json:"metadata"`
}

