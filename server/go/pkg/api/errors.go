package api

import (
	"encoding/json"
	"net/http"

	"github.com/segolab/relay-ref/server/go/pkg/middleware"
	"github.com/segolab/relay-ref/server/go/pkg/model"
)

func WriteError(w http.ResponseWriter, r *http.Request, status int, code, message string, details map[string]any) {
	reqID := middleware.RequestIDFromContext(r.Context())
	resp := model.ErrorResponse{
		Code:      code,
		Message:   message,
		Details:   details,
		RequestID: reqID,
	}

	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(resp)
}
