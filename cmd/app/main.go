package main

import (
	"fmt"
	"log"

	"be-knowledge/configs"
	"be-knowledge/internal/delivery/http/handler"
	"be-knowledge/internal/delivery/http/router"
	"be-knowledge/internal/repository"
	"be-knowledge/internal/usecases"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func main() {
	// 1. load config
	cfg := configs.LoadConfig()

	// 2. connect database
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true",
		cfg.DBUser,
		cfg.DBPass,
		cfg.DBHost,
		cfg.DBName,
	)

	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		log.Fatal("Cannot connect DB:", err)
	}

	// 3. repository
	userRepo := repository.NewUserRepository(db)

	// 4. service
	userService := usecases.NewUserService(userRepo)

	// 5. handler
	authHandler := handler.NewAuthHandler(userService)

	// 6. router
	r := router.SetupRouter(authHandler)

	// 7. run server
	port := fmt.Sprintf(":%s", cfg.AppPort)
	log.Printf("Server running on port %s (env: %s)", cfg.AppPort, cfg.AppEnv)

	r.Run(port)
}
