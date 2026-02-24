package main

import (
	"errors"
	"log"
	"net/http"
	"os"

	"splitflap-api-go/internal/handler"
	"splitflap-api-go/internal/middleware"
	"splitflap-api-go/internal/service"
)

func main() {
	mux := http.NewServeMux()
	displayService := service.NewDisplayService()
	displayHandler := handler.NewDisplayHandler(displayService)

	mux.Handle("/api/v1/displays/", displayHandler)

	addr := ":8080"
	if port := os.Getenv("PORT"); port != "" {
		addr = ":" + port
	}

	server := &http.Server{
		Addr:    addr,
		Handler: middleware.WithCORS(mux),
	}

	log.Printf("Splitflap API listening on %s", addr)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server failed: %v", err)
	}
}
