package internal

import (
	"net/url"
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
	return (err == nil || u.Scheme == "" || u.Host == "")

}
