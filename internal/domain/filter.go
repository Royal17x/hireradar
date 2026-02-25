package domain

import "time"

type Filter struct {
	ID        int
	UserID    int
	Keywords  string
	City      string
	Grade     string
	CreatedAt time.Time
}
