package models

type Payment struct {
	PaymentID    uint32
	ProfileID    uint32
	PaymentToken string
	IsDefault    bool
}

func NewPayment(profileID uint32, paymentToken string, isDefault bool) *Payment {
	return &Payment{
		ProfileID:    profileID,
		PaymentToken: paymentToken,
		IsDefault:    isDefault,
	}
}
