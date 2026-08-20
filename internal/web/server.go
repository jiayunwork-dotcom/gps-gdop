package web

import (
	"io/fs"
	"net/http"
)

// Assets bundles the embedded static files and example payloads the
// server exposes. The examples are served verbatim so the page can load
// them and post them straight back to the solver endpoints.
type Assets struct {
	WebFS    fs.FS
	Examples map[string][]byte
}

// Server wires the solver API and the static page onto one mux.
type Server struct {
	assets Assets
	mux    *http.ServeMux
}

// NewServer builds the HTTP handler with the solver endpoints and the
// static file handler for the web console.
func NewServer(assets Assets) http.Handler {
	s := &Server{assets: assets, mux: http.NewServeMux()}
	s.mux.HandleFunc("/api/dop", s.handleDop)
	s.mux.HandleFunc("/api/sky", s.handleSky)
	s.mux.HandleFunc("/api/examples", s.handleExamples)
	s.mux.Handle("/", staticHandler(assets.WebFS))
	return withRequestLog(s.mux)
}

// ServeHTTP dispatches to the registered routes.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}
