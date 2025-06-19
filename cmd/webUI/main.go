package main

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("./internal/web/static"))
	mux.Handle("/", fs)

	log.Printf("Server listening on port %s", port)
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), mux)
	if err != nil {
		log.Error().Msgf("Failed error : %v", err)
	}
}
