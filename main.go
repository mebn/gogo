package main

import (
	"log"

	"gogo/internal/pet"
	"gogo/internal/user"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(sqlite.Open("gogo.db?_foreign_keys=on"), &gorm.Config{})
	if err != nil {
		log.Fatalf("open sqlite database: %v", err)
	}

	err = db.AutoMigrate(
		&user.User{},
		&pet.Pet{},
	)
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
