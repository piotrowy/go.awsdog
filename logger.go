package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

const (
	DefaultLogPath = "/var/log/go.awsdog/"
)

type (
	Logger interface {
		Start()
		GetRoute() chan<- interface{}
		Close()
	}

	eventLogger struct {
		path  string
		data  chan interface{}
		close chan bool
	}

	errorLogger struct {
		path   string
		errors chan interface{}
		close  chan bool
	}
)

func NewEventLogger(bufferSize int, path string) Logger {
	return eventLogger{
		path:  path,
		data:  make(chan interface{}, bufferSize),
		close: make(chan bool),
	}
}

func (el eventLogger) Start() {
	var (
		y, m, d = time.Now().Date()
		fName   = fmt.Sprintf("%sawsDog_log_%d.%d.%d", el.path, y, m, d)
	)

	f, err := os.OpenFile(fName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0640)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	logger := log.New(f, "Info: ", log.Ldate|log.Ltime)
	for {
		select {
		case data := <-el.data:
			if ev, ok := data.(Event); ok {
				logger.Print(ev.String())
			}
		case <-el.close:
			return
		}
	}
}

func (el eventLogger) GetRoute() chan<- interface{} {
	return el.data
}

func (el eventLogger) Close() {
	el.close <- true
}

func NewErrorLogger(bufferSize int, path string) Logger {
	return errorLogger{
		path:   path,
		errors: make(chan interface{}, bufferSize),
		close:  make(chan bool),
	}
}

func (el errorLogger) Start() {
	var (
		y, m, d = time.Now().Date()
		fName   = fmt.Sprintf("%sawsDog_errorLog_%d.%d.%d", el.path, y, m, d)
	)

	f, err := os.OpenFile(fName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0640)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	logger := log.New(f, "Error: ", log.Ldate|log.Ltime)
	for {
		select {
		case data := <-el.errors:
			if err, ok := data.(error); ok {
				logger.Print(err.Error())
			}
		case <-el.close:
			return
		}
	}
}

func (el errorLogger) GetRoute() chan<- interface{} {
	return el.errors
}

func (el errorLogger) Close() {
	el.close <- true
}
