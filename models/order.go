package models

// OrderItemRequest represents a single item a customer wants to buy
type OrderItemRequest struct {
	ItemID   int `json:"item_id" binding:"required"`
	Quantity int `json:"quantity" binding:"required,gt=0"`
}

// OrderRequest represents the full checkout cart
type OrderRequest struct {
	Items []OrderItemRequest `json:"items" binding:"required,min=1"` // min=1 means cart can't be empty
}
