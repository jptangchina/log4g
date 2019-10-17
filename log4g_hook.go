package log4g

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Log4gHook struct {
	Writer *lumberjack.Logger
	Level  []logrus.Level
}

func (l *Log4gHook) Levels() []logrus.Level {
	return l.Level
}

func (l *Log4gHook) Fire(entry *logrus.Entry) error {
	serialized, err := entry.Logger.Formatter.Format(entry)
	if err != nil {
		return err
	}
	_, writeErr := l.Writer.Write(serialized)
	if nil != writeErr {
		return writeErr
	}
	return nil
}
