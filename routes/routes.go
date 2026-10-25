package routes

import (
	"myapi/handlers"
	"myapi/middleware"

	"github.com/gin-gonic/gin"
)

// InitializeRoutes sets up the Gin router with all the routes and middleware
func InitializeRoutes() *gin.Engine {
	router := gin.Default()
	router.Use(middleware.LoggerMiddleware())

	// Root endpoint for health check or welcome message
	router.GET("/", func(c *gin.Context) {
		c.String(200, "Welcome to Go Learn API!")
	})

	// Define routes
	router.POST("/login", handlers.Login)
	router.POST("/register", handlers.Register)

	user := router.Group("/user")
	user.Use(middleware.JWTMiddlewareGin())
	{
		user.GET("", handlers.GetUsers)
		user.GET(":id", handlers.GetUser)
		user.PUT(":id", handlers.UpdateUser)
		user.DELETE(":id", handlers.DeleteUser)
	}

	product := router.Group("/product")
	product.Use(middleware.JWTMiddlewareGin())
	{
		product.GET("", handlers.GetProducts)
		product.GET(":id", handlers.GetProduct)
		product.POST("", handlers.CreateProduct)
		product.DELETE(":id", handlers.DeleteProduct)
		product.PUT(":id", handlers.UpdateProduct)
	}

	category := router.Group("/category")
	category.Use(middleware.JWTMiddlewareGin())
	{
		category.GET("", handlers.GetCategories)
		category.GET(":id", handlers.GetCategory)
		category.POST("", handlers.CreateCategory)
		category.PUT(":id", handlers.UpdateCategory)
		category.DELETE(":id", handlers.DeleteCategory)
		category.PATCH(":id", handlers.ChangeStatus)
	}

	order := router.Group("/order")
	order.Use(middleware.JWTMiddlewareGin())
	{
		order.GET("", handlers.GetOrders)
		order.GET(":id", handlers.GetOrder)
		order.DELETE(":id", handlers.DeleteOrder)
		order.POST("", handlers.CreateOrder)
		order.PUT(":id/paid", handlers.PayOrder)
		order.PUT(":id/cancel", handlers.CancelOrder)
	}

	return router
}
