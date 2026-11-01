# URL Shortener REST API

A minimal HTTP service that accepts a URL, generates a short code, stores the mapping in memory, and redirects requests to the original URL.

## Run

```bash
go test ./...
go run .
```

The service listens on port 8080.

## Create a short URL

```bash
curl -i -X POST http://localhost:8080/shorten \
  -H 'Content-Type: application/json' \
  -d '{"url":"https://example.com/articles/go"}'
```

Example response:

```json
{
  "code": "a1B2c3D",
  "short_path": "/a1B2c3D"
}
```

Follow the returned path:

```bash
curl -i http://localhost:8080/a1B2c3D
```

## Endpoints

| Method | Path | Behavior |
|---|---|---|
| GET | `/healthz` | Health check |
| POST | `/shorten` | Validate URL and create short code |
| GET | `/{code}` | Redirect to original URL |

## Concepts

- RESTful HTTP handlers using the standard library.
- Request size limit and strict JSON decoding.
- URL scheme validation to allow only HTTP and HTTPS.
- Cryptographically generated short codes.
- Mutex-protected in-memory map for concurrent requests.
- Handler tests using `net/http/httptest`.

## Limitations and production extensions

This is intentionally a small demonstration. Data is lost on restart; storage is in-memory; there is no authentication, abuse prevention, analytics, expiry, or custom domain support. A production version should use a persistent database, enforce ownership and rate limits, add collision-safe atomic inserts, instrument the service, and define retention and abuse policies.
