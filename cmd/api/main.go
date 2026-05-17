package main

import (
	"log"

	"reflect-backend/internal/config"
	"reflect-backend/internal/database"
	"reflect-backend/internal/middleware"
	"reflect-backend/internal/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

	database.ConnectDatabase(cfg)
	database.AutoMigrate()

	router := gin.Default()

	router.Use(middleware.ErrorMiddleware())

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "success",
			"message": "Reflect backend API is running",
		})
	})

	router.Use(middleware.ErrorMiddleware())

	routes.SetupRoutes(router,cfg)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "success",
			"message": "Server is healthy",
		})
	})

	port := ":" + cfg.AppPort

	log.Println("Server running on port", port)

	if err := router.Run(port); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}