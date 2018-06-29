package logging

import (
	"github.com/aellwein/slf4go"
	"go.uber.org/zap"
)

type ZapLoggerAdaptor struct {
	slf4go.LoggerAdaptor
	logger *zap.SugaredLogger
}

func newZapLogger(name string, logger *zap.Logger) *ZapLoggerAdaptor {
	result := new(ZapLoggerAdaptor)
	result.SetName(name)
	result.logger = logger.Sugar()
	result.logger.Named(name)
	return result
}

func (z *ZapLoggerAdaptor) SetLevel(l slf4go.LEVEL) {
	z.LoggerAdaptor.SetLevel(l)
}

func (z *ZapLoggerAdaptor) Trace(args ...interface{}) {
	// trace is no-op in zap
}
func (z *ZapLoggerAdaptor) TraceF(format string, args ...interface{}) {
	// trace is no-op in zap
}
func (z *ZapLoggerAdaptor) Debug(args ...interface{}) {
	z.logger.Debug(args...)
}
func (z *ZapLoggerAdaptor) DebugF(format string, args ...interface{}) {
	z.logger.Debugf(format, args...)
}
func (z *ZapLoggerAdaptor) Info(args ...interface{}) {
	z.logger.Info(args...)
}
func (z *ZapLoggerAdaptor) InfoF(format string, args ...interface{}) {
	z.logger.Infof(format, args...)
}
func (z *ZapLoggerAdaptor) Warn(args ...interface{}) {
	z.logger.Warn(args...)
}
func (z *ZapLoggerAdaptor) WarnF(format string, args ...interface{}) {
	z.logger.Warnf(format, args...)
}
func (z *ZapLoggerAdaptor) Error(args ...interface{}) {
	z.logger.Error(args...)
}
func (z *ZapLoggerAdaptor) ErrorF(format string, args ...interface{}) {
	z.logger.Errorf(format, args...)
}
func (z *ZapLoggerAdaptor) Fatal(args ...interface{}) {
	z.logger.Fatal(args...)
}
func (z *ZapLoggerAdaptor) FatalF(format string, args ...interface{}) {
	z.logger.Fatalf(format, args...)
}
func (z *ZapLoggerAdaptor) Panic(args ...interface{}) {
	z.logger.Panic(args...)
}
func (z *ZapLoggerAdaptor) PanicF(format string, args ...interface{}) {
	z.logger.Panicf(format, args...)
}

type ZapLoggerFactory struct {
	logger *zap.Logger
}

func (factory ZapLoggerFactory) GetLogger(name string) slf4go.Logger {
	return newZapLogger(name, factory.logger)
}

func NewZapLoggerFactory(logger *zap.Logger) slf4go.LoggerFactory {
	return ZapLoggerFactory{logger: logger}
}
