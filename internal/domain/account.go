package domain

import "time"

type Account struct {
	ID           int
	Email        string
	PasswordHash string
	CreatedAt    time.Time
}
