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
	router.POST("/admin/login", handlers.Login)
	router.POST("/admin/register", handlers.Register)

	adminUser := router.Group("/admin/user")
	adminUser.Use(middleware.JWTMiddlewareGin())
	{
		adminUser.GET("", handlers.GetUsers)
		adminUser.GET(":id", handlers.GetUser)
		adminUser.PUT(":id", handlers.UpdateUser)
		adminUser.DELETE(":id", handlers.DeleteUser)
	}

	products := router.Group("/products")
	products.GET("", handlers.GetPublicProducts)
	products.GET(":slug", handlers.GetPublicProductBySlug)

	adminProducts := router.Group("/admin/products")
	adminProducts.Use(middleware.JWTMiddlewareGin())
	{
		adminProducts.GET("", handlers.GetProducts)
		adminProducts.GET(":id", handlers.GetProduct)
		adminProducts.POST("", handlers.CreateProduct)
		adminProducts.DELETE(":id", handlers.DeleteProduct)
		adminProducts.PUT(":id", handlers.UpdateProduct)
		adminProducts.PATCH(":id/publication_status", handlers.UpdatePublicationStatus)
	}

	categories := router.Group("/categories")
	categories.GET("", handlers.GetPublicCategories)
	categories.GET(":slug", handlers.GetPublicCategoriesBySlug)

	adminCategory := router.Group("/admin/categories")
	adminCategory.Use(middleware.JWTMiddlewareGin())
	{
		adminCategory.GET("", handlers.GetCategories)
		adminCategory.GET(":id", handlers.GetCategory)
		adminCategory.POST("", handlers.CreateCategory)
		adminCategory.PUT(":id", handlers.UpdateCategory)
		adminCategory.DELETE(":id", handlers.DeleteCategory)
		adminCategory.PATCH(":id", handlers.ChangeStatus)
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
