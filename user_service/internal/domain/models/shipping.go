package models

type Shipping struct {
	AddressID  uint32
	ProfileID  uint32
	Country    string
	City       string
	Street     string
	PostalCode string
	IsDefault  bool
}

func NewShipping(
	profileID uint32,
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
