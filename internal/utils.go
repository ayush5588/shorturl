package internal

import (
	"net/url"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// GetLogger ...
func GetLogger() *zap.SugaredLogger {
	// Setup logger
	zapProdConfig := zap.NewProductionConfig()
	// Modify the logger to show rfc3339 date & time format
	zapProdConfig.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)

	zapProd, _ := zapProdConfig.Build()

	logger := zapProd.Sugar()

	return logger
}

// IsValidURL ...
func IsValidURL(logger *zap.SugaredLogger, rawURL string) bool {
	_, err := url.ParseRequestURI(rawURL)
	if err != nil {
		logger.Errorw("error in parsing given url", "err", err)
		return false
	}

	u, err := url.Parse(rawURL)
	if err != nil {
		logger.Errorw("error in parsing given url", "err", err)
	}

	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}

	return true

}

// IsValidAlias ...
func IsValidAlias(logger *zap.SugaredLogger, alias string) bool {
	if alias == "" {
		return true
	}

	if len(alias) > 15 {
		return false
	}

	notAllowed := "!@#$%^&*()+={}[]|`/?.>,<:;'"

	for _, ch := range alias {
		if strings.ContainsRune(notAllowed, ch) {
			logger.Errorf("invalid alias. Contains %c special character", ch)
			return false
		}
	}

	return true
}
