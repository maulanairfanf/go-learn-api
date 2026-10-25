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

var userService = services.NewUserService()

func GetUsers(c *gin.Context) {
	var pagination utils.PaginationRequest
	c.ShouldBindQuery(&pagination)
	// pagination.SetDefault()

	search := c.Query("search")
	sort := c.Query("sort")
	order := c.Query("order")

	users, total, err := userService.GetAll(pagination, search, sort, order)
	if err != nil {
		ErrorResponse(ErrorParams{C: c, Status: 500, Message: err.Error()})
		return
	}

	var meta *utils.PaginationMeta
	if pagination.IsPaginated() {
		m := pagination.ComputeMeta(total)
		meta = &m
	}

	SuccessResponse(ResponseParams{C: c, Data: users, Message: "success", Meta: meta})
}

func GetUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		ErrorResponse(ErrorParams{C: c, Status: 400, Message: "Invalid User ID"})
		return
	}

	user, err := userService.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ErrorResponse(ErrorParams{C: c, Status: 404, Message: "User not found"})
		} else {
			ErrorResponse(ErrorParams{C: c, Status: 500, Message: err.Error()})
		}
		return
	}
	SuccessResponse(ResponseParams{C: c, Data: user, Message: "success"})
}

func DeleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		ErrorResponse(ErrorParams{C: c, Status: 400, Message: "Invalid User ID"})
		return
	}

	user, err := userService.Delete(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ErrorResponse(ErrorParams{C: c, Status: 404, Message: "User not found"})
		} else {
			ErrorResponse(ErrorParams{C: c, Status: 500, Message: err.Error()})
		}
		return
	}
	SuccessResponse(ResponseParams{C: c, Data: user, Message: "success"})
}

func UpdateUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		ErrorResponse(ErrorParams{C: c, Status: 400, Message: "Invalid User ID"})
		return
	}

	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(ErrorParams{C: c, Status: 400, Message: "Invalid Payload"})
		return
	}

	user, err := userService.Update(id, req)
	if err != nil {
		if errors.Is(err, services.ErrConflict) {
			ErrorResponse(ErrorParams{C: c, Status: 409, Message: "Username already taken"})
		} else if errors.Is(err, gorm.ErrRecordNotFound) {
			ErrorResponse(ErrorParams{C: c, Status: 404, Message: "User not found"})
		} else {
			ErrorResponse(ErrorParams{C: c, Status: 500, Message: err.Error()})
		}
		return
	}
	SuccessResponse(ResponseParams{C: c, Data: user, Message: "success"})
}
