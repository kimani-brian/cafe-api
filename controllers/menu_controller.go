package controllers

import (
	"log"
	"net/http"

	"cafe-api/database"
	"cafe-api/models"

	"github.com/gin-gonic/gin"
)

// GetMenu is for the Cashier's tablet. It only returns items that are actually in stock.
func GetMenu(c *gin.Context) {
	// Query the database for items with stock > 0
	rows, err := database.DB.Query(`
		SELECT id, name, price, stock_quantity 
		FROM inventory_items 
		WHERE stock_quantity > 0
		ORDER BY name ASC
	`)
	if err != nil {
		log.Println("Error querying menu:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch menu"})
		return
	}
	defer rows.Close()

	// Parse the SQL rows into a Go slice (array)
	var menu []models.InventoryItem
	for rows.Next() {
		var item models.InventoryItem
		if err := rows.Scan(&item.ID, &item.Name, &item.Price, &item.StockQuantity); err != nil {
			log.Println("Error scanning row:", err)
			continue
		}
		menu = append(menu, item)
	}

	// Return the menu as JSON. If it's empty, return an empty array instead of null
	if menu == nil {
		menu = []models.InventoryItem{}
	}
	c.JSON(http.StatusOK, menu)
}

// CreateInventoryItem lets admins add a new item to inventory.
func CreateInventoryItem(c *gin.Context) {
	var item models.InventoryItem
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid inventory item format"})
		return
	}

	err := database.DB.QueryRow(
		`INSERT INTO inventory_items (name, price, stock_quantity)
		 VALUES ($1, $2, $3) RETURNING id`,
		item.Name,
		item.Price,
		item.StockQuantity,
	).Scan(&item.ID)
	if err != nil {
		log.Println("Error creating inventory item:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create inventory item"})
		return
	}

	c.JSON(http.StatusCreated, item)
}
