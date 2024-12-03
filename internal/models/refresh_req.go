package models

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
	User             User   `json:"user"`
}
