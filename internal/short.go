package internal

import (
	"errors"

	"github.com/gin-gonic/gin"
)

// URL ...
type URL struct {
	OriginalURL string `json:"originalURL"`
	ShortURL    string `json:"shortURL"`
	Alias       string `json:"alias"`
}

var (
	// ErrNotSupportedMethod ...
	ErrNotSupportedMethod = errors.New("not supported method")
)

var (
	domain = "http://localhost:8080/"
)

// URLHandler ...
func (u *URL) URLHandler(c *gin.Context) error {
	switch c.Request.Method {
	case "GET", "":
		// Handle URL redirect request
	case "PUT":
		// Handle ShortURL generate request
	default:
		return ErrNotSupportedMethod
	}

	return nil
}
