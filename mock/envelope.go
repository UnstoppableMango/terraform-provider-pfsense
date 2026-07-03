package mock

import (
	"encoding/json"
	"net/http"
)

type apiResponse struct {
	Code       int    `json:"code"`
	Status     string `json:"status"`
	ResponseID string `json:"response_id"`
	Message    string `json:"message"`
	Data       any    `json:"data,omitempty"`
}

func writeResponse(w http.ResponseWriter, code int, data any) {
	status := "ok"
	if code >= 400 {
		status = "error"
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(apiResponse{
		Code:       code,
		Status:     status,
		ResponseID: "mock-response",
		Message:    http.StatusText(code),
		Data:       data,
	})
}
