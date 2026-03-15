package main

import (
	"log"
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
	log.Println("Dashboard running on :" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
