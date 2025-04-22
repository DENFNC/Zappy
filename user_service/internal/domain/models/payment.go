package models

type Payment struct {
	PaymentID    uint32 `db:"payment_id"`
	ProfileID    uint32 `db:"profile_id"`
	PaymentToken string `db:"payment_token"`
	IsDefault    bool   `db:"is_default"`
}

func NewPayment(profileID uint32, paymentToken string) *Payment {
	return &Payment{
		ProfileID:    profileID,
		PaymentToken: paymentToken,
	}
}
