package models

type ShippingAddress struct {
	AddressID  int
	ProfileID  int
	Country    string
	City       string
	Street     string
	PostalCode string
	IsDefault  bool
}
