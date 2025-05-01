package models

import "time"

type WishlistItem struct {
	ItemID    string
	ProfileID string
	ProductID string
	AddedAt   time.Time
	IsActive  bool
}
