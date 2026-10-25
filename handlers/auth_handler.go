package handlers

import (
	"errors"
	"myapi/models"
	"myapi/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var authService = services.NewAuthService()

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required,min=6"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func Login(c *gin.Context) {
	var loginReq LoginRequest
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		ErrorResponse(ErrorParams{C: c, Status: 400, Message: "Invalid request payload"})
		return
	}

	user, err := authService.Login(loginReq.Username, loginReq.Password)
	if err != nil {
		if errors.Is(err, services.ErrInvalidInput) || errors.Is(err, gorm.ErrRecordNotFound) {
			ErrorResponse(ErrorParams{C: c, Status: 401, Message: "Invalid username or password"})
		} else {
			ErrorResponse(ErrorParams{C: c, Status: 500, Message: err.Error()})
		}
		return
	}

	token, err := authService.GenerateToken(user.ID)
	if err != nil {
		ErrorResponse(ErrorParams{C: c, Status: 500, Message: "Failed to generate token"})
		return
	}
	SuccessResponse(ResponseParams{C: c, Data: LoginResponse{Token: token}, Message: "success"})
}

func Register(c *gin.Context) {
	var registerReq RegisterRequest
	if err := c.ShouldBindJSON(&registerReq); err != nil {
		ErrorResponse(ErrorParams{C: c, Status: 400, Message: "Invalid request payload"})
		return
	}

	user := models.User{
		Username: registerReq.Username,
		Password: registerReq.Password,
		Email:    registerReq.Email,
	}

	_, err := authService.Register(user)
	if err != nil {
		ErrorResponse(ErrorParams{C: c, Status: 500, Message: err.Error()})
		return
	}

	SuccessResponse(ResponseParams{C: c, Data: "Success Registration", Message: "success"})
}
