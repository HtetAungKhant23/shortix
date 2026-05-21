package shortener

import "errors"

var (
	ErrNotFound   = errors.New("short URL not found")
	ErrCodeExists = errors.New("short code already exists")
	ErrInvalidURL = errors.New("invalid URL")
)

type Store interface {
	Save(url *URL) error
	FindByCode(code string) (*URL, error)
	IncrementAccess(code string) error
	Update(code string, newURL string) (*URL, error)
}
