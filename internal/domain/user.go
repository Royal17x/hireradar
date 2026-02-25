package domain

import "time"

type User struct {
	UserID    int
	TgID      string
	Username  string
	IsActive  bool
	CreatedAt time.Time
}
