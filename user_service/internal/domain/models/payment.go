package models

type PaymentMethod struct {
	PaymentID    int
	ProfileID    int
	PaymentToken string
	IsDefault    bool
}
