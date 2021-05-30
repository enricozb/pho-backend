package api

import (
	"fmt"
	"net/http"
	"runtime"
	"strings"
)

func errorf(res http.ResponseWriter, status int, msg string, args ...interface{}) {
	pc, _, line, ok := runtime.Caller(1)
	f := runtime.FuncForPC(pc)

	if ok && f != nil {
		parts := strings.Split(f.Name(), ".")
		_log.Errorf(fmt.Sprintf("%s:%d> %s", parts[len(parts)-1], line, msg), args...)
	} else {
		_log.Errorf(msg, args...)
	}

	http.Error(res, fmt.Errorf(msg, args...).Error(), status)
}
