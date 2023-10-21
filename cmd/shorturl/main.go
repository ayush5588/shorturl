package main

import (
	"net/http"

	"github.com/ayush5588/shorturl/internal"
	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	logger := internal.GetLogger()

	router := gin.Default()

	router.GET("/healthz", func(c *gin.Context) {
		logger.Infof("Successfully served GET /healthz request")
		c.JSON(http.StatusOK, gin.H{"message": "Server is healthy"})
	})

	return router
}

func main() {
	r := setupRouter()
	r.Run(":8080")
}
