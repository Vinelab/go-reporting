package sentry

import (
	sdk "github.com/getsentry/sentry-go"
)

type Breadcrumb struct {
	Category string
	Message  string
	Data     map[string]interface{}
}

//wrapper for AddBreadcrumb function from sdk
func AddBreadcrumb(breadcrumb Breadcrumb) {

	b := sdk.Breadcrumb{
		Category: breadcrumb.Category,
		Message:  breadcrumb.Message,
		Data:     breadcrumb.Data,
		Level:    sdk.LevelInfo,
	}

	sdk.AddBreadcrumb(&b)
}

//capture and log error in case of panic
func LogPanic(err interface{}) {
	hub := sdk.CurrentHub().Clone()
	hub.Recover(err)
	hub.Flush(Timeout)
}

//flush queued sentry events
func Flush() {
	sdk.CurrentHub().Flush(Timeout)
}
