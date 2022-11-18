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

func LogIncoming(incoming interface{}) {
	if ShouldLogIncoming {
		s := fmt.Sprintf("request value: %v", pretty.JSON(incoming))
		Logger.Info(s)
	}
}

func LogError(s string) {
	if ShouldLog {
		Logger.Error(s)
	}
}

func LogInfo(s string) {
	if ShouldLog {
		Logger.Info(s)
	}
}

func LogWarn(s string) {
	if ShouldLog {
		Logger.Warn(s)
	}
}
