package handlers

import (
	"errors"
	"myapi/services"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var publicCategoriesService = services.NewPublicCategoriesService()

func GetPublicCategories(c *gin.Context) {
	publicCategories, err := publicCategoriesService.GetAll()
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
		Data:    publicCategories,
		Message: "Succes",
	})
}

func GetPublicCategoriesBySlug(c *gin.Context) {
	slug := strings.TrimSpace(c.Param("slug"))

	if slug == "" {
		ErrorResponse(ErrorParams{
			C:       c,
			Status:  400,
			Message: "Invalid Category",
		})
		return
	}

	publicCategories, err := publicCategoriesService.GetBySlug(slug)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ErrorResponse(ErrorParams{
				C:       c,
				Status:  404,
				Message: "Product not found",
			})
		} else {
			ErrorResponse(ErrorParams{
				C:       c,
				Status:  500,
				Message: err.Error(),
			})
		}

		return
	}

	SuccessResponse(ResponseParams{
		C:       c,
		Data:    publicCategories,
		Message: "Succes",
	})
}
