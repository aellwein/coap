package logging

import (
	"github.com/aellwein/slf4go"
	"go.uber.org/zap"
)

var LoggerFactory slf4go.LoggerFactory

// log4go facade initialization
func init() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	LoggerFactory = NewZapLoggerFactory(logger)
	slf4go.SetLoggerFactory(LoggerFactory)
}
