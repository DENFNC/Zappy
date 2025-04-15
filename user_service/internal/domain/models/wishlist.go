package models

import "time"

type WishlistItem struct {
	ItemID    int
	ProfileID int
	ProductID int
	AddedAt   time.Time
}
