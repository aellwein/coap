package logging

import "go.uber.org/zap"

var (
	Logger *zap.Logger
	Sugar  *zap.SugaredLogger
)

func InitLogger(development bool) {
	var err error
	if development {
		Logger, err = zap.NewDevelopment()
	} else {
		Logger, err = zap.NewProduction()
	}
	if err != nil {
		panic(err)
	}
	Sugar = Logger.Sugar()
}
