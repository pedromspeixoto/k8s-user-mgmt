package logger

import (
	"log"
	"strings"

	"go.uber.org/zap"

	"go.uber.org/fx/fxevent"
)

type StdoutLogger struct {
	level int
}

func NewStdoutLogger(level int) *StdoutLogger {
	return &StdoutLogger{
		level: level,
	}
}

func (logger *StdoutLogger) Trace() string {
	return ""
}

func (logger *StdoutLogger) Debug(a ...interface{}) {
	if logger.level <= LoggingLevelDebug {
		log.Println(a...)
	}
}

func (logger *StdoutLogger) Debugf(format string, a ...interface{}) {
	if logger.level <= LoggingLevelDebug {
		log.Printf(format, a...)
	}
}

func (logger *StdoutLogger) Info(a ...interface{}) {
	if logger.level <= LoggingLevelInfo {
		log.Println(a...)
	}
}

func (logger *StdoutLogger) Infof(format string, a ...interface{}) {
	if logger.level <= LoggingLevelInfo {
		log.Printf(format, a...)
	}
}

func (logger *StdoutLogger) Printf(format string, a ...interface{}) {
	if logger.level <= LoggingLevelInfo {
		logger.Infof(format, a...)
	}
}

func (logger *StdoutLogger) Notice(a ...interface{}) {
	if logger.level <= LoggingLevelInfo {
		log.Println(a...)
	}
}

func (logger *StdoutLogger) Noticef(format string, a ...interface{}) {
	if logger.level <= LoggingLevelInfo {
		log.Printf(format, a...)
	}
}

func (logger *StdoutLogger) Warning(a ...interface{}) {
	if logger.level <= LoggingLevelWarning {
		log.Println(a...)
	}
}

func (logger *StdoutLogger) Warningf(format string, a ...interface{}) {
	if logger.level <= LoggingLevelWarning {
		log.Printf(format, a...)
	}
}

func (logger *StdoutLogger) Error(a ...interface{}) {
	if logger.level <= LoggingLevelError {
		log.Println(a...)
	}
}

func (logger *StdoutLogger) Errorf(format string, a ...interface{}) {
	if logger.level <= LoggingLevelError {
		log.Printf(format, a...)
	}
}

func (logger *StdoutLogger) Fatal(a ...interface{}) {
	log.Fatal(a...)
}

func (logger *StdoutLogger) Fatalf(format string, a ...interface{}) {
	log.Fatalf(format, a...)
}

// LogEvent logs the given event to the provided console logger.
// Taken from https://github.com/uber-go/fx/blob/master/fxevent/console.go#L44
func (logger *StdoutLogger) LogEvent(event fxevent.Event) {
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

func (logger *StdoutLogger) ZapInfo(msg string, fields ...zap.Field) {
	logger.Info(msg)
}
