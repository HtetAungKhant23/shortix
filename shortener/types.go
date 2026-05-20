package shortener

import "time"

type URL struct {
	ID          string    `json:"id"`
	URL         string    `json:"url"`
	ShortCode   string    `json:"short_code"`
	AccessCount int64     `json:"access_count,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateRequest struct {
	URL string `json:"url"`
}
