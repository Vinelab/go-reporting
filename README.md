# [![:walrus:](https://camo.githubusercontent.com/0e39ba1d71e6818ceeccd5b6f7f70bddec32e768/687474703a2f2f692e696d6775722e636f6d2f68546556776d4a2e706e67)](https://camo.githubusercontent.com/0e39ba1d71e6818ceeccd5b6f7f70bddec32e768/687474703a2f2f692e696d6775722e636f6d2f68546556776d4a2e706e67) Logging and error reporting 

This module provides hooks registrator for [Logrus](https://github.com/sirupsen/logrus) to simplify logging to external error tracking systems. Currently only [Sentry](https://sentry.io/welcome/) system supported based on [official SDK for Go](https://github.com/getsentry/sentry-go).

## Requirements

Supported Go versions are:

- 1.10
- 1.11
- 1.12

## Installation

`go-reporting` can be installed like any other Go library through `go get`:

```
$ go get github.com/Vinelab/go-reporting
```

## Configuration

### Sentry

To register Sentry hook for Logrus 

```go
package main

import (
	"github.com/Vinelab/go-reporting"
	"github.com/Vinelab/go-reporting/sentry"
	log "github.com/sirupsen/logrus"

)

func main() {
_ = reporting.RegisterSentry(
		log.InfoLevel,
		sentry.Options{},
}
```

Minimum logging level and Sentry SDK options can be specified on hook registration. 

Default Sentry SDK options are:

```go
var DefaultSentryConfig = Options{
	Dsn:         os.Getenv("SENTRY_DSN"),
	Environment: os.Getenv("APP_ENV"),
}
```

**Environment variables**

- `APP_ENV` - application environment, used by sentry events

- `SENTRY_DSN` - dsn for authorization in sentry

- `SENTRY_SYNC_DELIVERY` - true/false for sync message delivery to sentry service, concurrent async sending by default

- `SENTRY_TIMEOUT` - connection timeout to sentry server 15 sec by default

  

Sentry DSN for authorization and application environment can be passed via .env file using `SENTRY_DSN` and `APP_ENV` variables.



## Injectors

### Sentry

Module supports injectors for sentry event data using `BeforeSend` hook. This allows to pass dynamically  generated data to sentry right before sending an event. This can be useful for integrations with other modules. Currently only `TagsInjector` is supported

```go
package main

import (
    "github.com/Vinelab/go-reporting"
    "github.com/Vinelab/go-reporting/sentry"
	log "github.com/sirupsen/logrus"
    "time"
)

func main() {
_ = reporting.RegisterSentry(
    log.InfoLevel,
    sentry.SentryOptions{},
    sentry.TagInjector{
    Tags: func() map[string]string {
        return map[string]string{
            "date":  time.Now().String(),
            "exchange_rate": currensy.USDToEUR(),
        }
    }}
}

```



## Logging HTTP Request data

### Sentry

To capture http request data into sentry event `LogResponseMiddleware` should be added to the defined routes

```go
import (
	"github.com/go-chi/chi"
	"github.com/Vinelab/go-reporting/sentry"
	"net/http"
)

// Register holds the routes to be registered
// when the server starts listening
func Register() *chi.Mux {
	router := chi.NewRouter()

	//add sentry middleware
	router.Use(sentry.LogResponseMiddleware())

	router.Get("/user", handler)

	return router
}
```



## Creating Breadcrumbs

### Sentry

Sentry supports `breadcrumbs` for customization events and better application flow visibility. Module provides simple way for adding breadcrumbs to logs

```go
sentry.AddBreadcrumb(sentry.Breadcrumb{
    Category: "Log initialised",
    Message: "Sentry hook was registered in Logrus",
    Data: map[string]interface{}{
        "key" : "value",
    },
})
```


## Basic Usage

Module works with Logrus logger so logging and reporting should be performed via Logrus methods.

**Info message**

```go
log.WithFields(log.Fields{
    "product_id": 10,
    "description": "Amazon Kindle",
}).Info("Product was purchased")
```

Fields will be recorded as `extra` data in sentry event

**Error message**

```go
log.WithError(err).Error("Something went wrong")
```

Stacktrace and error type (if any) will be automatically extracted into sentry event 



## Handle Panic

Module provides recovery function for logging panic errors into all registered error reporting services. Function logs panic error as `Fatal` and re-panic at the end. Deffer call should be added to main routine and to any other goroutine since panic is thrown in goroutine scope only.

```go
package main

import (
	"github.com/Vinelab/go-reporting"

)

func main() {
    defer reporting.LogPanic()
    
    //application code here

}
```


