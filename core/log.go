package core

import (
	"fmt"

	"github.com/goinbox/golog"
)

var Logger golog.Logger = new(golog.NoopLogger)
var RequestLogger golog.Logger = new(golog.NoopLogger)

func EmergencyLog(title string, msgs ...interface{}) {
	Logger.Emergency(makeLogMsg(title, msgs))
}

func AlertLog(title string, msgs ...interface{}) {
	Logger.Alert(makeLogMsg(title, msgs))
}

func CriticalLog(title string, msgs ...interface{}) {
	Logger.Critical(makeLogMsg(title, msgs))
}

func ErrorLog(title string, msgs ...interface{}) {
	Logger.Error(makeLogMsg(title, msgs))
}

func WarningLog(title string, msgs ...interface{}) {
	Logger.Warning(makeLogMsg(title, msgs))
}

func NoticeLog(title string, msgs ...interface{}) {
	Logger.Notice(makeLogMsg(title, msgs))
}

func InfoLog(title string, msgs ...interface{}) {
	Logger.Info(makeLogMsg(title, msgs))
}

func DebugLog(title string, msgs ...interface{}) {
	Logger.Debug(makeLogMsg(title, msgs))
}

func makeLogMsg(title string, msgs []interface{}) string {
	return title + "\t" + fmt.Sprint(msgs...)
}
