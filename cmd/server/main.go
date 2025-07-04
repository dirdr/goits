package main

import (
	"github.com/dirdr/goits/internal/handler"
	"github.com/dirdr/goits/pkg/logger"
)

func main() {
	log := logger.New("debug")

	log.Info("Starting server...")

	r := handler.NewRouter()

	if err := r.Run(":8080"); err != nil {
		log.Error("Failed to start server", "error", err)
	}
}
