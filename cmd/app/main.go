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
	//  load config
	cfg := configs.LoadConfig()

	// connect database
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

	//Auth Handler
	userRepo := repository.NewUserRepository(db)
	userService := usecases.NewUserService(userRepo)
	authHandler := handler.NewAuthHandler(userService)
	//User Management Handler
	userManagementRepo := repository.NewUserManagementRepository(db)
	userManagementService := usecases.NewUserManagementService(userManagementRepo)
	userManagementHandler := handler.NewUserManagementHandler(userManagementService)

	// router
	r := router.SetupRouter(authHandler, userManagementHandler)

	// run server
	port := fmt.Sprintf(":%s", cfg.AppPort)
	log.Printf("Server running on port %s (env: %s)", cfg.AppPort, cfg.AppEnv)

	r.Run(port)
}
