package logger

import (
	"fmt"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/rs/zerolog"
)

func init() {
	zerolog.TimestampFieldName = "@timestamp"
}

var (
	defaultLogger = New()
)

func GetLogger() *Logger {
	return defaultLogger
}

type Logger struct {
	zl zerolog.Logger
}

func New() *Logger {
	return newLogger(
		zerolog.New(os.Stdout).
			Level(zerolog.InfoLevel).
			With().
			Timestamp().
			Logger(),
	)
}

func ParseLogLevel(str string) {
	level, err := zerolog.ParseLevel(str)
	if err != nil {
		defaultLogger.WithError(err).Error("logger.ParseLogLevel")
		return
	}

	defaultLogger.zl.Level(level)
	defaultLogger.Warn("log level: ", level.String())
}

func newLogger(zl zerolog.Logger) *Logger {
	return &Logger{
		zl: zl,
	}
}

func (l *Logger) Level(level Level) {
	l.zl.Level(level)
}

func (l *Logger) Spew(args ...any) {
	l.zl.WithLevel(zerolog.NoLevel).Msg(spew.Sdump(args...))
}

func (l *Logger) Debug(args ...any) {
	l.zl.WithLevel(zerolog.DebugLevel).Msg(fmt.Sprint(args...))
}

func (l *Logger) Info(args ...any) {
	l.zl.WithLevel(zerolog.InfoLevel).Msg(fmt.Sprint(args...))
}

func (l *Logger) Warn(args ...any) {
	l.zl.WithLevel(zerolog.WarnLevel).Msg(fmt.Sprint(args...))
}

func (l *Logger) Warning(args ...any) {
	l.Warn(args...)
}

func (l *Logger) Error(args ...any) {
	l.zl.WithLevel(zerolog.ErrorLevel).Msg(fmt.Sprint(args...))
}

func (l *Logger) Fatal(args ...any) {
	l.zl.WithLevel(zerolog.FatalLevel).Msg(fmt.Sprint(args...))
	os.Exit(1)
}

func (l *Logger) Panic(args ...any) {
	msg := fmt.Sprint(args...)
	l.zl.WithLevel(zerolog.PanicLevel).Msg(msg)
	panic(msg)
}

func (l *Logger) Debugf(format string, args ...any) {
	l.zl.WithLevel(zerolog.DebugLevel).Msgf(format, args...)
}

func (l *Logger) Infof(format string, args ...any) {
	l.zl.WithLevel(zerolog.InfoLevel).Msgf(format, args...)
}

func (l *Logger) Warnf(format string, args ...any) {
	l.zl.WithLevel(zerolog.WarnLevel).Msgf(format, args...)
}

func (l *Logger) Warningf(format string, args ...any) {
	l.Warnf(format, args...)
}

func (l *Logger) Errorf(format string, args ...any) {
	l.zl.WithLevel(zerolog.ErrorLevel).Msgf(format, args...)
}

func (l *Logger) Fatalf(format string, args ...any) {
	l.zl.WithLevel(zerolog.FatalLevel).Msgf(format, args...)
	os.Exit(1)
}

func (l *Logger) Panicf(format string, args ...any) {
	l.zl.WithLevel(zerolog.PanicLevel).Msgf(format, args...)
}

func (l *Logger) Msgf(format string, args ...any) {
	l.zl.WithLevel(zerolog.InfoLevel).Msgf(format, args...)
}

func (l *Logger) WithLevel(level Level) *Logger {
	return newLogger(l.zl.Level(level))
}

func (l *Logger) WithError(err error) *Logger {
	return newLogger(l.zl.With().Err(err).Logger())
}
