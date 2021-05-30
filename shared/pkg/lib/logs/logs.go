package logs

import "github.com/op/go-logging"

func MustGetLogger(module string) *logging.Logger {
	logging.MustStringFormatter(`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`)
	return logging.MustGetLogger(module)
}
