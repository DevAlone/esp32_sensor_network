package logger

import (
	"github.com/go-errors/errors"
	logging "github.com/op/go-logging"
)

var Log = logging.MustGetLogger("esp32_sensor_network")
var LogFormat = logging.MustStringFormatter(
	`%{color}%{module} %{pid} %{level:.5s} %{time:15:04:05.000} %{shortfile} %{shortfunc} ▶ %{id:03x}%{color:reset} %{message}`,
)

func LogError(err error, message string) {
	if err == nil {
		return
	}

	if len(message) > 0 {
		Log.Errorf("%v:%v", message, err)
	} else {
		Log.Error(err)
	}
	// if it has stack
	if errWithStack, ok := err.(*errors.Error); ok {
		Log.Error(errWithStack.ErrorStack())
	}
}

func init() {
}
