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

func TestShorten_UniqueShortCodes(t *testing.T) {
	svc := newSvc()

	codes := make(map[string]struct{})
	for i := 1; i < 100; i++ {
		entity, err := svc.Shorten("https://roadmap.sh")
		if err != nil {
			t.Fatalf("iteration %d: %v", i, err)
		}

		if _, duplicate := codes[entity.ShortCode]; duplicate {
			t.Fatalf("duplicate short code %s at iteration %d", entity.ShortCode, i)
		}
		codes[entity.ShortCode] = struct{}{}
	}
}

func TestRetrieve_GetOringailURL(t *testing.T) {
	url := "https://roadmap.sh"

	svc := newSvc()
	created, _ := svc.Shorten(url)

	entity, err := svc.Retrieve(created.ShortCode)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entity.URL != url {
		t.Errorf("got %s, want %s", entity.URL, url)
	}
}

func TestRetrieve_ShortCodeNotFound(t *testing.T) {
	svc := newSvc()
	_, err := svc.Retrieve("does-not-exist")

	if !errors.Is(err, shortener.ErrNotFound) {
		t.Errorf("expected not found error message, got %v", err)
	}
}

func TestUpdate_Success(t *testing.T) {
	svc := newSvc()

	oldURL := "https://roadmap.sh"
	entity, err := svc.Shorten(oldURL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	newURL := "https://example.sh"
	updatedURL, err := svc.Update(entity.ShortCode, newURL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if updatedURL.ShortCode != entity.ShortCode {
		t.Errorf("expected same short code, but got %s", updatedURL.ShortCode)
	}

	if updatedURL.URL != newURL {
		t.Errorf("expected URL %s but got %s", newURL, updatedURL.URL)
	}
}

func TestUpdate_NotFound(t *testing.T) {
	svc := newSvc()

	_, err := svc.Update("hsgeqh", "https://roadmap.sh")
	if !errors.Is(err, shortener.ErrNotFound) {
		t.Errorf("expected not found error message, but got %v", err)
	}
}

func TestUpdate_InvalidURL(t *testing.T) {
	svc := newSvc()
	old, _ := svc.Shorten("https://example.com")

	_, err := svc.Update(old.ShortCode, "bad-url")
	if !errors.Is(err, shortener.ErrInvalidURL) {
		t.Errorf("expected invalid URL error message, but got %v", err)
	}
}
