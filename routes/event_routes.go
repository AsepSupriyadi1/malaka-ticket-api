package routes

import (
	"case_study_api/container"
	"case_study_api/controllers"
	"case_study_api/middleware"

	"github.com/gin-gonic/gin"
)

func EventRoutes(rg *gin.RouterGroup, container *container.Container) {
	eventController := controllers.NewEventController(container.EventService)

	event := rg.Group("/events")
	event.GET("", eventController.GetEventsPaginated)
	event.GET("/:id", eventController.GetEventByID)
	event.POST("", middleware.RoleAuth("admin"), eventController.CreateEvent)
	event.PUT("/:id", middleware.RoleAuth("admin"), eventController.UpdateEvent)
	event.DELETE("/:id", middleware.RoleAuth("admin"), eventController.DeleteEvent)
}
