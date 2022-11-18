package logging

import (
	"fmt"
	"github.com/ohler55/ojg/pretty"
	"go.uber.org/zap"
)

var Logger *zap.Logger

func SetZapLogger(zapLogger *zap.Logger) {
	Logger = zapLogger
}

var ShouldLog = false
var ShouldLogIncoming = false

func LogIncoming[T any](incoming T) {
	if ShouldLogIncoming {
		s := fmt.Sprintf("request value: %v", pretty.JSON(incoming))
		Logger.Info(s)
	}
}

func LogError[T any](s string) {
	if ShouldLog {
		Logger.Error(s)
	}
}

func LogInfo[T any](s string) {
	if ShouldLog {
		Logger.Info(s)
	}
}

func LogWarn[T any](s string) {
	if ShouldLog {
		Logger.Warn(s)
	}
}
