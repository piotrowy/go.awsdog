package main

import (
	"fmt"
	"os"
	"regexp"
	"testing"
	"time"
)

const buffSize = 10

func TestNewEventLogger(t *testing.T) {
	l := NewEventLogger(buffSize, DefaultLogPath)

	if l == nil {
		t.Error("wrong initialization of logger")
	}
}

func TestEventLogger_Start(t *testing.T) {

}

func TestEventLogger_Log(t *testing.T) {

}

func TestEventLogger_Close(t *testing.T) {

}

func TestNewErrorLogger(t *testing.T) {
	l := NewErrorLogger(buffSize, DefaultLogPath)

	if l == nil {
		t.Error("wrong initialization of errorLogger")
	}
}

func TestErrorLogger_Start(t *testing.T) {
	//var (
	//	y, m, d    = time.Now().Date()
	//	fName      = fmt.Sprintf("%s.%d.%d.%d", DefaultLogPath, y, m, d)
	//	l          = NewErrorLogger(buffSize, DefaultLogPath)
	//	exampleErr = fmt.Errorf("whatever error")
	//	data       = make([]byte, 100)
	//)
	//
	//go l.Start()
	//
	//if _, err := os.Stat(fName); os.IsNotExist(err) {
	//	t.Error("log file is not created")
	//}
	//
	//l.Log(exampleErr)
	//f, _ := os.Open(fName)
	//count, _ := f.Read(data)
	//if ok, err := regexp.Match(exampleErr.Error(), data[:count]); !ok && err != nil {
	//	t.Error("wrong error logged")
	//}
	//l.Close()
}

func TestErrorLogger_Log(t *testing.T) {

}

func TestErrorLogger_Close(t *testing.T) {
	//var (
	//	l          = NewErrorLogger(buffSize, DefaultLogPath)
	//	exampleErr = fmt.Errorf("whatever error")
	//)
	//
	//go l.Start()
	//l.Close()
	//
	//select {
	//case l.Log(exampleErr):
	//	t.Error("channel not closed, goroutines not ended")
	//default:
	//}
}
