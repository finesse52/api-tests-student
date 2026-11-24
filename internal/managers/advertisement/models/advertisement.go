package models

type CreateRequest struct {
	Title        string
	Description  string
	Price        float64
	Quantity     int
	PriceText    string
	QuantityText string
	PhotoFiles   []string
}

type Photo struct {
	ID        string `json:"id"`
	URL       string `json:"url"`
	SortOrder int    `json:"sort_order"`
}

type Advertisement struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
	UserID      string  `json:"user_id"`
	Photos      []Photo `json:"photos"`
}
