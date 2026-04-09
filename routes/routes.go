package routes

import (
	"cafe-api/controllers"
	"cafe-api/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	auth := router.Group("/auth")
	{
		auth.POST("/login", controllers.Login)
	}

	api := router.Group("/api")
	api.Use(middleware.RequireAuth())
	{
		// === MENU & INVENTORY ROUTES (Phase 4) ===

		// Both Cashiers and Admins can see the menu
		api.GET("/menu", middleware.RequireRole("admin", "cashier"), controllers.GetMenu)

		// ONLY Admins can add new inventory
		api.POST("/inventory/items", middleware.RequireRole("admin"), controllers.CreateInventoryItem)
		// Note: CreateInventoryItem is the function we discussed in Phase 4
		// (e.g., INSERT INTO inventory_items (name, price, stock) VALUES ...)

		// === ORDER ROUTES (Phase 5) ===

		// Both Cashiers and Admins can process an order
		api.POST("/orders", middleware.RequireRole("admin", "cashier"), controllers.CreateOrder)
	}

	return router
}
