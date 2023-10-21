package internal

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

// Database ...
type Database struct {
	client *redis.Client
}

// URL ...
type URL struct {
	OriginalURL string `json:"originalURL"`
	ShortURL    string `json:"shortURL"`
	Alias       string `json:"alias"`
	Database
}

var (
	// ErrNotSupportedMethod ...
	ErrNotSupportedMethod = errors.New("not supported method")
)

var (
	domain = "http://localhost:8080/"
)

// URLHandler ...
func (u *URL) URLHandler(c *gin.Context, logger *zap.SugaredLogger) error {
	switch c.Request.Method {
	case "GET", "":
		// Handle URL redirect request
		logger.Info("Inside GET of URLHandler")
		return nil
	case "PUT":
		// Handle ShortURL generate request
		logger.Info("Inside PUT of URLHandler", u.OriginalURL)

		return nil
	default:
		return ErrNotSupportedMethod
	}
}