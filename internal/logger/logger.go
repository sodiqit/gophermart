package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	Debugw(msg string, keysAndValues ...interface{})
	Infow(msg string, keysAndValues ...interface{})
	Infoln(args ...interface{})
	Info(args ...interface{})
	Warnw(msg string, keysAndValues ...interface{})
	Errorw(msg string, keysAndValues ...interface{})
	DPanicw(msg string, keysAndValues ...interface{})
	Panicw(msg string, keysAndValues ...interface{})
	Fatalw(msg string, keysAndValues ...interface{})
	With(args ...interface{}) Logger

	Sync() error
}

type zapLoggerWrapper struct {
	*zap.SugaredLogger
}

func (z *zapLoggerWrapper) With(args ...interface{}) Logger {
	return &zapLoggerWrapper{z.SugaredLogger.With(args...)}
}

func New(level string) Logger {
	parsedLevel, err := zap.ParseAtomicLevel(level)

	if err != nil {
		panic(err)
	}

	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.Level = zap.NewAtomicLevelAt(parsedLevel.Level())

	logger := zap.Must(config.Build())
	sugar := logger.Sugar()

	return &zapLoggerWrapper{sugar}
}
