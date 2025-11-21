package handlers

import (
	"myapi/db"
	"myapi/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetCategories handles retrieving all categories
func GetCategories(c *gin.Context) {
	var categories []models.Category
	if err := db.DB.Find(&categories).Error; err != nil {
		ErrorResponse(c, 500, err.Error())
		return
	}
	SuccessResponse(c, categories)
}


// GetCategory handles retrieving a single category by ID
func GetCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		ErrorResponse(c, 400, "Invalid category ID")
		return
	}
	var category models.Category
	if err := db.DB.First(&category, id).Error; err != nil {
		ErrorResponse(c, 404, "Category not found")
		return
	}
	SuccessResponse(c, category)
}


// CreateCategory handles the creation of a new category
func CreateCategory(c *gin.Context) {
	var category models.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		ErrorResponse(c, 400, "Invalid request payload")
		return
	}
	if err := db.DB.Create(&category).Error; err != nil {
		ErrorResponse(c, 500, err.Error())
		return
	}
	SuccessResponse(c, category)
}
