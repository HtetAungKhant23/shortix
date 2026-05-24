# Shortix

A URL shortening service written in Go — usable **as an importable package** or **as a standalone REST API server**.

---

## Table of Contents

- [Project Structure](#project-structure)
- [Usage as a Go Package](#usage-as-a-go-package)
- [Usage as an API Server](#usage-as-an-api-server)
- [API Reference](#api-reference)
- [Implementing a Custom Store](#implementing-a-custom-store)
- [Design Decisions](#design-decisions)

---

## Project Structure

```
shortix/
├── shortener/               #    Core package — import this in your project
│   ├── shortener.go         #    Service: Shorten, Resolve, Update, Delete, Stats
│   ├── store.go             #    Store interface + sentinel errors
│   ├── types.go             #    URL, CreateRequest, UpdateRequest types
│   └── shortener_test.go    #    Full unit test suite
│
├── store/
│   └── memory.go            #    Built-in in-memory Store implementation
│
├── api/
│   ├── helper               #    HTTP request/response helper
│   ├── handler.go           #    HTTP handlers (use shortener.Service internally)
│   └── router.go            #    Route registration (stdlib ServeMux, Go 1.22+)
│
├── cmd/
│   └── main.go              #    Entry point — wires everything and starts server
│
├── go.mod
└── README.md
```

---

## Usage as a Go Package

```go
import (
    "fmt"
    "github.com/HtetAungKhant23/shortix/shortener"
    "github.com/HtetAungKhant23/shortix/store"
)

func main() {
    // 1. Pick a store (or implement your own)
    memStore := store.NewMemoryStore()

    // 2. Create the service
    svc := shortener.New(memStore)

    // 3. Shorten a URL
    entity, err := svc.Shorten("https://www.example.com/some/very/long/path")
    if err != nil {
        // err is shortener.ErrInvalidURL for bad input
    }
    fmt.Println(entity.ShortCode) // e.g. "aB3xY7"

    // 4. Retrieve — also increments the access counter
    entity, err = svc.Retrieve("aB3xY7")
    fmt.Println(entity.URL) // "https://www.example.com/some/very/long/path"

    // 5. Update
    entity, err = svc.Update("aB3xY7", "https://www.example.com/updated/path")

    // 6. Stats
    stats, err := svc.Stats("aB3xY7")
    fmt.Println(stats.AccessCount) // 1

    // 7. Delete
    err = svc.Delete("aB3xY7")
}
```

### Handling Errors

```go
import "errors"

entry, err := svc.Retrieve("unknown")
switch {
case errors.Is(err, shortener.ErrNotFound):
    fmt.Println("code does not exist")
case errors.Is(err, shortener.ErrInvalidURL):
    fmt.Println("URL failed validation")
case err != nil:
    fmt.Println("unexpected error:", err)
}
```

---

## Usage as an API Server

```bash
### Run directly
make run

### Build and run
make build
./bin/shortix
```

The server starts on `:8080` by default (or `$PORT` if set).

---

## API Reference

### Create Short URL

```
POST /api/v1/shorten
Content-Type: application/json

{"url": "https://www.example.com/some/long/url"}
```

**201 Created**

```json
{
  "status": "success",
  "data": {
    "id": "1",
    "url": "https://www.example.com/some/long/url",
    "short_code": "_WAW79",
    "created_at": "2026-05-22T03:48:44.862898Z",
    "updated_at": "2026-05-22T03:48:44.862898Z"
  }
}
```

**400 Bad Request** — missing or invalid URL

```json
{
  "status": "error",
  "error": {
    "code": "generic",
    "message": "invalid or missing URL"
  }
}
```

---

### Retrieve Original URL

```
GET /api/v1/shorten/{code}
```

**200 OK**

```json
{
  "status": "success",
  "data": {
    "id": "1",
    "url": "https://www.example.com/some/long/url",
    "short_code": "_WAW79",
    "created_at": "2026-05-22T03:48:44.862898Z",
    "updated_at": "2026-05-22T03:48:44.862898Z"
  }
}
```

**404 Not Found**

```json
{
  "status": "error",
  "error": {
    "code": "generic",
    "message": "short URL not found"
  }
}
```

---

### Get URL Statistics

```
GET /api/v1/shorten/{code}/stats
```

**200 OK**

```json
{
  "status": "success",
  "data": {
    "id": "1",
    "url": "https://www.example.com/some/long/url",
    "short_code": "_WAW79",
    "access_count": 1,
    "created_at": "2026-05-22T03:48:44.862898Z",
    "updated_at": "2026-05-22T03:48:44.862898Z"
  }
}
```

**404 Not Found**

```json
{
  "status": "error",
  "error": {
    "code": "generic",
    "message": "short URL not found"
  }
}
```

---

### Update Short URL

```
PUT /api/v1/shorten/{code}
Content-Type: application/json

{"url": "https://www.example.com/some/updated/url"}
```

**200 OK**

```json
{
  "status": "success",
  "data": {
    "id": "1",
    "url": "https://www.example.com/some/updated/url",
    "short_code": "_WAW79",
    "created_at": "2026-05-22T03:48:44.862898Z",
    "updated_at": "2026-05-22T03:48:44.862898Z"
  }
}
```

**400 Bad Request** | **404 Not Found**

---

### Delete Short URL

```
DELETE /api/v1/shorten/{code}
```

**204 No Content** — successfully deleted

**404 Not Found**

```json
{
  "status": "error",
  "error": {
    "code": "generic",
    "message": "short URL not found"
  }
}
```

---

## Implementing a Custom Store

To use a real database, implement the `shortener.Store` interface:

```go
// shortener/store.go
type Store interface {
    Save(url *URL) error
    FindByCode(code string) (*URL, error)
    IncrementAccess(code string) error
    Update(code string, newURL string) (*URL, error)
    Delete(code string) error
}
```

Example skeleton for a Postgres store:

```go
package pgstore

import (
    "database/sql"
    "errors"

    "github.com/HtetAungKhant23/shortix/shortener"
)

type PostgresStore struct {
    db *sql.DB
}

func New(db *sql.DB) *PostgresStore {
    return &PostgresStore{db: db}
}

func (p *PostgresStore) Save(url *shortener.URL) error {
    _, err := p.db.Exec(
        `INSERT INTO urls (id, url, short_code, created_at, updated_at)
         VALUES ($1, $2, $3, $4, $5)`,
        url.ID, url.URL, url.ShortCode, url.CreatedAt, url.UpdatedAt,
    )
    return err
}

func (p *PostgresStore) FindByCode(code string) (*shortener.URL, error) {
    row := p.db.QueryRow(`SELECT id, url, short_code, created_at, updated_at, access_count
                          FROM urls WHERE short_code = $1`, code)
    var u shortener.URL
    if err := row.Scan(&u.ID, &u.URL, &u.ShortCode, &u.CreatedAt, &u.UpdatedAt, &u.AccessCount); err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, shortener.ErrNotFound
        }
        return nil, err
    }
    return &u, nil
}

// ... implement Update, Delete, IncrementAccess similarly
```

Then wire it in `main.go`:

```go
db, _ := sql.Open("postgres", os.Getenv("DATABASE_URL"))
pgStore := pgstore.New(db)
svc := shortener.New(pgStore)
```

Zero changes to `shortener` or `api` packages required.

---

## Running Tests

```bash
make test
```

### What's tested

- Valid and invalid URL creation
- Code uniqueness across 100 iterations
- Retrieve: found and not-found paths
- Access counter incremented on retrieve
- Update: success, not-found, invalid URL
- Delete: success and not-found
- Stats retrieval
- Concurrent access safety (50 goroutines)

---

## Design Decisions

### Why is the core package HTTP-free?

The `shortener` package imports only the standard library (`crypto/rand`, `net/url`, `time`, `sync`). This means you can embed `shortix` into a CLI tool, a gRPC service, a Lambda function, or a desktop app — no HTTP stack required.

### Why the `Store` interface?

Persistence is a plugin. Today it's in-memory; tomorrow it can be Redis, DynamoDB, or SQLite. The interface boundary means none of the business logic or HTTP code ever needs to change when you swap the store.

### Why sentinel errors instead of custom error types?

Sentinel errors (`var ErrNotFound = errors.New(...)`) are the idiomatic Go way to communicate "expected" failure modes. Callers use `errors.Is()` for type-safe matching. Custom error types would add boilerplate with no benefit here.

---

## Project Origin

Built as a solution to the [Roadmap.sh - URL Shortening Service Challenge](https://roadmap.sh/projects/url-shortening-service).

## **Build with ❤️ using Go**

---
