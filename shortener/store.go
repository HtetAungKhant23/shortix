package shortener

import "errors"

var (
	ErrCodeExists = errors.New("short code already exists")
	ErrInvalidURL = errors.New("invalid URL")
)

type Store interface {
	Save(url *URL) error
}
