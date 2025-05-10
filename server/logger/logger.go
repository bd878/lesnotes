package logger

import (
	"sync/atomic"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	DebugLevel string = "debug"
	InfoLevel string = "info"
	WarnLevel  string = "warn"
	ErrorLevel string = "error"
)

type Config struct {
	LogPath    string
	NodeName   string
	SkipCaller int
	LogLevel   string
}

type Logger struct {
	*zap.SugaredLogger
	conf Config
}

var defaultLogger atomic.Pointer[Logger]

func init() {
	defaultLogger.Store(New(Config{SkipCaller: 1}))
}

func Default() *Logger {
	return defaultLogger.Load()
}

func SetDefault(l *Logger) {
	defaultLogger.Store(l)
}

func New(cfg Config) *Logger {
	var options []zap.Option
	if cfg.SkipCaller > 0 {
		options = append(options, zap.AddCallerSkip(cfg.SkipCaller))
	}

	zapConfig := zap.NewDevelopmentConfig()
	switch cfg.LogLevel {
	case ErrorLevel:
		zapConfig.Level.SetLevel(zapcore.ErrorLevel)
	case DebugLevel:
		zapConfig.Level.SetLevel(zapcore.DebugLevel)
	case InfoLevel:
		zapConfig.Level.SetLevel(zapcore.InfoLevel)
	case WarnLevel:
		zapConfig.Level.SetLevel(zapcore.WarnLevel)
	}

	return &Logger{
		SugaredLogger: zap.Must(zapConfig.Build(options...)).Sugar(),
		conf: cfg,
	}
}

func Error(args ...any) {
	Default().Error(args...)
}

func Errorw(msg string, args ...any) {
	Default().Errorw(msg, args...)
}

func Errorf(msg string, args ...any) {
	Default().Errorf(msg, args...)
}

func Errorln(args ...any) {
	Default().Errorln(args...)
}

func Info(args ...any) {
	Default().Info(args...)
}

func Infow(msg string, args ...any) {
	Default().Infow(msg, args...)
}

func Infoln(args ...any) {
	Default().Infoln(args...)
}

func Infof(msg string, args ...any) {
	Default().Infof(msg, args...)
}

func Warn(args ...any) {
	Default().Warn(args...)
}

func Warnf(msg string, args ...any) {
	Default().Warnf(msg, args...)
}

func Warnln(args ...any) {
	Default().Warnln(args...)
}

func Debugw(msg string, args ...any) {
	Default().Debugw(msg, args...)
}