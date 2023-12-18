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
	// ErrNotSupportedMethod ...
	ErrNotSupportedMethod = errors.New("not supported method")
	// ErrOriginalURLDoesNotExist ...
	ErrOriginalURLDoesNotExist = errors.New("original url for the given short url does not exist")
	// ErrAliasExist ...
	ErrAliasExist = errors.New("given alias already exist")
)

// HandleError ...
func HandleError(c *gin.Context, err error, errorKey string, logger *zap.SugaredLogger) {

	switch {

	case errors.Is(err, ErrEmptyURLField):
		logger.Error(ErrEmptyURLField)
		c.HTML(http.StatusBadRequest, "index.html", gin.H{"longURLOutput": ErrEmptyURLField.Error()})
		break

	case errors.Is(err, ErrInvalidURL):
		logger.Error(ErrInvalidURL)
		c.HTML(http.StatusBadRequest, "index.html", gin.H{"longURLOutput": ErrInvalidURL.Error()})
		break

	case errors.Is(err, ErrInvalidAlias):
		logger.Error(ErrInvalidAlias)
		c.HTML(http.StatusBadRequest, "index.html", gin.H{"longURLOutput": ErrInvalidAlias.Error()})
		break

	case errors.Is(err, ErrNotSupportedMethod):
		logger.Error(ErrNotSupportedMethod)
		c.JSON(http.StatusBadRequest, gin.H{"message": ErrNotSupportedMethod.Error()})
		break

	case errors.Is(err, ErrAliasExist):
		logger.Error(ErrAliasExist)
		c.HTML(http.StatusBadRequest, "index.html", gin.H{"longURLOutput": ErrAliasExist.Error()})
		break

	case errors.Is(err, ErrOriginalURLDoesNotExist):
		logger.Error(ErrOriginalURLDoesNotExist)
		c.JSON(http.StatusNotFound, gin.H{"message": ErrOriginalURLDoesNotExist.Error()})
		break

	default:
		logger.Errorf("%s failed with error: %s", errorKey, err.Error())
		c.HTML(http.StatusInternalServerError, "index.html", gin.H{"longURLOutput": "Please try again after some time"})
		break
	}

	return
}
