package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

const (
	MediaTypeJSON = "application/json"
	MaxBodySize   = 1024 * 1024 // 1MB
)

// JSONResponse is a generic JSON response.
type JSONResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// WriteJSON writes a JSON response.
// It sets the provided HTTP headers, defaulting to "Content-Type: application/json".
func WriteJSON(w http.ResponseWriter, status int, data any, headers ...http.Header) error {
	out, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	w.Header().Set("Content-Type", MediaTypeJSON)
	if len(headers) > 0 {
		for k, v := range headers[0] {
			w.Header().Set(k, v[0])
		}
	}

	w.WriteHeader(status)
	if _, err := w.Write(out); err != nil {
		return fmt.Errorf("failed to write response: %w", err)
	}

	return nil
}

// ReadJSON reads a JSON request with a maximum size of 1MB. It disallows unknown fields.
func ReadJSON(w http.ResponseWriter, r *http.Request, data any) error {
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodySize)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(data); err != nil {
		return fmt.Errorf("failed to decode JSON: %w", err)
	}

	if err := dec.Decode(&struct{}{}); err != io.EOF {
		return errors.New("body must only contain a single JSON value")
	}

	return nil
}

// ErrorJSON writes an error JSON response.
// The default status is 400, but can be overridden by providing an optional status code as a second argument.
func ErrorJSON(w http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest
	if len(status) > 0 {
		statusCode = status[0]
	}

	payload := JSONResponse{
		Error:   true,
		Message: err.Error(),
	}

	return WriteJSON(w, statusCode, payload)
}
