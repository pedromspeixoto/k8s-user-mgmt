package logger

import (
	"fmt"
	"strings"

	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

type ZapLogger struct {
	Zap   *zap.Logger
	Level int
}

func NewZapLogger(level int) *ZapLogger {
	zap, err := zap.NewProduction()
	if err != nil {
		return nil
	}
	return &ZapLogger{
		Zap:   zap,
		Level: level,
	}
}

func (logger *ZapLogger) Trace() string {
	return ""
}

func (logger *ZapLogger) Debug(a ...interface{}) {
	if logger.Level <= LoggingLevelDebug {
		logger.Zap.Debug(fmt.Sprint(a...))
	}
}

func (logger *ZapLogger) Debugf(format string, a ...interface{}) {
	if logger.Level <= LoggingLevelDebug {
		logger.Zap.Debug(fmt.Sprintf(format, a...))
	}
}

func (logger *ZapLogger) Info(a ...interface{}) {
	if logger.Level <= LoggingLevelInfo {
		logger.Zap.Info(fmt.Sprint(a...))
	}
}

func (logger *ZapLogger) Infof(format string, a ...interface{}) {
	if logger.Level <= LoggingLevelInfo {
		logger.Zap.Info(fmt.Sprintf(format, a...))
	}
}

func (logger *ZapLogger) Printf(format string, a ...interface{}) {
	if logger.Level <= LoggingLevelInfo {
		logger.Zap.Info(fmt.Sprintf(format, a...))
	}
}

func (logger *ZapLogger) Notice(a ...interface{}) {
	if logger.Level <= LoggingLevelInfo {
		logger.Zap.Info(fmt.Sprint(a...))
	}
}

func (logger *ZapLogger) Noticef(format string, a ...interface{}) {
	if logger.Level <= LoggingLevelInfo {
		logger.Zap.Info(fmt.Sprintf(format, a...))
	}
}

func (logger *ZapLogger) Warning(a ...interface{}) {
	if logger.Level <= LoggingLevelWarning {
		logger.Zap.Warn(fmt.Sprint(a...))
	}
}

func (logger *ZapLogger) Warningf(format string, a ...interface{}) {
	if logger.Level <= LoggingLevelWarning {
		logger.Zap.Warn(fmt.Sprintf(format, a...))
	}
}

func (logger *ZapLogger) Error(a ...interface{}) {
	if logger.Level <= LoggingLevelError {
		logger.Zap.Error(fmt.Sprint(a...))
	}
}

func (logger *ZapLogger) Errorf(format string, a ...interface{}) {
	if logger.Level <= LoggingLevelError {
		logger.Zap.Error(fmt.Sprintf(format, a...))
	}
}

func (logger *ZapLogger) Fatal(a ...interface{}) {
	logger.Zap.Fatal(fmt.Sprint(a...))
}

func (logger *ZapLogger) Fatalf(format string, a ...interface{}) {
	logger.Zap.Fatal(fmt.Sprintf(format, a...))
}

// LogEvent logs the given event to the provided console logger.
// Taken from https://github.com/uber-go/fx/blob/master/fxevent/console.go#L44
func (logger *ZapLogger) LogEvent(event fxevent.Event) {
	switch e := event.(type) {
	case *fxevent.OnStartExecuting:
		logger.Infof("HOOK OnStart\t\t%s executing (caller: %s)", e.FunctionName, e.CallerName)
	case *fxevent.OnStartExecuted:
		if e.Err != nil {
			logger.Infof("HOOK OnStart\t\t%s called by %s failed in %s: %+v", e.FunctionName, e.CallerName, e.Runtime, e.Err)
		} else {
			logger.Infof("HOOK OnStart\t\t%s called by %s ran successfully in %s", e.FunctionName, e.CallerName, e.Runtime)
		}
	case *fxevent.OnStopExecuting:
		logger.Infof("HOOK OnStop\t\t%s executing (caller: %s)", e.FunctionName, e.CallerName)
	case *fxevent.OnStopExecuted:
		if e.Err != nil {
			logger.Infof("HOOK OnStop\t\t%s called by %s failed in %s: %+v", e.FunctionName, e.CallerName, e.Runtime, e.Err)
		} else {
			logger.Infof("HOOK OnStop\t\t%s called by %s ran successfully in %s", e.FunctionName, e.CallerName, e.Runtime)
		}
	case *fxevent.Supplied:
		if e.Err != nil {
			logger.Infof("ERROR\tFailed to supply %v: %+v", e.TypeName, e.Err)
		} else if e.ModuleName != "" {
			logger.Infof("SUPPLY\t%v from module %q", e.TypeName, e.ModuleName)
		} else {
			logger.Infof("SUPPLY\t%v", e.TypeName)
		}
	case *fxevent.Provided:
		for _, rtype := range e.OutputTypeNames {
			if e.ModuleName != "" {
				logger.Infof("PROVIDE\t%v <= %v from module %q", rtype, e.ConstructorName, e.ModuleName)
			} else {
				logger.Infof("PROVIDE\t%v <= %v", rtype, e.ConstructorName)
			}
		}
		if e.Err != nil {
			logger.Infof("Error after options were applied: %+v", e.Err)
		}

	case *fxevent.Replaced:
		for _, rtype := range e.OutputTypeNames {
			if e.ModuleName != "" {
				logger.Infof("REPLACE\t%v from module %q", rtype, e.ModuleName)
			} else {
				logger.Infof("REPLACE\t%v", rtype)
			}
		}
		if e.Err != nil {
			logger.Infof("ERROR\tFailed to replace: %+v", e.Err)
		}
	case *fxevent.Decorated:
		for _, rtype := range e.OutputTypeNames {
			if e.ModuleName != "" {
				logger.Infof("DECORATE\t%v <= %v from module %q", rtype, e.DecoratorName, e.ModuleName)
			} else {
				logger.Infof("DECORATE\t%v <= %v", rtype, e.DecoratorName)
			}
		}
		if e.Err != nil {
			logger.Infof("Error after options were applied: %+v", e.Err)
		}
	case *fxevent.Invoking:
		if e.ModuleName != "" {
			logger.Infof("INVOKE\t\t%s from module %q", e.FunctionName, e.ModuleName)
		} else {
			logger.Infof("INVOKE\t\t%s", e.FunctionName)
		}
	case *fxevent.Invoked:
		if e.Err != nil {
			logger.Infof("ERROR\t\tfx.Invoke(%v) called from:\n%+vFailed: %+v", e.FunctionName, e.Trace, e.Err)
		}
	case *fxevent.Stopping:
		logger.Infof("%v", strings.ToUpper(e.Signal.String()))
	case *fxevent.Stopped:
		if e.Err != nil {
			logger.Infof("ERROR\t\tFailed to stop cleanly: %+v", e.Err)
		}
	case *fxevent.RollingBack:
		logger.Infof("ERROR\t\tStart failed, rolling back: %+v", e.StartErr)
	case *fxevent.RolledBack:
		if e.Err != nil {
			logger.Infof("ERROR\t\tCouldn't roll back cleanly: %+v", e.Err)
		}
	case *fxevent.Started:
		if e.Err != nil {
			logger.Infof("ERROR\t\tFailed to start: %+v", e.Err)
		} else {
			logger.Infof("RUNNING")
		}
	case *fxevent.LoggerInitialized:
		if e.Err != nil {
			logger.Infof("ERROR\t\tFailed to initialize custom logger: %+v", e.Err)
		} else {
			logger.Infof("LOGGER\tInitialized custom logger from %v", e.ConstructorName)
		}
	}
}

func (logger *ZapLogger) ZapInfo(msg string, fields ...zap.Field) {
	logger.Zap.Info(msg, fields...)
}
