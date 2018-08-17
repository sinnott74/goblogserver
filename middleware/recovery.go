package middleware

// The original work was derived from Goji's middleware, source:
// https://github.com/zenazn/goji/tree/master/web/middleware

import (
	"fmt"
	"net/http"
	"os"
	"reflect"
	"runtime/debug"

	chiMiddleware "github.com/go-chi/chi/middleware"
)

var responseMap = make(map[reflect.Type]int)

// MapErrorResponseStatus configures the a statusCode for a given Error type
// An unrecovered panic will trigger a look to set a http response
func MapErrorResponseStatus(err error, statusCode int) {
	_, ok := err.(error)
	if !ok {
		panic(err)
	}
	errType := reflect.TypeOf(err)
	responseMap[errType] = statusCode
}

func getErrorAndStatus(err interface{}) (int, string) {
	errType := reflect.TypeOf(err)
	statusCode := responseMap[errType]
	// err not mapped, return default 500 Internal Server Error
	if statusCode == 0 {
		return http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)
	}
	return statusCode, err.(error).Error()
}

// Recoverer is a middleware that recovers from panics, logs the panic (and a
// backtrace), and returns a HTTP 500 (Internal Server Error) status if
// possible. Recoverer prints a request ID if one is provided.
//
// Alternatively, look at https://github.com/pressly/lg middleware pkgs.
func Recoverer(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil {

				logEntry := chiMiddleware.GetLogEntry(r)
				if logEntry != nil {
					logEntry.Panic(rvr, debug.Stack())
				} else {
					fmt.Fprintf(os.Stderr, "Panic: %+v\n", rvr)
					debug.PrintStack()
				}
				statusCode, responseText := getErrorAndStatus(rvr)
				http.Error(w, responseText, statusCode)
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
