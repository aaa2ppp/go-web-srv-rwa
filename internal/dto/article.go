package dto

import "time"

type Article struct {
	// ID             uint64    `json:"id,omitempty"`
	CreatedAt      time.Time `json:"createdAt,omitempty"`
	UpdatedAt      time.Time `json:"updatedAt,omitempty"`
	Author         *Profile  `json:"author,omitempty"`
	Slug           string    `json:"slug,omitempty"`
	Title          string    `json:"title,omitempty"`
	TagList        []string  `json:"tagList,omitempty"`
	Description    string    `json:"description,omitempty"`
	Body           string    `json:"body,omitempty"`
	Favorited      bool      `json:"favorited,omitempty"`
	FavoritesCount uint      `json:"favoritesCount,omitempty"`
}

type Profile struct {
	BIO       string `json:"bio,omitempty"`
	Following bool   `json:"following,omitempty"`
	Image     string `json:"image,omitempty"`
	Username  string `json:"username,omitempty"`
}

type ListArticle struct {
	Tag       string
	Author    string
	Favorited string
	Limit     uint64
	Offset    uint64
}

type CreateArticle struct {
	Body        string   `json:"body,omitempty"`
	Description string   `json:"description,omitempty"`
	TagList     []string `json:"tagList,omitempty"`
	Title       string   `json:"title,omitempty"`
	Token       string   `json:"-"`
}
