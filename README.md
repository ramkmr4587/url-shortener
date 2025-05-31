📦 URL Shortener Service
========================

This is a simple, API-only URL Shortener service written in Golang as part of a Product Engineer assignment.

Features
-----------
- Shortens long URLs via REST API
- Returns the same short URL for duplicate original URLs
- Redirects short URLs to their original links
- Stores URL mappings in memory
- Provides metrics for top 3 most frequently shortened domains
- [BONUS] Dockerized for easy deployment

Assignment Reference
-----------------------
This project was built in response to the following requirements:

Build a simple URL shortener service with a REST API that:
- Accepts a URL and returns a shortened version
- Remembers and returns the same short URL for repeated inputs
- Redirects short URLs to the original URL
- Maintains in-memory storage
- Exposes a metrics API for top 3 domains shortened most frequently
- [BONUS] Provide Dockerfile for containerized execution

Technologies Used
---------------------
- Golang — Core language
- net/http — HTTP server
- chi router (or your router of choice)
- Go modules — Dependency management
- Docker — Containerization (Bonus)

Installation & Run
----------------------

▶️ Run using Go:
    git clone https://github.com/your-username/url-shortener.git
    cd url-shortener
    go mod tidy
    go run main.go

Run using Docker:
    docker build -t url-shortener .
    docker run -p 8080:8080 url-shortener

API Endpoints
-----------------

1. Shorten URL
   POST /shorten
   Body:
       {
         "url": "https://example.com/some/very/long/url"
       }

   Response:
       {
         "short_url": "http://localhost:8080/abc123"
       }

2. Redirect URL
   GET /abc123
   → Redirects to original URL.

3. Get Metrics
   GET /metrics

   Response:
       {
         "top_domains": {
           "udemy.com": 6,
           "youtube.com": 4,
           "wikipedia.org": 2
         }
       }

Tests
--------
To run unit tests:

    go test ./...

Project Structure
--------------------
url-shortener/
├── main.go
├── handlers/
│   └── url_handler.go
├── storage/
│   └── memory_store.go
├── utils/
│   └── shortener.go
├── Dockerfile
├── url-shortener-service.yaml
└── README.txt

Best Practices Followed
--------------------------
- Modular and readable code
- Well-structured project layout
- Proper naming conventions
- Unit tests included
- Docker support for easy deployment

👤 Author
---------
Ramkumar
Golang Developer