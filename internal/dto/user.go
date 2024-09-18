package dto

import "time"

type User struct {
	ID        uint64    `json:"id,omitempty"`
	BIO       string    `json:"bio,omitempty"`
	Email     string    `json:"email,omitempty"`
	Image     string    `json:"image,omitempty"`
	Token     string    `json:"token,omitempty"`
	Username  string    `json:"username,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
}

type UpdateUser struct {
	ID       uint64  `json:"id,omitempty"`
	BIO      *string `json:"bio,omitempty"`
	Email    *string `json:"email,omitempty"`
	Image    *string `json:"image,omitempty"`
	Token    *string `json:"token,omitempty"`
	Username *string `json:"username,omitempty"`
}

type CreateUser struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
	Username string `json:"username,omitempty"`
}

type LoginUser struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}
