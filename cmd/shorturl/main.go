package main

import (
	"errors"
	"net/http"
	"reflect"

	"github.com/ayush5588/shorturl/internal"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var (
	// ErrEmptyReqBody ...
	ErrEmptyReqBody = errors.New("request body cannot be empty")
	// ErrEmptyURLField ...
	ErrEmptyURLField = errors.New("no url was provided to shorten")
	// ErrUnmarshallingReqBody ...
	ErrUnmarshallingReqBody = errors.New("error in unmarshalling the request body")
)

func preShortenValidation(c *gin.Context, url internal.URL, logger *zap.SugaredLogger) error {
	if reflect.DeepEqual(url, internal.URL{}) {
		logger.Error(ErrEmptyReqBody)
		return ErrEmptyReqBody
	}

	err := c.BindJSON(&url)
	if err != nil {
		logger.Errorw(ErrUnmarshallingReqBody.Error(), "err", err)
		return ErrUnmarshallingReqBody
	}

	if url.OriginalURL == "" {
		logger.Errorw("invalid req body", "err", ErrEmptyURLField)
		c.JSON(http.StatusBadRequest, gin.H{"message": ErrEmptyURLField.Error() + ". Please provide url"})
		return ErrEmptyURLField
	}
	return nil
}

func setupRouter() *gin.Engine {
	logger := internal.GetLogger()

	router := gin.Default()

	/*
		Method: GET
		Path: /healthz
		Definition: Represents server health
	*/
	router.GET("/healthz", func(c *gin.Context) {
		logger.Infof("Successfully served GET /healthz request")
		c.JSON(http.StatusOK, gin.H{"message": "Server is healthy"})
		return
	})

	/*
		Method: PUT
		Path: /short
		Body: {
			OriginalURL string (Required)
			Alias string (To be handled later)
		}
		Definition: Returns a shortened URL for the given URL
	*/
	router.PUT("/short", func(c *gin.Context) {
		var url internal.URL

		err := preShortenValidation(c, url, logger)
		if err != nil {
			if errors.Is(err, ErrEmptyReqBody) {
				c.JSON(http.StatusBadRequest, gin.H{"message": ErrEmptyReqBody})
				return
			} else if errors.Is(err, ErrUnmarshallingReqBody) {
				c.JSON(http.StatusBadRequest, gin.H{"message": "Please try again"})
				return
			} else if errors.Is(err, ErrEmptyURLField) {
				c.JSON(http.StatusBadRequest, gin.H{"message": ErrEmptyURLField.Error() + ". Please provide url"})
				return
			}
			logger.Errorw("preShortenValidation failed", "err", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Please try again after some time"})
			return

		}

		// URLHandler handles the url shortening operation
		err = url.URLHandler(c, logger)
		if err != nil {
			if errors.Is(err, internal.ErrNotSupportedMethod) {
				logger.Error(internal.ErrNotSupportedMethod)
				c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request method."})
				return
			}
			logger.Errorw("URLHandler operation failed.", "err", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Please try again after some time."})
			return
		}
		c.JSON(http.StatusOK, gin.H{"originalURL": url.OriginalURL, "shortURL": url.ShortURL})
		return

	})

	router.GET("/:id", func(ctx *gin.Context) {

	})

	return router
}

func main() {
	r := setupRouter()
	r.Run(":8080")
}
