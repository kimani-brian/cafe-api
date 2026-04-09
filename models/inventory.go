package models

// InventoryItem represents a product in our café
type InventoryItem struct {
	ID            int     `json:"id"`
	Name          string  `json:"name" binding:"required"`
	Price         float64 `json:"price" binding:"required,gt=0"`  // gt=0 means "greater than zero"
	StockQuantity int     `json:"stock_quantity" binding:"gte=0"` // gte=0 means "greater than or equal to zero"
}
