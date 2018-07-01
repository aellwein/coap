package logging

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestLogging(t *testing.T) {
	logger := LoggerFactory.GetLogger("test")
	logger.Debug("debug")
	logger.DebugF("debug: %d", 42)
	logger.Info("info")
	logger.InfoF("info: %d", 42)
	logger.Warn("warn")
	logger.WarnF("warn: %d", 42)
	logger.Error("error")
	logger.ErrorF("error: %d", 42)
}

func TestPanicLogging(t *testing.T) {
	Convey("Given a logger", t, func() {
		logger := LoggerFactory.GetLogger("test")
		Convey("When logged on panic level", func() {
			Convey("Then panic is expected", func() {
				So(func() { logger.Panic("panic!") }, ShouldPanic)
			})
		})
	})
}

func TestPanicFLogging(t *testing.T) {
	Convey("Given a logger", t, func() {
		logger := LoggerFactory.GetLogger("test")
		Convey("When logged on panic level", func() {
			Convey("Then panic is expected", func() {
				So(func() { logger.PanicF("panic with arg: %d!", 42) }, ShouldPanic)
			})
		})
	})
}
