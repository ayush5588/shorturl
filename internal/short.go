package internal

import (
	"encoding/json"
	"errors"

	"github.com/ayush5588/shorturl/db"
	"github.com/ayush5588/shorturl/internal/pkg/algo"
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

// URLInfo ...
type URLInfo struct {
	OriginalURL string `json:"originalURL" redis:"originalURL"`
	Alias       string `json:"alias" redis:"alias"`
}

var (
	// ErrNotSupportedMethod ...
	ErrNotSupportedMethod = errors.New("not supported method")
)

var (
	domain         = "http://localhost:8080/"
	origToShortKey = "original:to:short"
	shortToOrigKey = "short:to:original"
)

// URLHandler ...
func (u *URL) URLHandler(c *gin.Context, logger *zap.SugaredLogger) error {
	redisClient, err := db.NewRedisConnection()
	if err != nil {
		return err
	}
	u.client = redisClient

	switch c.Request.Method {
	case "GET", "":
		// Handle URL redirect request
		logger.Info("Inside GET of URLHandler")
		return nil
	case "PUT":
		// Handle ShortURL generate request
		logger.Info("Inside PUT of URLHandler", u.OriginalURL)
		return u.shortenURLHandler(logger)
	default:
		return ErrNotSupportedMethod
	}
}

func (u *URL) shortenURLHandler(logger *zap.SugaredLogger) error {
	origURL := u.OriginalURL

	// Check in db if there exist an entry for the given originalURL
	val, err := u.client.HGet(origToShortKey, origURL).Result()
	if err != nil && err != redis.Nil {
		logger.Errorf("error in getting the value from db for %s key & %s field", origToShortKey, origURL)
		return err
	}

	if val != "" {
		logger.Infof("value exist for originalURL: %s", origURL)
		u.ShortURL = domain + val
		return nil
	}

	// This is a first time shortening request for the url
	uid, err := algo.Hashing(origURL)
	if err != nil {
		logger.Errorf("error in getting uid from hashing of originalURL: %s", origURL)
		return err
	}

	urlMap := URLInfo{
		OriginalURL: origURL,
		Alias:       u.Alias,
	}

	urlMapbytes, err := json.Marshal(urlMap)
	if err != nil {
		logger.Errorf("error in marshaling of urlMap: %+v", urlMap)
		return err
	}

	_, err = u.client.HSet(shortToOrigKey, uid, urlMapbytes).Result()
	if err != nil {
		logger.Errorf("error in entering urlMap: %+v for uid: %s", urlMap, uid)
		return err
	}

	_, err = u.client.HSet(origToShortKey, origURL, uid).Result()
	if err != nil {
		logger.Errorf("error in entering mapping between originalURL: %s & uid: %s", origURL, uid)
		return err
	}

	u.ShortURL = domain + uid

	logger.Info("Success URLHandler operation")

	return nil

}
