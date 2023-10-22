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
	UID         string `json:"uid"`
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
	// ErrOriginalURLDoesNotExist ...
	ErrOriginalURLDoesNotExist = errors.New("original url for the given short url does not exist")
)

var (
	// Domain ...
	Domain         = "http://localhost:8080/"
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
		return u.redirectToOriginalURL(logger)
	case "POST":
		// Handle ShortURL generate request
		logger.Info("Inside PUT of URLHandler", u.OriginalURL)
		return u.shortenURLHandler(logger)
	default:
		return ErrNotSupportedMethod
	}
}

func (u *URL) redirectToOriginalURL(logger *zap.SugaredLogger) error {
	uid := u.UID

	// Check for the original URL in redis shortToOrigKey
	val, err := u.client.HGet(shortToOrigKey, uid).Result()
	if err != nil && err != redis.Nil {
		logger.Error("error in getting the value from db for %s field in %s key ", uid, shortToOrigKey)
		return err
	}

	// Original URL exist for the given short URL
	if val != "" {
		var urlInfo URLInfo
		err := json.Unmarshal([]byte(val), &urlInfo)
		if err != nil {
			return err
		}
		logger.Infof("original URL for uid: %s is %s", uid, urlInfo.OriginalURL)
		u.OriginalURL = urlInfo.OriginalURL
		return nil
	}

	return ErrOriginalURLDoesNotExist

}

func (u *URL) shortenURLHandler(logger *zap.SugaredLogger) error {
	origURL := u.OriginalURL

	// Check in db if there exist an entry for the given originalURL
	val, err := u.client.HGet(origToShortKey, origURL).Result()
	if err != nil && err != redis.Nil {
		logger.Errorf("error in getting the value from db for %s field in %s key ", origURL, origToShortKey)
		return err
	}

	if val != "" {
		logger.Infof("value exist for originalURL: %s", origURL)
		u.UID = val
		return nil
	}

	// This is a first time shortening request for the url
	uid := algo.UniqueID(origURL)

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

	u.UID = uid

	logger.Info("Success URLHandler operation")

	return nil

}
