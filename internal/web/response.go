package web

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
)

const maxBodyBytes = 1 << 20

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("encode response: %v", err)
	}
}

func writeError(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, ErrorResponse{Error: err.Error()})
}

func badRequest(w http.ResponseWriter, err error) {
	writeError(w, http.StatusBadRequest, err)
}

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

func requireMethod(w http.ResponseWriter, r *http.Request, method string) bool {
	if r.Method == method {
		return true
	}
	writeJSON(w, http.StatusMethodNotAllowed, ErrorResponse{
		Error: "method " + r.Method + " not allowed on this endpoint",
	})
	return false
}
