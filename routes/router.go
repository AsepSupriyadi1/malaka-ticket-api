package routes

import (
	"case_study_api/container"
	"case_study_api/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, container *container.Container) {
	auth := r.Group("/auth")
	AuthRoutes(auth, container)

	api := r.Group("/api")
	api.Use(middleware.JWTAuth())

	EventRoutes(api, container)
	TicketRoutes(api, container)
	ReportRoutes(api, container)
}
