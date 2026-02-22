package logger

import (
	"context"

	"github.com/sirupsen/logrus"
)

type Logger interface {
	Trace(args ...interface{})

	Debug(args ...interface{})

	Info(args ...interface{})

	Warn(args ...interface{})

	Error(args ...interface{})

	Fatal(args ...interface{})

	Panic(args ...interface{})

	Tracef(format string, args ...interface{})

	Debugf(format string, args ...interface{})

	Infof(format string, args ...interface{})

	Printf(format string, args ...interface{})

	Warnf(format string, args ...interface{})

	Errorf(format string, args ...interface{})

	Fatalf(format string, args ...interface{})

	Panicf(format string, args ...interface{})

	WithField(key string, value interface{}) Logger

	WithFields(fields Fields) Logger
}

func Init(_ context.Context) {
	logrus.SetReportCaller(true)
	logFormat := &logrus.TextFormatter{
		ForceColors:     false,
		DisableColors:   false,
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	}
	logrus.SetFormatter(logFormat)
}

func GetLogger(ctx context.Context) Logger {
	logger := ctx.Value(loggerCtxKey)
	if logger == nil {
		return loggerLogrus{
			Entry: logrus.NewEntry(logrus.StandardLogger()),
		}
	}

	return logger.(Logger)
}

func WithLogger(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, loggerCtxKey, logger)
}
