package api

import (
	"fmt"
	"net/http"
	"runtime"
	"strings"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_log.Debugf("%s: %s", r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func errorf(w http.ResponseWriter, status int, msg string, args ...interface{}) {
	pc, _, line, ok := runtime.Caller(1)
	f := runtime.FuncForPC(pc)

	if ok && f != nil {
		parts := strings.Split(f.Name(), ".")
		_log.Errorf(fmt.Sprintf("%s:%d> %s", parts[len(parts)-1], line, msg), args...)
	} else {
		_log.Errorf(msg, args...)
	}

	http.Error(w, fmt.Errorf(msg, args...).Error(), status)
}
