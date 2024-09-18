package dto

import "time"

type Comment struct {
	ID        ID
	CreatedAt time.Time
	UpdatedAt time.Time
	Author    *User
	Body      string
}
