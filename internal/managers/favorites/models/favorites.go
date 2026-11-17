package models

import "api-tests-template/internal/managers/advertisement/models"

const FavoriteLimit = 50

type AddRequest struct {
	AdvertisementID string `json:"advertisement_id"`
}

type BulkAddRequest struct {
	AdvertisementIDs []string `json:"advertisement_ids"`
}

type BulkAddResult struct {
	Added        []string `json:"added"`
	NotFound     []string `json:"not_found"`
	Skipped      []string `json:"skipped"`
	LimitReached bool     `json:"limit_reached"`
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
