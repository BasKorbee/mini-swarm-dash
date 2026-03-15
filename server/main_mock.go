//go:build mock

package main

import (
	"net/http"
)

func main() {
	logger.Info("mock mode enabled — using fake swarm data")

	http.HandleFunc("/api/services", mockHandleServices())
	http.HandleFunc("/api/nodes", mockHandleNodes())
	http.HandleFunc("/api/local-stats", mockHandleLocalStats())
	http.Handle("/", gzipFileServer("client/dist"))

	listenAndServe(getPort())
}
