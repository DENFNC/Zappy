package models

type Payment struct {
	PaymentID    string
	ProfileID    string
	PaymentToken string
	IsDefault    bool
}
