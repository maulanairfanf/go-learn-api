package handlers

import (
	"myapi/utils"

	"github.com/gin-gonic/gin"
)

type ResponseParams struct {
	C       *gin.Context
	Data    interface{}
	Message string
	Meta    *utils.PaginationMeta
}

type ErrorParams struct {
	C       *gin.Context
	Status  int
	Data    interface{}
	Message string
}

type APIResponse struct {
	Status  int                   `json:"status"`
	Data    interface{}           `json:"data"`
	Message string                `json:"message"`
	Meta    *utils.PaginationMeta `json:"meta,omitempty"`
}

func SuccessResponse(p ResponseParams) {
	p.C.JSON(200, APIResponse{
		Status:  200,
		Data:    p.Data,
		Message: p.Message,
		Meta:    p.Meta,
	})
}

func ErrorResponse(p ErrorParams) {
	p.C.JSON(p.Status, APIResponse{
		Status:  p.Status,
		Data:    p.Data,
		Message: p.Message,
	})
}
