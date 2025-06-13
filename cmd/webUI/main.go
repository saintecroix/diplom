package main

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/saintecroix/diplom/cmd/webUI/internal/web"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not specified
	}

	mux := http.NewServeMux()

	// Serve static files
	fs := http.FileServer(http.Dir("./webUI/cmd/webUI/internal/web/static"))
	mux.Handle("/", fs)

	// Register handlers
	web.RegisterHandlers(mux)

	log.Printf("Server listening on port %s", port)
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), mux)
	if err != nil {
		log.Fatal(err)
	}
}
