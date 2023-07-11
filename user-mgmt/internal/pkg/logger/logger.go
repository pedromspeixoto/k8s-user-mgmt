package logger

import (
	"github.com/pedromspeixoto/users-api/internal/config"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

const (
	TypeStdout = "stdout"
	TypeZap    = "zap"
)

const (
	LoggingLevelDebug = 1 << iota
	LoggingLevelInfo
	LoggingLevelWarning
	LoggingLevelError
	LoggingLevelNone
)

func ProvideLogger() fx.Option {
	return fx.Provide(
		NewLoggingClient,
	)
}

type loggerDeps struct {
	fx.In

	Config *config.Config
}

type Logger interface {
	Trace() string

	Debug(a ...interface{})
	Debugf(format string, a ...interface{})
	Info(a ...interface{})
	Infof(format string, a ...interface{})
	Notice(a ...interface{})
	Noticef(format string, a ...interface{})
	Warning(a ...interface{})
	Warningf(format string, a ...interface{})
	Error(a ...interface{})
	Errorf(format string, a ...interface{})
	Fatal(a ...interface{})
	Fatalf(format string, a ...interface{})
	Printf(format string, a ...interface{})
	LogEvent(event fxevent.Event)
	ZapInfo(msg string, fields ...zap.Field)
}

type LoggingClient struct {
	deps       loggerDeps
	loggerType string
	logLevel   int
}

type ClientOptions struct {
	LoggerType string
	LogLevel   int
	ProjectId  string
}

func NewLoggingClient(deps loggerDeps) (*LoggingClient, error) {
	lm := &LoggingClient{
		loggerType: deps.Config.LoggerType,
		logLevel:   deps.Config.LoggerLevel,
		deps:       deps,
	}
	return lm, nil
}

func (lm *LoggingClient) GetLogger() Logger {
	switch lm.loggerType {
	case TypeZap:
		return NewZapLogger(lm.logLevel)
	}
	return NewStdoutLogger(lm.logLevel)
}
