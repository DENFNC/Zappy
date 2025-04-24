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

func NewShipping(
	profileID string,
	country string,
	city string,
	street string,
	postalCode string,
) *Shipping {
	return &Shipping{
		ProfileID:  profileID,
		Country:    country,
		City:       city,
		Street:     street,
		PostalCode: postalCode,
		IsDefault:  true,
	}
}
