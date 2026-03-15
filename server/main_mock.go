//go:build mock

package main

import (
	"log"
	"net/http"
)

func main() {
	log.Println("Mock mode enabled — using fake swarm data")

	http.HandleFunc("/api/services", mockHandleServices())
	http.HandleFunc("/api/nodes", mockHandleNodes())
	http.HandleFunc("/api/local-stats", mockHandleLocalStats())
	http.Handle("/", gzipFileServer("client/dist"))

	log.Println("Swarm dashboard running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
