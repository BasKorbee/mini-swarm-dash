//go:build !mock

package main

import (
	"net/http"

	"github.com/moby/moby/client"
)

func main() {
	// Connect to the local Docker socket (/var/run/docker.sock must be mounted)
	cli, err := client.New(client.FromEnv)
	if err != nil {
		logger.Error("failed to create Docker client", "err", err)
		return
	}

	http.HandleFunc("/api/services", handleServices(cli))
	http.HandleFunc("/api/nodes", handleNodes(cli))
	http.HandleFunc("/api/local-stats", handleLocalStats(cli))
	http.Handle("/", gzipFileServer("public"))

	listenAndServe(getPort())
}
