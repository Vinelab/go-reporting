package reporting

import (
	sdk "github.com/getsentry/sentry-go"
	"github.com/imdario/mergo"
	"github.com/Vinelab/go-reporting/sentry"
	log "github.com/sirupsen/logrus"
)

type services struct {
	Sentry string
}

//services supported by module
var Services = services{
	Sentry: "sentry",
}

//panic loggers list for each supported service
var panicLoggers = map[string]func(err interface{}){
	Services.Sentry: sentry.LogPanic,
}

//list of registered services for logging
var registeredServices = make(map[string]bool)

//register sentry hook for logrus
func RegisterSentry(level log.Level, config sentry.Options, injectors ...sentry.Injector) error {

	sentry.InitDefaultConfig()

	//merge config with defaults
	err := mergo.Merge(&config, &sentry.DefaultConfig)
	if err != nil {
		return err
	}

	setupSentryInjectors(&config, injectors)

	//init sentry client
	err = sdk.Init(sdk.ClientOptions(config))
	if err != nil {
		return err
	}

	//register hook for logrus
	hook := sentry.NewHook(level)
	log.AddHook(hook)

	//add to registered services list
	registeredServices[Services.Sentry] = true

	return nil
}

//log panic error for all registered services and repanic
func LogPanic() {
	err := recover()

	if err != nil {
		//log panic
		for service := range registeredServices {
			panicLoggers[service](err)
		}
		//repanic
		panic(err)
	}
}

func setupSentryInjectors(config *sentry.Options, injectors []sentry.Injector) {
	config.BeforeSend = func(event *sdk.Event, hint *sdk.EventHint) *sdk.Event {
		for _, injector := range injectors {
			injector.Inject(event)
		}

		return event
	}
}
