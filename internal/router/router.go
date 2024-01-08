package router

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/ayush5588/shorturl/db"
	"github.com/ayush5588/shorturl/internal"
	"github.com/ayush5588/shorturl/internal/handler"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var tmplt *template.Template

var (
	domain = os.Getenv("DOMAIN_NAME")
)

func preShortenValidation(c *gin.Context, url *handler.URL, logger *zap.SugaredLogger) error {

	url.OriginalURL = c.PostForm("originalURL")
	url.Alias = c.PostForm("alias")

	if url.OriginalURL == "" {
		return internal.ErrEmptyURLField
	}

	if !internal.IsValidURL(logger, url.OriginalURL) {
		return internal.ErrInvalidURL
	}

	if !internal.IsValidAlias(logger, url.Alias) {
		return internal.ErrInvalidAlias
	}

	return nil
}

// SetupRouter ...
func SetupRouter() *gin.Engine {
	logger := internal.GetLogger()

	// Establish a redis connection at the start of the router setup
	redisClient, err := db.NewRedisConnection(logger)
	if err != nil {
		log.Fatal(err)
	}

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
		Method: POST
		Path: /short
		Body: {
			OriginalURL string (Required)
			Alias string (To be handled later)
		}
		Definition: Returns a shortened URL for the given URL
	*/
	router.POST("/short", func(c *gin.Context) {
		var url handler.URL
		url.Client = redisClient
		err := preShortenValidation(c, &url, logger)
		if err != nil {
			internal.HandleError(c, err, "preShortenValidation", logger)
			return
		}

		// URLHandler handles the url shortening GET (i.e. redirect) & POST (i.e. shortening) operation
		err = url.URLHandler(c, logger)
		if err != nil {
			internal.HandleError(c, err, "URLHandler", logger)
			return
		}

		var urlExistMessage string
		if url.URLExist {
			urlExistMessage = "Shortened URL of the given URL already exists"
		}
		if domain == "" {
			domain = "http://localhost:8080/"
		}

		c.HTML(http.StatusOK, "index.html", gin.H{"longURLOutput": urlExistMessage, "output": domain + url.UID})
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
		url := handler.URL{
			UID:      id,
			Database: handler.Database{Client: redisClient},
		}

		err := url.URLHandler(c, logger)
		if err != nil {
			internal.HandleError(c, err, "URLHandler Redirect operation", logger)
			return
		}

		c.Redirect(http.StatusFound, url.OriginalURL)
		return
	})

	return router
}
