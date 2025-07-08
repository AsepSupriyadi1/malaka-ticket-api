package main

import (
	"case_study_api/config"
	"case_study_api/container"
	"case_study_api/middleware"
	"case_study_api/routes"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()
	db := config.ConnectDatabase(cfg)
	config.ResetDatabase(db)

	// Initialize dependency injection container
	appContainer := container.NewContainer(db)

	r := gin.New()
	r.Use(
		gin.Logger(),
		gin.Recovery(),
		middleware.CustomLogger(),
	)

	routes.RegisterRoutes(r, appContainer)

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
