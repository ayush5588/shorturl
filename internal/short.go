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
	// ErrAliasExist ...
	ErrAliasExist = errors.New("given alias already exist")
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

// checkAliasExist ...
func (u *URL) checkAliasExist(alias string) (bool, error) {
	var exist bool
	valExist, err := u.client.HGet(shortToOrigKey, alias).Result()
	if err != nil && err != redis.Nil {
		return exist, err
	}
	if valExist != "" {
		exist = true
	}

	return exist, nil
}

func (u *URL) shortenURLHandler(logger *zap.SugaredLogger) error {
	origURL := u.OriginalURL

	// Check in db if there exist an entry for the given originalURL
	val, err := u.client.HGet(origToShortKey, origURL).Result()
	if err != nil && err != redis.Nil {
		logger.Errorf("error in getting the value from db for %s field in %s key ", origURL, origToShortKey)
		return err
	}

	// Shortened URL already exist for the user given original URL
	if val != "" {
		logger.Infof("value exist for originalURL: %s", origURL)
		u.UID = val
		return nil
	}

	// This is a first time shortening request for the url
	var uid string

	if u.Alias != "" {
		uid = u.Alias
		// check if the given alias is unique
		aliasExist, err := u.checkAliasExist(uid)
		if err != nil && err != redis.Nil {
			logger.Errorf("error in getting the alias value from db for %s field in %s key ", uid, shortToOrigKey)
			return err
		}
		// Given alias exist
		if aliasExist {
			logger.Errorf("given alias %s exist", uid)
			return ErrAliasExist
		}
	} else {
		// Generate a unique id for the given original URL as user has not given any alias
		uid = algo.UniqueID(origURL)
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

	// Store mapping between unique id (either valid alias or generated id) AND original URL, alias(if given)
	_, err = u.client.HSet(shortToOrigKey, uid, urlMapbytes).Result()
	if err != nil {
		logger.Errorf("error in entering urlMap: %+v for uid: %s", urlMap, uid)
		return err
	}

	// Store Original URL mapping with the unique id
	_, err = u.client.HSet(origToShortKey, origURL, uid).Result()
	if err != nil {
		logger.Errorf("error in entering mapping between originalURL: %s & uid: %s", origURL, uid)
		return err
	}

	u.UID = uid

	logger.Info("Success URLHandler operation")

	return nil

}
