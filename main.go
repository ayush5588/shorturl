package main

import (
	"errors"
	"html/template"
	"net/http"
	"path"

	"github.com/ayush5588/shorturl/internal"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var (
	// ErrEmptyReqBody ...
	ErrEmptyReqBody = errors.New("request body cannot be empty")
	// ErrEmptyURLField ...
	ErrEmptyURLField = errors.New("URL field cannot be empty")
	// ErrInvalidURL ...
	ErrInvalidURL = errors.New("invalid url")
	// ErrInvalidAlias ...
	ErrInvalidAlias = errors.New("invalid alias")
)

var (
	tmplt *template.Template
)

func preShortenValidation(c *gin.Context, url *internal.URL, logger *zap.SugaredLogger) error {

	url.OriginalURL = c.PostForm("originalURL")
	url.Alias = c.PostForm("alias")

	if url.OriginalURL == "" {
		logger.Errorw("invalid req body", "err", ErrEmptyURLField)
		return ErrEmptyURLField
	}

	if !internal.IsValidURL(logger, url.OriginalURL) {
		return ErrInvalidURL

	}

	if !internal.IsValidAlias(logger, url.Alias) {
		return ErrInvalidAlias
	}

	return nil
}

func setupRouter() *gin.Engine {
	logger := internal.GetLogger()

	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.Static("/templates", "./templates/")

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
		Method: GET
		Path: /
		Definition: Serves the home page
	*/

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
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
	router.POST("/short", func(c *gin.Context) {
		var url internal.URL		
		err := preShortenValidation(c, &url, logger)
		if err != nil {
			if errors.Is(err, ErrEmptyURLField) {
				c.HTML(http.StatusBadRequest, "index.html", gin.H{"longURLOutput": ErrEmptyURLField.Error()})
				return
			} else if errors.Is(err, ErrInvalidURL) {
				c.HTML(http.StatusBadRequest, "index.html", gin.H{"longURLOutput": ErrInvalidURL.Error()})
				return
			} else if errors.Is(err, ErrInvalidAlias) {
				c.HTML(http.StatusBadRequest, "index.html", gin.H{"longURLOutput": ErrInvalidAlias.Error()})
				return
			}
			logger.Errorw("preShortenValidation failed", "err", err)
			c.HTML(http.StatusInternalServerError, "index.html", gin.H{"longURLOutput": "Please try again after some time"})
			return

		}

		// URLHandler handles the url shortening GET (i.e. redirect) & POST (i.e. shortening) operation
		err = url.URLHandler(c, logger)
		if err != nil {
			if errors.Is(err, internal.ErrNotSupportedMethod) {
				logger.Error(internal.ErrNotSupportedMethod)
				c.JSON(http.StatusBadRequest, gin.H{"message": internal.ErrNotSupportedMethod.Error()})
				return
			} else if errors.Is(err, internal.ErrAliasExist) {
				logger.Error(internal.ErrAliasExist)
				c.HTML(http.StatusBadRequest, "index.html", gin.H{"longURLOutput": internal.ErrAliasExist.Error()})
				return
			}
			logger.Errorw("URLHandler operation failed.", "err", err)
			c.HTML(http.StatusInternalServerError, "index.html", gin.H{"longURLOutput": "Please try again after some time"})
			return
		}
		var urlExistMessage string
		if url.URLExist {
			urlExistMessage = "Shortened URL of the given URL already exists"
		}
		c.HTML(http.StatusOK, "index.html", gin.H{"longURLOutput": urlExistMessage, "output": internal.Domain + url.UID})
		return

	})

	/*
		Method: GET
		Path: /:id
		Definition: Redirects the shortURL to originalURL
	*/
	router.GET("/:id", func(c *gin.Context) {
		reqURL := c.Request.RequestURI
		id := path.Base(reqURL)
		url := internal.URL{UID: id}

		err := url.URLHandler(c, logger)
		if err != nil {
			if errors.Is(err, internal.ErrOriginalURLDoesNotExist) {
				c.JSON(http.StatusNotFound, gin.H{"message": internal.ErrOriginalURLDoesNotExist.Error()})
				return
			}
			logger.Errorw("URLHandler Redirect operation failed", "err", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Please try again after some time"})
			return
		}
		c.Redirect(http.StatusTemporaryRedirect, url.OriginalURL)
		return
	})

	return router
}

func main() {
	r := setupRouter()
	r.Run(":8080")
}
