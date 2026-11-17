package models

type AdvertisementsList struct {
	Items      []Advertisement `json:"items"`
	HasNext    bool            `json:"has_next"`
	NextCursor string          `json:"next_cursor"`
}

type Advertisement struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	UserID string `json:"user_id"`
}
