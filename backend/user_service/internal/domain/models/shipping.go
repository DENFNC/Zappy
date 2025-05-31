package models

type Shipping struct {
	AddressID  string
	ProfileID  string
	Country    string
	City       string
	Street     string
	PostalCode string
	IsDefault  bool
}
