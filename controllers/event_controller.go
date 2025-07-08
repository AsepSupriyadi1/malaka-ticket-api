package controllers

import (
	"case_study_api/dto"
	"case_study_api/services"
	"case_study_api/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type EventController struct {
	eventService services.EventService
}

func NewEventController(eventService services.EventService) *EventController {
	return &EventController{
		eventService: eventService,
	}
}

func (ec *EventController) GetEventsPaginated(c *gin.Context) {
	pagination := utils.GetPaginationFromQuery(c)

	result, err := ec.eventService.GetAllPaginated(pagination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.BuildErrorResponse("failed to fetch events"))
		return
	}

	c.JSON(http.StatusOK, utils.BuildSuccessResponse("success", result))
}

func (ec *EventController) GetEventByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildErrorResponse("invalid event id"))
		return
	}

	event, err := ec.eventService.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, utils.BuildErrorResponse("event not found"))
		return
	}
	c.JSON(http.StatusOK, utils.BuildSuccessResponse("success", event))
}

func (ec *EventController) CreateEvent(c *gin.Context) {
	var req dto.CreateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildErrorResponse("invalid request body"))
		return
	}

	creator := c.MustGet("user_id").(uint)
	created, err := ec.eventService.Create(req, creator)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusCreated, utils.BuildSuccessResponse("event created", created))
}

func (ec *EventController) UpdateEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildErrorResponse("invalid id"))
		return
	}

	var req dto.UpdateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildErrorResponse("invalid body"))
		return
	}

	updated, err := ec.eventService.Update(uint(id), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.BuildSuccessResponse("event updated", updated))
}

func (ec *EventController) DeleteEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildErrorResponse("invalid id"))
		return
	}

	if err := ec.eventService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.BuildSuccessResponse("event deleted", nil))
}
