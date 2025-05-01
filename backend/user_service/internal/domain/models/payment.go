package models

type Payment struct {
	PaymentID    string
	ProfileID    string
	PaymentToken string
	IsDefault    bool
}

func NewPayment(profileID string, paymentToken string) *Payment {
	return &Payment{
		ProfileID:    profileID,
		PaymentToken: paymentToken,
	}
}
