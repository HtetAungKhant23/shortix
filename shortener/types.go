package shortener

import (
	"time"

	"github.com/HtetAungKhant23/shortix/api/helper"
)

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

type CreateShortURLResponse struct {
	helper.ResponseBase
	Data *URL `json:"data,omitempty"`
}

type UpdateRequest struct {
	URL string `json:"url"`
}

type GetOriginalURLResponse struct {
	helper.ResponseBase
	Data *URL `json:"data,omitempty"`
}

type UpdateShortURLResponse struct {
	helper.ResponseBase
	Data *URL `json:"data,omitempty"`
}

type GetShortURLStatsResponse struct {
	helper.ResponseBase
	Data *URL `json:"data,omitempty"`
}
