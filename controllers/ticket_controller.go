package controllers

import (
	"case_study_api/dto"
	"case_study_api/services"
	"case_study_api/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TicketController struct {
	ticketService services.TicketService
}

func NewTicketController(ticketService services.TicketService) *TicketController {
	return &TicketController{
		ticketService: ticketService,
	}
}

func (tc *TicketController) GetTicketsPaginated(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	pagination := utils.GetPaginationFromQuery(c)

	result, err := tc.ticketService.GetUserTicketsPaginated(userID, pagination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.BuildErrorResponse("failed to fetch tickets"))
		return
	}

	c.JSON(http.StatusOK, utils.BuildSuccessResponse("success", result))
}

func (tc *TicketController) GetTicket(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildErrorResponse("invalid id"))
		return
	}

	ticket, err := tc.ticketService.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, utils.BuildErrorResponse("ticket not found"))
		return
	}
	c.JSON(http.StatusOK, utils.BuildSuccessResponse("success", ticket))
}

func (tc *TicketController) BookTicket(c *gin.Context) {
	var req dto.CreateTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildErrorResponse("invalid request body"))
		return
	}

	userID := c.MustGet("user_id").(uint)
	ticket, err := tc.ticketService.BookTicket(req, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusCreated, utils.BuildSuccessResponse("ticket booked", ticket))
}

func (tc *TicketController) CancelTicket(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildErrorResponse("invalid id"))
		return
	}

	var req dto.CancelTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildErrorResponse("invalid request body"))
		return
	}

	userID := c.MustGet("user_id").(uint)
	if err := tc.ticketService.CancelTicket(uint(id), userID, req); err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.BuildSuccessResponse("ticket cancelled", nil))
}
