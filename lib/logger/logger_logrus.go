package logger

import "github.com/sirupsen/logrus"

type Fields map[string]interface{}

type ctxKey string

const loggerCtxKey ctxKey = "logger"

type loggerLogrus struct {
	*logrus.Entry
}

func (l loggerLogrus) WithField(key string, value interface{}) Logger {
	return loggerLogrus{
		Entry: l.Entry.WithField(key, value),
	}
}

func (l loggerLogrus) WithFields(fields Fields) Logger {
	return loggerLogrus{
		Entry: l.Entry.WithFields(logrus.Fields(fields)),
	}
}
