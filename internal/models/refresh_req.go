package models

type RefreshRequest struct {
	RefreshTokenHash string `json:"refresh_token_hash"`
	User             User   `json:"user"`
}
