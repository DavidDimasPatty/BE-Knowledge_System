package main

import (
	"fmt"
	"log"

	"example/be-knowledge/configs"
	httpServer "example/be-knowledge/internal/delivery/http/router"
)

func main() {
	cfg := configs.LoadConfig()

	r := httpServer.NewRouter()

	port := fmt.Sprintf(":%s", cfg.AppPort)

	log.Printf("Server running on port %s (env: %s)", cfg.AppPort, cfg.AppEnv)
	r.Run(port)
}