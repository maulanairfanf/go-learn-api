package handlers

import (
	"errors"
	"mime/multipart"
	"myapi/models"
	"myapi/services"
	"myapi/storage"
	"myapi/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var productService = services.NewProductService(storage.NewVercelBlobStorage())

func GetProducts(c *gin.Context) {
	var pagination utils.PaginationRequest
	c.ShouldBindQuery(pagination)

	search := c.Query("search")
	sort := c.Query("sort")
	order := c.Query("order")

	products, total, err := productService.GetAll(pagination, search, sort, order)
	if err != nil {
		ErrorResponse(ErrorParams{C: c, Status: 500, Message: err.Error()})
		return
	}

	var meta *utils.PaginationMeta
	if pagination.IsPaginated() {
		m := pagination.ComputeMeta(total)
		meta = &m
	}
	SuccessResponse(ResponseParams{C: c, Data: products, Message: "success", Meta: meta})
}

func GetProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		ErrorResponse(ErrorParams{C: c, Status: 400, Message: "Invalid product ID"})
		return
	}

	product, err := productService.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ErrorResponse(ErrorParams{C: c, Status: 404, Message: "Product not found"})
		} else {
			ErrorResponse(ErrorParams{C: c, Status: 500, Message: err.Error()})
		}
		return
	}
	SuccessResponse(ResponseParams{C: c, Data: product, Message: "success"})
}

func CreateProduct(c *gin.Context) {
	var req models.CreateProductRequest
	if err := c.ShouldBind(&req); err != nil {
		ErrorResponse(ErrorParams{C: c, Status: 400, Message: "Invalid product input"})
		return
	}

	file, err := c.FormFile("image")
	if !ValidationFileSize(c, file, err) {
		return
	}

	req.Image = file

	product, err := productService.Create(c.Request.Context(), req)
	if err != nil {
		var validationErr *services.CategoryValidationError
		if errors.As(err, &validationErr) {
			ErrorResponse(ErrorParams{
				C:       c,
				Status:  400,
				Data:    validationErr,
				Message: validationErr.Error(),
			})
			return
		} else if errors.Is(err, services.ErrConflict) {
			ErrorResponse(ErrorParams{
				C:       c,
				Status:  409,
				Message: err.Error(),
			})
			return
		} else {
			ErrorResponse(ErrorParams{
				C:       c,
				Status:  500,
				Message: err.Error(),
			})
		}
		return
	}
	SuccessResponse(ResponseParams{C: c, Data: product, Message: "success"})
}

func UpdateProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		ErrorResponse(ErrorParams{C: c, Status: 400, Message: "Invalid product ID"})
		return
	}

	var req models.UpdateProductRequest
	if err := c.ShouldBind(&req); err != nil {
		ErrorResponse(ErrorParams{C: c, Status: 400, Message: "Invalid request payload"})
		return
	}

	product, err := productService.Update(c.Request.Context(), id, req)
	if err != nil {
		var validationErr *services.CategoryValidationError
		if errors.As(err, &validationErr) {
			ErrorResponse(ErrorParams{
				C:       c,
				Status:  400,
				Data:    validationErr,
				Message: validationErr.Error(),
			})
		} else if errors.Is(err, services.ErrConflict) {
			ErrorResponse(ErrorParams{
				C:       c,
				Status:  409,
				Message: err.Error(),
			})
		} else if errors.Is(err, gorm.ErrRecordNotFound) {
			ErrorResponse(ErrorParams{C: c, Status: 404, Message: "Product not found"})
		} else if errors.Is(err, services.ErrInvalidInput) {
			ErrorResponse(ErrorParams{C: c, Status: 400, Message: err.Error()})
		} else {
			ErrorResponse(ErrorParams{C: c, Status: 500, Message: err.Error()})
		}
		return
	}
	SuccessResponse(ResponseParams{C: c, Data: product, Message: "success"})
}

func DeleteProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		ErrorResponse(ErrorParams{C: c, Status: 400, Message: "Invalid product ID"})
		return
	}

	if err := productService.Delete(c.Request.Context(), id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ErrorResponse(ErrorParams{C: c, Status: 404, Message: "Product not found"})
		} else {
			ErrorResponse(ErrorParams{C: c, Status: 500, Message: err.Error()})
		}
		return
	}
	SuccessResponse(ResponseParams{C: c, Data: "Product deleted successfully", Message: "success"})
}

func ValidationFileSize(c *gin.Context, file *multipart.FileHeader, err error) bool {
	if err != nil {
		ErrorResponse(ErrorParams{
			C:       c,
			Status:  400,
			Message: "Image is required",
		})

		return false
	}

	const maxImageSize = 2 * 1024 * 1024

	if file.Size > maxImageSize {
		ErrorResponse(ErrorParams{
			C:       c,
			Status:  400,
			Message: "Image must be smaller than 2 MB",
		})
		return false
	}

	return true
}
