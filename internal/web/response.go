package web

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
)

const maxBodyBytes = 1 << 20

// writeJSON encodes v as JSON with the given status code.
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("encode response: %v", err)
	}
}

// writeError emits a JSON error body whose message comes from the
// backend error chain, so the page always displays the real failure.
func writeError(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, ErrorResponse{Error: err.Error()})
}

// badRequest is the shorthand used for solver and decoding failures.
func badRequest(w http.ResponseWriter, err error) {
	writeError(w, http.StatusBadRequest, err)
}

// decodeBody reads and parses the request body with a size limit.
func decodeBody(r *http.Request, v any) error {
	if r.Body == nil {
		return errors.New("empty request body")
	}
	defer r.Body.Close()
	if r.ContentLength > maxBodyBytes {
		return errors.New("request body too large")
	}
	data, err := io.ReadAll(io.LimitReader(r.Body, maxBodyBytes+1))
	if err != nil {
		return errors.New("read body: " + err.Error())
	}
	if int64(len(data)) > maxBodyBytes {
		return errors.New("request body too large")
	}
	if len(data) == 0 {
		return errors.New("empty request body")
	}
	if err := json.Unmarshal(data, v); err != nil {
		return errors.New("invalid JSON: " + err.Error())
	}
	return nil
}

// requireMethod rejects requests that do not use the expected verb with a
// 405 JSON error body.
func requireMethod(w http.ResponseWriter, r *http.Request, method string) bool {
	if r.Method == method {
		return true
	}
	writeJSON(w, http.StatusMethodNotAllowed, ErrorResponse{
		Error: "method " + r.Method + " not allowed on this endpoint",
	})
	return false
}
