package shortener

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/url"
	"strconv"
	"sync/atomic"
	"time"
)

var counter uint64

type Service struct {
	store   Store
	codeLen int
}

func New(store Store) *Service {
	return &Service{
		store:   store,
		codeLen: 6,
	}
}

func (s *Service) Shorten(url string) (*URL, error) {
	if err := validateURL(url); err != nil {
		return nil, err
	}

	code, err := s.generateCode()
	if err != nil {
		return nil, fmt.Errorf("failed to generate short code: %s", err.Error())
	}

	id := atomic.AddUint64(&counter, 1)
	time := time.Now().UTC()

	entity := &URL{
		ID:        strconv.FormatUint(id, 10),
		URL:       url,
		ShortCode: code,
		CreatedAt: time,
		UpdatedAt: time,
	}

	return entity, nil
}

func (s *Service) generateCode() (string, error) {
	b := make([]byte, s.codeLen)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	code := base64.URLEncoding.EncodeToString(b)[:s.codeLen]
	return code, nil
}

func validateURL(rawURL string) error {
	if rawURL == "" {
		return ErrInvalidURL
	}

	u, err := url.ParseRequestURI(rawURL)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") {
		return ErrInvalidURL
	}

	return nil
}
