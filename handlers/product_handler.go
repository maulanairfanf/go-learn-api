package handlers

import (
	"errors"
	"myapi/db"
	"myapi/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetProducts handles retrieving all products along with their categories
func GetProducts(c *gin.Context) {
	var products []models.Product
	if err := db.DB.Preload("Categories").Find(&products).Error; err != nil {
		ErrorResponse(c, 500, err.Error())
		return
	}
	SuccessResponse(c, products)
}


// GetProduct handles retrieving a single product by ID along with its categories
func GetProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		ErrorResponse(c, 400, "Invalid product ID")
		return
	}
	var product models.Product
	if err := db.DB.Preload("Categories").First(&product, id).Error; err != nil {
		ErrorResponse(c, 404, "Product not found")
		return
	}
	SuccessResponse(c, product)
}


// CreateProduct handles the creation of a new product
func CreateProduct(c *gin.Context) {
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		ErrorResponse(c, 400, "Invalid request payload")
		return
	}
	if err := db.DB.Create(&product).Error; err != nil {
		ErrorResponse(c, 500, err.Error())
		return
	}
	SuccessResponse(c, product)
}


// UpdateProduct handles the update of an existing product
func UpdateProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		ErrorResponse(c, 400, "Invalid product ID")
		return
	}
	var product models.Product
	if err := db.DB.First(&product, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ErrorResponse(c, 404, "Product not found")
		} else {
			ErrorResponse(c, 500, err.Error())
		}
		return
	}
	if err := c.ShouldBindJSON(&product); err != nil {
		ErrorResponse(c, 400, "Invalid request payload")
		return
	}
	if err := db.DB.Save(&product).Error; err != nil {
		ErrorResponse(c, 500, err.Error())
		return
	}
	SuccessResponse(c, product)
}


// DeleteProduct handles the deletion of a product by ID
func DeleteProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		ErrorResponse(c, 400, "Invalid product ID")
		return
	}
	var product models.Product
	if err := db.DB.Preload("Categories").First(&product, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ErrorResponse(c, 404, "Product not found")
		} else {
			ErrorResponse(c, 500, err.Error())
		}
		return
	}
	if err := db.DB.Model(&product).Association("Categories").Clear(); err != nil {
		ErrorResponse(c, 500, err.Error())
		return
	}
	if err := db.DB.Delete(&product).Error; err != nil {
		ErrorResponse(c, 500, err.Error())
		return
	}
	SuccessResponse(c, "Product deleted successfully")
}
