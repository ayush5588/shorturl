package main

import (
	"errors"
	"net/http"

	"github.com/ayush5588/shorturl/internal"
	"github.com/gin-gonic/gin"
)

var (
	// ErrURLFieldEmpty ...
	ErrURLFieldEmpty = errors.New("no url was provided to shorten")
)

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
		err := c.BindJSON(&url)
		if err != nil {
			logger.Errorw("error in unmarshalling the req body", "err", err)
			c.JSON(http.StatusBadRequest, gin.H{"message": "Please try again."})
			return
		}

		if url.OriginalURL == "" {
			logger.Errorw("invalid req body", "err", ErrURLFieldEmpty)
			c.JSON(http.StatusBadRequest, gin.H{"message": ErrURLFieldEmpty.Error() + ". Please provide url."})
			return
		}

		err = url.URLHandler(c, logger)
		if err != nil {
			if errors.Is(err, internal.ErrNotSupportedMethod) {
				logger.Error(internal.ErrNotSupportedMethod)
				c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request method."})
				return
			}
			logger.Errorw("internal error", "err", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Please try again after some time."})
			return
		}

	})

	return router
}

func main() {
	r := setupRouter()
	r.Run(":8080")
}
