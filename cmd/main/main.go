package main

import (
	"log"

	"gogo/internal/auth"
	"gogo/internal/database"
	"gogo/internal/pet"
	"gogo/internal/user"

	"github.com/gin-gonic/gin"
)

func main() {
	db, err := database.InitSQLite()
	if err != nil {
		log.Fatalf("initialize SQLite database: %v", err)
	}

	err = database.Migrate(db)
	if err != nil {
		log.Fatalf("migrate database: %v", err)
	}

	router := gin.Default()
	apiRoute := router.Group("/api/v1")

	authService := auth.NewService(db)
	authHandler := auth.NewHandler(authService)
	authHandler.RegisterRoutes(apiRoute)

	protectedApiRoute := apiRoute.Group("")
	protectedApiRoute.Use(authService.RequireAuthenticatedUser())

	userService := user.NewService(db)
	userHandler := user.NewHandler(userService)
	userHandler.RegisterRoutes(protectedApiRoute)

	petService := pet.NewService(db)
	petHandler := pet.NewHandler(petService)
	petHandler.RegisterRoutes(protectedApiRoute)

	if err := router.Run(":1337"); err != nil {
		log.Fatalf("run server: %v", err)
	}
}
