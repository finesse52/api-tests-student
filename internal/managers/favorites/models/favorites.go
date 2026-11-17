package models

import "api-tests-template/internal/managers/advertisement/models"

type AddRequest struct {
	AdvertisementID string `json:"advertisement_id"`
}

type Item struct {
	FavoriteID    string               `json:"favorite_id"`
	Advertisement models.Advertisement `json:"advertisement"`
}

type List struct {
	Items      []Item `json:"items"`
	Page       int    `json:"page"`
	PageSize   int    `json:"page_size"`
	Total      int    `json:"total"`
	TotalPages int    `json:"total_pages"`
}
