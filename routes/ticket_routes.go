package routes

import (
	"case_study_api/container"
	"case_study_api/controllers"

	"github.com/gin-gonic/gin"
)

func TicketRoutes(rg *gin.RouterGroup, container *container.Container) {
	ticketController := controllers.NewTicketController(container.TicketService)

	ticket := rg.Group("/tickets")
	ticket.GET("", ticketController.GetTicketsPaginated)
	ticket.GET("/:id", ticketController.GetTicket)
	ticket.POST("", ticketController.BookTicket)
	ticket.PATCH("/:id", ticketController.CancelTicket)
}
