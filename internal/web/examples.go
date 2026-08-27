package web

import (
	"encoding/json"
	"net/http"
)

func (s *Server) handleExamples(w http.ResponseWriter, r *http.Request) {
	if !requireMethod(w, r, http.MethodGet) {
		return
	}
	if len(s.assets.Examples) == 0 {
		writeJSON(w, http.StatusOK, map[string]string{})
		return
	}
	out := make(map[string]json.RawMessage, len(s.assets.Examples))
	for name, payload := range s.assets.Examples {
		out[name] = payload
	}
	writeJSON(w, http.StatusOK, out)
}
