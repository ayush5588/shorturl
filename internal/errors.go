package internal

import (
	"errors"
	"net/http"

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

// HandleError ...
func HandleError(c *gin.Context, err error, errorKey string, logger *zap.SugaredLogger) {
	if errors.Is(err, ErrEmptyURLField) {
		c.HTML(http.StatusBadRequest, "index.html", gin.H{"longURLOutput": ErrEmptyURLField.Error()})
		return
	} else if errors.Is(err, ErrInvalidURL) {
		c.HTML(http.StatusBadRequest, "index.html", gin.H{"longURLOutput": ErrInvalidURL.Error()})
		return
	} else if errors.Is(err, ErrInvalidAlias) {
		c.HTML(http.StatusBadRequest, "index.html", gin.H{"longURLOutput": ErrInvalidAlias.Error()})
		return
	} else if errors.Is(err, ErrNotSupportedMethod) {
		logger.Error(ErrNotSupportedMethod)
		c.JSON(http.StatusBadRequest, gin.H{"message": ErrNotSupportedMethod.Error()})
		return
	} else if errors.Is(err, ErrAliasExist) {
		logger.Error(ErrAliasExist)
		c.HTML(http.StatusBadRequest, "index.html", gin.H{"longURLOutput": ErrAliasExist.Error()})
		return
	} else if errors.Is(err, ErrOriginalURLDoesNotExist) {
		c.JSON(http.StatusNotFound, gin.H{"message": ErrOriginalURLDoesNotExist.Error()})
		return
	}
	logger.Errorf("%s failed with error: %s", errorKey, err.Error())
	c.HTML(http.StatusInternalServerError, "index.html", gin.H{"longURLOutput": "Please try again after some time"})
	return
}
