package logger

import (
	"go.uber.org/zap"
	"http2ws/conf"
	"log"
)

var logger *zap.Logger

func Init() {
	var err error
	loggerCfg := zap.NewProductionConfig()
	loggerCfg.OutputPaths = []string{conf.LogFile}
	loggerCfg.ErrorOutputPaths = []string{conf.LogFile}
	logger, err = loggerCfg.Build()
	zap.NewProduction()
	if err != nil {
		log.Fatal(err)
		return
	}
}

func Debug(v ...interface{}) {
	logger.Sugar().Debug(v...)
}
func Debugf(f string, v ...interface{}) {
	logger.Sugar().Debugf(f, v...)

}
func Info(v ...interface{}) {
	logger.Sugar().Info(v...)

}
func Infof(f string, v ...interface{}) {
	logger.Sugar().Infof(f, v...)

}
func Warn(v ...interface{}) {
	logger.Sugar().Warn(v...)

}
func Warnf(f string, v ...interface{}) {
	logger.Sugar().Warnf(f, v...)

}
func Error(v ...interface{}) {
	logger.Sugar().Error(v...)

}
func Errorf(f string, v ...interface{}) {
	logger.Sugar().Errorf(f, v...)

}
func Fatal(v ...interface{}) {
	logger.Sugar().Fatal(v...)

}
func Fatalf(f string, v ...interface{}) {
	logger.Sugar().Fatalf(f, v...)

}
func Panic(v ...interface{}) {
	logger.Sugar().Panic(v...)

}
func Panicf(f string, v ...interface{}) {
	logger.Sugar().Panicf(f, v...)
}
