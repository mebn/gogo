package main

import (
	"log"

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

	userService := user.NewService(db)
	userHandler := user.NewHandler(userService)
	userHandler.RegisterRoutes(router)

	petService := pet.NewService(db)
	petHandler := pet.NewHandler(petService)
	petHandler.RegisterRoutes(router)

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("run server: %v", err)
	}
}
