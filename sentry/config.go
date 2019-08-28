package sentry

import (
	sdk "github.com/getsentry/sentry-go"
	"os"
	"strconv"
	"time"
)

//type aliases for external use
type Options sdk.ClientOptions

var DefaultConfig Options
var Timeout time.Duration

func InitDefaultConfig() {

	DefaultConfig = Options{
		Dsn:         os.Getenv("SENTRY_DSN"),
		Environment: os.Getenv("APP_ENV"),
	}

	Timeout = getConnectionTimeout(15)

	if isSyncTransport() {
		sentrySyncTransport := sdk.NewHTTPSyncTransport()
		sentrySyncTransport.Timeout = Timeout
		DefaultConfig.Transport = sentrySyncTransport
	}
}

func getConnectionTimeout(defaultTimeout time.Duration) time.Duration {
	value := os.Getenv("SENTRY_TIMEOUT")

	if t, err := strconv.ParseInt(value, 10, 64); err == nil && t > 0 {
		return time.Second * time.Duration(t)
	}

	return time.Second * defaultTimeout
}

func isSyncTransport() bool {
	if sync, err := strconv.ParseBool(os.Getenv("SENTRY_SYNC_DELIVERY")); err == nil && sync {
		return true
	}

	return false
}
