package controllers

import (
	"database/sql"
	"log"
	"net/http"

	"cafe-api/database"
	"cafe-api/models"

	"github.com/gin-gonic/gin"
)

func CreateOrder(c *gin.Context) {
	var req models.OrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order format"})
		return
	}

	// 1. Get the Cashier's ID from the JWT token (set by our Auth Middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in token"})
		return
	}

	// ==========================================
	// 2. BEGIN THE TRANSACTION
	// ==========================================
	tx, err := database.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}

	// defer a Rollback. If the transaction successfully Commits later,
	// this Rollback does nothing. If the function crashes/returns early, this ensures
	// the database doesn't stay locked forever.
	defer tx.Rollback()

	var totalAmount float64

	// We will temporarily store the processed items here to insert into the receipt later
	type processedItem struct {
		itemID   int
		quantity int
		price    float64
	}
	var processedItems []processedItem

	// 3. Process each item in the cart
	for _, item := range req.Items {
		var currentStock int
		var currentPrice float64

		// A. Check stock and price WITH A LOCK (FOR UPDATE)
		// This freezes this specific row so no other cashier can buy this exact item at this exact millisecond
		err := tx.QueryRow(`
			SELECT stock_quantity, price 
			FROM inventory_items 
			WHERE id = $1 FOR UPDATE
		`, item.ItemID).Scan(&currentStock, &currentPrice)

		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Item ID not found"})
			return // Triggering return triggers the deferred Rollback!
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		// B. Check if we have enough stock
		if currentStock < item.Quantity {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Insufficient stock",
				"item_id": item.ItemID,
			})
			return // ROLLBACK! Transaction aborted. No harm done.
		}

		// C. Deduct the inventory
		_, err = tx.Exec(`
			UPDATE inventory_items 
			SET stock_quantity = stock_quantity - $1 
			WHERE id = $2
		`, item.Quantity, item.ItemID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update inventory"})
			return // ROLLBACK!
		}

		// D. Calculate total cost and save for the receipt
		totalAmount += currentPrice * float64(item.Quantity)
		processedItems = append(processedItems, processedItem{
			itemID:   item.ItemID,
			quantity: item.Quantity,
			price:    currentPrice, // We freeze the price here for historical accuracy!
		})
	}

	// 4. Create the main Order record
	var orderID int
	err = tx.QueryRow(`
		INSERT INTO orders (cashier_id, total_amount) 
		VALUES ($1, $2) RETURNING id
	`, userID, totalAmount).Scan(&orderID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order record"})
		return // ROLLBACK!
	}

	// 5. Create the Order Items (The Receipt)
	for _, pItem := range processedItems {
		_, err = tx.Exec(`
			INSERT INTO order_items (order_id, item_id, quantity, price_at_purchase) 
			VALUES ($1, $2, $3, $4)
		`, orderID, pItem.itemID, pItem.quantity, pItem.price)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save order items"})
			return // ROLLBACK! Even if the order was created, it gets wiped completely.
		}
	}

	// ==========================================
	// 6. COMMIT THE TRANSACTION
	// ==========================================
	if err = tx.Commit(); err != nil {
		log.Println("Commit failed:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to finalize order"})
		return
	}

	// 7. Success!
	c.JSON(http.StatusCreated, gin.H{
		"message":      "Order completed successfully",
		"order_id":     orderID,
		"total_amount": totalAmount,
	})
}
