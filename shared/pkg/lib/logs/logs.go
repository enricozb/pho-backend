package logs

import "github.com/op/go-logging"

func MustGetLogger(module string) *logging.Logger {
	f := logging.MustStringFormatter(`%{color}%{time:15:04:05.000} [%{level:.5s}] %{id:03x}%{color:reset} %{message}`)
	logging.SetFormatter(f)
	return logging.MustGetLogger(module)
}
