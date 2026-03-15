package main

import (
	"net/http"
	"os"
)

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		return "8080"
	}
	return port
}

func listenAndServe(port string) {
	logger.Info("dashboard starting", "port", port)
	if err := http.ListenAndServe(":"+port, logMiddleware(http.DefaultServeMux)); err != nil {
		logger.Error("server failed", "err", err)
		os.Exit(1)
	}
}
