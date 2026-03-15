package main

import (
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// gzipFileServer wraps http.FileServer to serve pre-compressed .gz files when
// the client accepts gzip and a .gz sibling exists on disk.
func gzipFileServer(dir string) http.Handler {
	fs := http.FileServer(http.Dir(dir))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			fs.ServeHTTP(w, r)
			return
		}
		p := filepath.Join(dir, filepath.Clean("/"+r.URL.Path))
		if _, err := os.Stat(p + ".gz"); err == nil {
			r2 := r.Clone(r.Context())
			r2.URL.Path += ".gz"
			w.Header().Set("Content-Encoding", "gzip")
			w.Header().Set("Vary", "Accept-Encoding")
			// Strip Content-Type sniffing from the .gz extension
			ext := filepath.Ext(p)
			if ct := mime.TypeByExtension(ext); ct != "" {
				w.Header().Set("Content-Type", ct)
			}
			fs.ServeHTTP(w, r2)
			return
		}
		fs.ServeHTTP(w, r)
	})
}
