package handlers

import (
	"myapi/services"
	"strings"

	"github.com/gin-gonic/gin"
)

var publicProductService = services.NewPublicProductService()

func GetPublicProducts(c *gin.Context) {
	publicProducts, err := publicProductService.GetAll()
	if err != nil {
		ErrorResponse(ErrorParams{
			C:       c,
			Status:  500,
			Message: err.Error(),
		})
		return
	}

	SuccessResponse(ResponseParams{
		C:       c,
		Data:    publicProducts,
		Message: "Succes",
	})
}

func GetPublicProductBySlug(c *gin.Context) {
	slug := strings.TrimSpace(c.Param("slug"))

	if slug == "" {
		ErrorResponse(ErrorParams{
			C:       c,
			Status:  400,
			Message: "Invalid Product",
		})
		return
	}

	publicProducts, err := publicProductService.GetBySlug(slug)

	if err != nil {
		ErrorResponse(ErrorParams{
			C:       c,
			Status:  500,
			Message: err.Error(),
		})
		return
	}

	SuccessResponse(ResponseParams{
		C:       c,
		Data:    publicProducts,
		Message: "Succes",
	})
}
