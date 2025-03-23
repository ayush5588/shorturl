# ShortURL - URL Shortener Service

A modern, fast, and reliable URL shortening service built with Go and Redis. This service allows you to create shortened URLs with optional custom aliases, perfect for sharing links in a more manageable way.

## Features

- Create shortened URLs instantly
- Custom alias support
- Fast redirects using Redis
- URL validation
- Clean and simple web interface
- Docker support
- Detailed logging

## Tech Stack

- Go (Gin Web Framework)
- Redis
- Docker & Docker Compose
- HTML/Templates

## Prerequisites

- Go 1.x
- Docker and Docker Compose
- Redis (if running locally)

## Quick Start

1. Clone the repository:
   ```bash
   git clone https://github.com/ayush5588/shorturl.git
   cd shorturl
   ```

2. Run with Docker Compose:
   ```bash
   docker-compose up
   ```

   This will start both the web server and Redis container.

3. Access the application at `http://localhost:8080`

## Environment Variables

- `REDIS_HOST`: Redis host address (default: "redis")
- `REDIS_PORT`: Redis port (default: 6379)
- `DOMAIN_NAME`: Domain name for shortened URLs (default: "http://localhost:8080/")

## API Endpoints

- `GET /`: Home page with URL shortening interface
- `POST /short`: Create a shortened URL
  - Parameters:
    - `originalURL`: The URL to shorten (required)
    - `alias`: Custom alias (optional)
- `GET /:id`: Redirect to the original URL
- `GET /healthz`: Health check endpoint

## URL Shortening Rules

- Original URLs must be valid and include scheme (http/https)
- Custom aliases:
  - Must be unique
  - Maximum length: 15 characters
  - Cannot contain special characters: !@#$%^&*()+={}[]|`/?.>,<:;'
  - Optional (system will generate a unique ID if not provided)

## Development

To run locally without Docker:

1. Start Redis server
2. Set environment variables
3. Run the application:
   ```bash
   go run main.go
   ```

## Testing

Run the tests:
```bash
go test ./...
```

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.
