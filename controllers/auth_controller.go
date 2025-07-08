package controllers

import (
	"case_study_api/dto"
	"case_study_api/services"
	"case_study_api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService services.AuthService
}

func NewAuthController(authService services.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

func (ac *AuthController) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildErrorResponse("invalid request body"))
		return
	}

	res, err := ac.authService.Register(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, utils.BuildSuccessResponse("registration successful", res))
}

func (ac *AuthController) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.BuildErrorResponse("invalid request body"))
		return
	}

	res, err := ac.authService.Login(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, utils.BuildErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.BuildSuccessResponse("login successful", res))
}
