package models

import "time"

type SubscriptionStatus string

const (
	SubActive   SubscriptionStatus = "ACTIVE"
	SubExpired  SubscriptionStatus = "EXPIRED"
	SubCanceled SubscriptionStatus = "CANCELED"
)

type Subscription struct {
	UserID    string             `json:"userId"`
	Status    SubscriptionStatus `json:"status"`
	ExpiresAt time.Time          `json:"expiresAt"`
	UpdatedAt time.Time          `json:"updatedAt"`
	CreatedAt time.Time          `json:"createdAt"`
}
