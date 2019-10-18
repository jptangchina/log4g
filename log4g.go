package log4g

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
)

var loggers map[logrus.Level]*logrus.Logger

func Trace(args ...interface{}) {
	loggers[logrus.TraceLevel].Trace(args)
}

func Tracef(format string, args ...interface{}) {
	loggers[logrus.TraceLevel].Tracef(format, args)
}

func Traceln(args ...interface{}) {
	loggers[logrus.TraceLevel].Traceln(args)
}

func Debug(args ...interface{}) {
	loggers[logrus.DebugLevel].Debug(args)
}

func Degbugf(format string, args ...interface{}) {
	loggers[logrus.DebugLevel].Debugf(format, args)
}

func Debugln(args ...interface{}) {
	loggers[logrus.DebugLevel].Debugln(args)
}

func Info(args ...interface{}) {
	loggers[logrus.InfoLevel].Info(args)
}

func Infof(format string, args ...interface{}) {
	loggers[logrus.InfoLevel].Infof(format, args)
}

func Infoln(args ...interface{}) {
	loggers[logrus.InfoLevel].Infoln(args)
}

func Warn(args ...interface{}) {
	loggers[logrus.WarnLevel].Warn(args)
}

func Warnf(format string, args ...interface{}) {
	loggers[logrus.WarnLevel].Warnf(format, args)
}

func Warnln(args ...interface{}) {
	loggers[logrus.WarnLevel].Warnln(args)
}

func Error(args ...interface{}) {
	loggers[logrus.ErrorLevel].Error(args)
}

func Errorf(format string, args ...interface{}) {
	loggers[logrus.ErrorLevel].Errorf(format, args)
}

func Errorln(args ...interface{}) {
	loggers[logrus.ErrorLevel].Errorln(args)
}

func Fatal(args ...interface{}) {
	loggers[logrus.FatalLevel].Fatal(args)
}

func Fatalf(format string, args ...interface{}) {
	loggers[logrus.FatalLevel].Fatalf(format, args)
}

func Fatalln(args ...interface{}) {
	loggers[logrus.FatalLevel].Fatalln(args)
}

func Panic(args ...interface{}) {
	loggers[logrus.PanicLevel].Panic(args)
}

func Panicf(format string, args ...interface{}) {
	loggers[logrus.PanicLevel].Panicf(format, args)
}

func Panicln(args ...interface{}) {
	loggers[logrus.PanicLevel].Panicln(args)
}

func init() {
	loggers = make(map[logrus.Level]*logrus.Logger, len(logrus.AllLevels))
	initLogrusWithDefaultLogger()
	appenders, err := InitAppenders()
	if err != nil {
		logrus.Errorf("Can not decode config, caused by: %v", err)
		return
	}
	appenderMap := make(map[logrus.Level]*lumberjack.Logger, len(logrus.AllLevels))
	hooks := make(map[*lumberjack.Logger][]logrus.Level)
	for _, appender := range appenders.Appender {
		var (
			minLevel, maxLevel logrus.Level
			err                error
		)
		if minLevel, err = logrus.ParseLevel(appender.MaxLevel); err != nil {
			logrus.Fatalf("Unable to parse levels,use min instead. caused by: %v", err)
			minLevel = 0
		}
		if maxLevel, err = logrus.ParseLevel(appender.MinLevel); err != nil {
			logrus.Fatalf("Unable to parse levels, use max instead. caused by: %v", err)
			maxLevel = logrus.Level(uint(len(logrus.AllLevels) - 1))
		}
		if minLevel > maxLevel {
			minLevel, maxLevel = maxLevel, minLevel
		}
		logger := newLogger(&appender)
		for level := minLevel; level <= maxLevel; level++ {
			// 如果非空，说明同一个日志级别有不同的appender与之对应，多于的appender采用hook方式输出
			if appenderMap[level] != nil {
				hooks[logger] = append(hooks[logger], level)
				continue
			}
			appenderMap[level] = logger
		}
	}
	for level, output := range appenderMap {
		logger := logrus.New()
		logger.SetOutput(output)
		loggers[level] = logger
	}

	for logger, levels := range hooks {
		hook := newLog4gHook(logger, levels)
		for _, level := range levels {
			loggers[level].AddHook(hook)
		}
	}
}

func initLogrusWithDefaultLogger() {
	if nil == loggers {
		loggers = make(map[logrus.Level]*logrus.Logger, len(logrus.AllLevels))
	}
	// 确保所有级别日志都有对应输出
	for _, level := range logrus.AllLevels {
		if loggers[level] == nil {
			loggers[level] = logrus.New()
			loggers[level].SetOutput(getDefaultOutput(level))
		}
	}
}

func newLog4gHook(logger *lumberjack.Logger, level []logrus.Level) *Log4gHook {
	return &Log4gHook{
		Writer: logger,
		Level:  level,
	}
}

func newLogger(appender *Appender) *lumberjack.Logger {
	if appender.Filename == "" {
		return nil
	}
	return &lumberjack.Logger{
		Filename:   appender.Filename,
		MaxSize:    appender.MaxSize,
		MaxAge:     appender.MaxAge,
		MaxBackups: appender.MaxBackups,
		LocalTime:  appender.LocalTime,
		Compress:   appender.Compress,
	}
}

func getDefaultOutput(level logrus.Level) io.Writer {
	switch level {
	case logrus.PanicLevel, logrus.FatalLevel, logrus.ErrorLevel:
		return os.Stderr
	default:
		return os.Stdout
	}
}

func GetLogger(level logrus.Level) *logrus.Logger {
	return loggers[level]
}
