package main

import (
	"fmt"
	"time"
)

type (
	ServiceDown struct {
		Name          string
		Time          time.Time
		NumOfAttempts int
		NumOfSecWait  int64
	}

	ServiceDownHandler struct {
		router   EventRouter
		logCh    chan<- interface{}
		eventBus chan Event
		close    chan bool
	}
)

func (s ServiceDown) String() string {
	return fmt.Sprintf("Service %v is down.", s.Name)
}

func NewServiceDownHandler(router EventRouter, logCh chan<- interface{}, bufferSize int) EventHandler {
	return ServiceDownHandler{
		router:   router,
		logCh:    logCh,
		eventBus: make(chan Event, bufferSize),
		close:    make(chan bool),
	}
}

func (sdh ServiceDownHandler) Start() {
	for {
		select {
		case e := <-sdh.eventBus:
			go func() {
				if serviceDown, ok := e.(ServiceDown); ok {
					sdh.logCh <- serviceDown
					sdh.notify(serviceDown)
					sdh.router.Route(ServiceRestartReq{
						Name:          serviceDown.Name,
						NumOfAttempts: serviceDown.NumOfAttempts,
						NumOfSecWait:  serviceDown.NumOfSecWait,
					})
				}
			}()
		case <-sdh.close:
			return
		}
	}
}

func (sdh ServiceDownHandler) CanHandle(e Event) bool {
	_, ok := e.(ServiceDown)
	return ok
}

func (sdh ServiceDownHandler) GetRoute() chan<- Event {
	return sdh.eventBus
}

func (sdh ServiceDownHandler) Close() {
	sdh.close <- true
}

func (sdh ServiceDownHandler) notify(event ServiceDown) {
	//todo notifier
}
