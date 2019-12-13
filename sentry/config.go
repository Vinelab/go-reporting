package sentry

import (
	"os"
	"strconv"
	"strings"
	"time"

	sdk "github.com/getsentry/sentry-go"
)

//type aliases for external use
type Options sdk.ClientOptions

var DefaultConfig Options
var Timeout time.Duration

func InitDefaultConfig() {

	DefaultConfig = Options{
		Dsn:         getEnv("SENTRY_DSN", ""),
		Environment: getEnv("APP_ENV", ""),
	}

	Timeout = getConnectionTimeout(15)

	if isSyncTransport() {
		sentrySyncTransport := sdk.NewHTTPSyncTransport()
		sentrySyncTransport.Timeout = Timeout
		DefaultConfig.Transport = sentrySyncTransport
	}
}

func getConnectionTimeout(defaultTimeout time.Duration) time.Duration {
	value := getEnv("SENTRY_TIMEOUT", "")

	if t, err := strconv.ParseInt(value, 10, 64); err == nil && t > 0 {
		return time.Second * time.Duration(t)
	}

	return time.Second * defaultTimeout
}

func isSyncTransport() bool {
	if sync, err := strconv.ParseBool(getEnv("SENTRY_SYNC_DELIVERY", "")); err == nil && sync {
		return true
	}

	return false
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return strings.TrimSuffix(value, "\n")
	}

	return defaultVal
}
