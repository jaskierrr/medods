package models

import "time"

type RefreshToken struct {
	UserID         string
	UserIP         string
	TokenHash      string
	ExpirationTime time.Time
}
