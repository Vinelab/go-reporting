package sentry

import (
	"github.com/getsentry/sentry-go"
	"github.com/sirupsen/logrus"
)

var (
	levelsMap = map[logrus.Level]sentry.Level{
		logrus.PanicLevel: sentry.LevelFatal,
		logrus.FatalLevel: sentry.LevelFatal,
		logrus.ErrorLevel: sentry.LevelError,
		logrus.WarnLevel:  sentry.LevelWarning,
		logrus.InfoLevel:  sentry.LevelInfo,
		logrus.DebugLevel: sentry.LevelDebug,
		logrus.TraceLevel: sentry.LevelDebug,
	}
)

type Hook struct {
	level logrus.Level
}

func (hook *Hook) Levels() []logrus.Level {
	levels := make([]logrus.Level, 0, len(logrus.AllLevels))

	for _, l := range logrus.AllLevels {
		if hook.level >= l {
			levels = append(levels, l)
		}
	}

	return levels
}

func (hook *Hook) Fire(entry *logrus.Entry) error {
	exceptions := []sentry.Exception{}

	if err, ok := entry.Data[logrus.ErrorKey].(error); ok && err != nil {
		stacktrace := sentry.ExtractStacktrace(err)
		if stacktrace == nil {
			stacktrace = sentry.NewStacktrace()
		}
		exceptions = append(exceptions, sentry.Exception{
			Type:       entry.Message,
			Value:      err.Error(),
			Stacktrace: stacktrace,
		})
	}

	event := sentry.Event{
		Level:     levelsMap[entry.Level],
		Message:   entry.Message,
		Extra:     map[string]interface{}(entry.Data),
		Exception: exceptions,
	}

	hub := sentry.CurrentHub().Clone()
	hub.CaptureEvent(&event)

	Flush()

	return nil
}

func NewHook(level logrus.Level) *Hook {

	hook := Hook{
		level: level,
	}

	return &hook
}
