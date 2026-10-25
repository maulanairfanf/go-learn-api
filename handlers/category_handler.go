package handlers

import (
	"errors"
	"myapi/models"
	"myapi/services"
	"myapi/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var categoryService = services.NewCategoryService()

func GetCategories(c *gin.Context) {
	var pagination utils.PaginationRequest
	c.ShouldBindQuery(&pagination)
	// pagination.SetDefault()

	search := c.Query("search")
	sort := c.Query("sort")
	order := c.Query("order")

	categories, total, err := categoryService.GetAll(pagination, search, sort, order)
	if err != nil {
		ErrorResponse(ErrorParams{C: c, Status: 500, Message: err.Error()})
		return
	}
	var meta *utils.PaginationMeta
	if pagination.IsPaginated() {
		m := pagination.ComputeMeta(total)
		meta = &m
	}
	SuccessResponse(ResponseParams{
		C:       c,
		Data:    categories,
		Message: "success",
		Meta:    meta,
	})
}

func GetCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		ErrorResponse(ErrorParams{C: c, Status: 400, Message: "Invalid Category ID"})
		return
	}

	category, err := categoryService.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ErrorResponse(ErrorParams{C: c, Status: 404, Message: "Category not found"})
		} else {
			ErrorResponse(ErrorParams{C: c, Status: 500, Message: err.Error()})
		}
		return
	}
	SuccessResponse(ResponseParams{C: c, Data: category, Message: "success"})
}

func CreateCategory(c *gin.Context) {
	var req models.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(ErrorParams{C: c, Status: 400, Message: "Invalid request payload"})
		return
	}

	category, err := categoryService.Create(req)
	if err != nil {
		ErrorResponse(ErrorParams{C: c, Status: 500, Message: err.Error()})
		return
	}
	SuccessResponse(ResponseParams{C: c, Data: category, Message: "success"})
}

func UpdateCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		ErrorResponse(ErrorParams{C: c, Status: 400, Message: "Invalid Category ID"})
		return
	}

	var req models.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(ErrorParams{C: c, Status: 400, Message: "Invalid request payload"})
		return
	}

	category, err := categoryService.Update(id, req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ErrorResponse(ErrorParams{C: c, Status: 404, Message: "Category not found"})
		} else {
			ErrorResponse(ErrorParams{C: c, Status: 500, Message: err.Error()})
		}
		return
	}
	SuccessResponse(ResponseParams{C: c, Data: category, Message: "success"})
}

func DeleteCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		ErrorResponse(ErrorParams{C: c, Status: 400, Message: "Invalid Category Id"})
		return
	}

	_, err = categoryService.Delete(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ErrorResponse(ErrorParams{C: c, Status: 404, Message: "Category not found"})
		} else {
			ErrorResponse(ErrorParams{C: c, Status: 500, Message: err.Error()})
		}
		return
	}
	SuccessResponse(ResponseParams{C: c, Data: "Category deleted successfully", Message: "success"})
}

func ChangeStatus(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		ErrorResponse(ErrorParams{C: c, Status: 400, Message: "Invalid Category Id"})
		return
	}

	_, err = categoryService.ChangeStatus(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ErrorResponse(ErrorParams{C: c, Status: 404, Message: "Category not found"})
		} else {
			ErrorResponse(ErrorParams{C: c, Status: 500, Message: err.Error()})
		}
		return
	}
	SuccessResponse(ResponseParams{C: c, Data: "Category status changed successfully", Message: "success"})
}
