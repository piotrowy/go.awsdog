package main

import (
	"flag"
	"os"
	"os/signal"
)

type (
	Service interface {
		Start()
		Close()
	}
)

func main() {
	dynamoPK := flag.String("dynamoPK", "", "Dynamo primary key to AWS dynamo database.")

	var (
		settings = GetSettings(*dynamoPK)
		router   = NewEventRouter()
	)

	errorLogger := NewErrorLogger(len(settings.Services)*5, DefaultLogPath)
	eventLogger := NewEventLogger(len(settings.Services)*2, DefaultLogPath)

	serviceDownHandler := NewServiceDownHandler(router, eventLogger.GetRoute(), len(settings.Services))
	router.RegisterHandler(serviceDownHandler)

	serviceRestartHandler := NewServiceRestartHandler(router, eventLogger.GetRoute(), len(settings.Services))
	router.RegisterHandler(serviceRestartHandler)

	watchDog := NewWatchDog(router, WatchDogConfig{
		ConfigAWS: ConfigAWS{},
		Settings:  settings,
	})
	router.RegisterHandler(&watchDog)

	settingsLoader := NewSettingsLoader(router, errorLogger.GetRoute(), ConfigAWS{})

	services := []Service{errorLogger, eventLogger, serviceDownHandler, serviceRestartHandler, settingsLoader, &watchDog}
	for _, service := range services {
		go service.Start()
	}

	defer func() {
		for _, service := range services {
			service.Close()
		}
	}()

	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, os.Kill)
	<-sig

	//todo go-daemon
}
