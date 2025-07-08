package routes

import (
	"case_study_api/container"
	"case_study_api/controllers"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(rg *gin.RouterGroup, container *container.Container) {
	authController := controllers.NewAuthController(container.AuthService)

	rg.POST("/register", authController.Register)
	rg.POST("/login", authController.Login)
}
