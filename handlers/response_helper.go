package handlers

import (
	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
}

func SuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(200, APIResponse{
		Status: 200,
		Data:   data,
	})
}

func ErrorResponse(c *gin.Context, status int, message string) {
	c.JSON(status, APIResponse{
		Status: status,
		Data:   gin.H{"error": message},
	})
}