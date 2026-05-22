package shortener_test

import (
	"errors"
	"testing"

	"github.com/HtetAungKhant23/shortix/shortener"
	"github.com/HtetAungKhant23/shortix/store"
)

func newSvc() *shortener.Service {
	return shortener.New(store.NewInMemoryStore())
}

func TestShorten_ValidURL(t *testing.T) {
	svc := newSvc()

	originalURL := "https://roadmap.sh"
	entity, err := svc.Shorten(originalURL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if entity.ShortCode == "" {
		t.Errorf("expected non-empty short code")
	}

	if entity.URL != originalURL {
		t.Errorf("expected same URL: got '%s', expected '%s'", entity.URL, originalURL)
	}
}

func TestShorten_InvalidURL(t *testing.T) {
	svc := newSvc()

	urls := []string{"", "abcd-efgh", "abc://test.com"}
	for _, url := range urls {
		_, err := svc.Shorten(url)
		if !errors.Is(err, shortener.ErrInvalidURL) {
			t.Errorf("expected invalid url error message: got %v", err)
		}
	}
}
