package main

import (
	"time"
)

type (
	WatchDogConfig struct {
		ConfigAWS
		Settings
	}

	WatchDog struct {
		settings Settings
		router   EventRouter
		eventBus chan Event
		close    chan bool
	}
)

func NewWatchDog(router EventRouter, config WatchDogConfig) WatchDog {
	return WatchDog{
		settings: config.Settings,
		router:   router,
		eventBus: make(chan Event, len(config.Settings.Services)),
		close:    make(chan bool),
	}
}

func (w *WatchDog) Start() {
settingsChanged:
	for range time.Tick(time.Duration(int64(time.Second) * w.settings.NumOfSecCheck)) {
		select {
		case ev := <-w.eventBus:
			switch v := ev.(type) {
			case SettingsChanged:
				w.settings = v.Settings
				break settingsChanged
			case ServiceRestarted:
				w.settings.Services = append(w.settings.Services, v.Name)
			}
		case <-w.close:
			return
		default:
			w.check()
		}
	}
}

func (w *WatchDog) GetRoute() chan<- Event {
	return w.eventBus
}

func (w *WatchDog) CanHandle(e Event) bool {
	switch e.(type) {
	case SettingsChanged, ServiceRestarted:
		return true
	}
	return false
}

func (w *WatchDog) Close() {
	w.close <- true
}

func (w *WatchDog) check() {
	for _, v := range w.settings.Services {
		go func(service string) {
			if checkService(service) {
				w.router.Route(ServiceDown{
					Name:          service,
					Time:          time.Now(),
					NumOfAttempts: w.settings.NumOfAttempts,
					NumOfSecWait:  w.settings.NumOfSecWait,
				})
			}
		}(v)
	}
}

func checkService(service string) bool {
	return true //todo checker
}
