package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"gogo/internal/user"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(sqlite.Open("gogo.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("open sqlite database: %v", err)
	}

	if err := db.AutoMigrate(&user.User{}); err != nil {
		log.Fatalf("migrate database: %v", err)
	}

	service := user.NewService(db)
	handler := user.NewHandler(service)

	router := gin.Default()
	handler.RegisterRoutes(router)

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("run server: %v", err)
	}
}
