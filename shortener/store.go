package shortener

import "errors"

var (
	ErrInvalidURL = errors.New("invalid URL")
)

type Store interface {
	Save(url *URL) error
}
