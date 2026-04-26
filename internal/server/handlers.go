package server

import (
	apphttp "aegis-ai-gateway/internal/app_http"
	"aegis-ai-gateway/internal/domain"
	"aegis-ai-gateway/internal/provider"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type Handler interface {
	ServerHttp(http.ResponseWriter, *http.Request)
}

type ChatHandler struct {
	logger *slog.Logger
	provider provider.LLMProvider
	requestTimeout time.Duration
}

type LogRequest struct {
	RequestID uuid.UUID
	Method    string
	Path      string
}

func (lr LogRequest) String() string {
	return "RequestID: " + string(lr.RequestID.String()) + " Method: " + lr.Method + " Path: " + lr.Path
}

func (h *ChatHandler) ServeHttp(w http.ResponseWriter, r *http.Request) {
	
	ctx, cancel := context.WithTimeout(r.Context(), h.requestTimeout)
	defer cancel()
	
	if r.Method != http.MethodPost {
		
		writeAPIError(w, http.StatusMethodNotAllowed, "405", "Method not allowed")
		h.logger.Error("method not allowed")
		return
	}
	
	if r.Body == nil {
		writeAPIError(w, http.StatusBadRequest, "400", "Request body is empty")
		h.logger.Error("request body is empty")
		return
	}
	
	var requestID uuid.UUID
	if r.Header.Get("X-Request-ID") != "" {
		requestId, err := uuid.Parse(r.Header.Get("X-Request-ID"))
		if err != nil {
			
			requestID = generateRequestID()
			h.logger.Debug("Generate Request ID: " + string(requestID.String()))
		} else {
			requestID = requestId
			
		}
	} else {
		requestID = generateRequestID()
		h.logger.Debug("Generate Request ID: " + string(requestID.String()))
	}
	
	h.logger.InfoContext(ctx, LogRequest{
		RequestID: requestID,
		Method:    r.Method,
		Path:      r.URL.Path,
	}.String())
	
	bodyData, err := io.ReadAll(r.Body)
	if err != nil {
		writeAPIError(w, http.StatusBadRequest, "500", "Invalid request body")
		h.logger.Error("invalid request body", "request_id", requestID)
		return
	}
	var data apphttp.RequestDTO
	if err := json.Unmarshal(bodyData, &data); err != nil {
		writeAPIError(w, http.StatusBadRequest, "400", "Invalid request body")
		return
	}
	defer r.Body.Close()
	providerReq := domain.GenerateRequest{
		RequestID: requestID.String(),
		Model:     data.Model,
		Messages:  []domain.Message{data.Message},
		ResponseFormat: data.ResponseFormat,
		Metadata:      data.Metadata,
	}
	providerResponse, err := h.provider.Generate(ctx, providerReq)
	if err != nil {
		writeAPIError(w, http.StatusInternalServerError, "500", "LLM Failed to do its job :(")
		return
	}
	
	responseDTO := apphttp.ResponseDTO {
		ID:      requestID,
		Status:  domain.StatusCompleted,
		Provider: providerResponse.Provider,
		Model:   providerResponse.Model,
		Output:  providerResponse.Output,
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(responseDTO); err != nil {
		writeAPIError(w, http.StatusInternalServerError, "500", "Failed to encode data")
	}
	
	h.logger.Info("request processed", "request_id", requestID)
	
}

func generateRequestID () uuid.UUID {
	return uuid.New()
}

func writeAPIError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	resp := domain.APIErrorResponse{
		Error: domain.APIError{
			Code:    code,
			Message: message,
		},
	}
	_ = json.NewEncoder(w).Encode(resp)
}