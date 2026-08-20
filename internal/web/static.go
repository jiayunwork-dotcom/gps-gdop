package web

import (
	"io"
	"io/fs"
	"net/http"
	"path"
	"strings"
)

// staticHandler serves the embedded web console from the web/ subtree.
// Unknown paths fall back to index.html so the page can be reached at the
// root of the server.
func staticHandler(webFS fs.FS) http.Handler {
	if webFS == nil {
		return http.NotFoundHandler()
	}
	sub, err := fs.Sub(webFS, "web")
	if err != nil {
		return http.NotFoundHandler()
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clean := path.Clean(strings.TrimPrefix(r.URL.Path, "/"))
		if clean == "." || clean == "/" {
			clean = "index.html"
		}
		if _, err := fs.Stat(sub, clean); err != nil {
			clean = "index.html"
		}
		serveFile(w, r, sub, clean)
	})
}

// serveFile streams one file from the embedded filesystem.
func serveFile(w http.ResponseWriter, r *http.Request, root fs.FS, name string) {
	f, err := root.Open(name)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil {
		http.NotFound(w, r)
		return
	}
	rs, ok := f.(io.ReadSeeker)
	if !ok {
		http.NotFound(w, r)
		return
	}
	http.ServeContent(w, r, info.Name(), info.ModTime(), rs)
}
