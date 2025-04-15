package models

import "time"

type WishlistItem struct {
	ItemID     int
	WishlistID int
	ProductID  int
	AddedAt    time.Time
}
