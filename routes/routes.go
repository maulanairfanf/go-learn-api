package routes

import (
	"myapi/handlers"
	"myapi/middleware"

	"github.com/gin-gonic/gin"
)

// InitializeRoutes sets up the Gin router with all the routes and middleware
func InitializeRoutes() *gin.Engine {
	router := gin.Default()

	// Root endpoint for health check or welcome message
	router.GET("/", func(c *gin.Context) {
		c.String(200, "Welcome to Go Learn API!")
	})

	// Define routes
	router.POST("/login", handlers.Login)

	product := router.Group("/product")
	product.Use(middleware.JWTMiddlewareGin())
	{
		product.GET("", handlers.GetProducts)
		product.GET(":id", handlers.GetProduct)
		product.POST("", handlers.CreateProduct)
		router.DELETE(":id", middleware.JWTMiddlewareGin(), handlers.DeleteProduct)
		router.PUT(":id", middleware.JWTMiddlewareGin(), handlers.UpdateProduct)
	}


	category := router.Group("/category")
	category.Use(middleware.JWTMiddlewareGin())
	{
		category.GET("", handlers.GetCategories)
		category.GET(":id", handlers.GetCategory)
		category.POST("", handlers.CreateCategory)
		category.PUT(":id", handlers.UpdateCategory)
		category.DELETE(":id", handlers.DeleteCategory)
	}

	return router
}


