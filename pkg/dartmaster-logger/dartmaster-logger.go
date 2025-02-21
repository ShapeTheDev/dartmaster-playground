package dartmasterlogger

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type DartmasterLogger struct {
	log      *log.Logger
	statuses *statuses
}

type statuses struct {
	err string
}

// NewDartmasterLogger returns a new Dartmaster-logger
func NewDartmasterLogger(prefix string) *DartmasterLogger {
	dartmasterLogger := &DartmasterLogger{
		log:      log.New(os.Stdout, prefix, log.LstdFlags),
		statuses: newStatuses(),
	}
	return dartmasterLogger
}

func newStatuses() *statuses {
	statuses := &statuses{
		err: "❌",
	}
	return statuses
}

// Println prints a message
func (dl *DartmasterLogger) Println(v ...any) {
	dl.log.Println(v...)
}

// Printf prints a customized message
func (dl *DartmasterLogger) Printf(format string, v ...any) {
	dl.log.Printf(format, v...)
}

func (dl *DartmasterLogger) PrintlnErr(v ...any) {
	dl.log.Println(append([]any{dl.statuses.err}, v...)...)
}

// PrintfErr prints a customized message with error status
func (dl *DartmasterLogger) PrintfErr(format string, v ...any) {
	dl.log.Printf("%s "+format, append([]any{dl.statuses.err}, v...)...)
}

// LogHttpRequest logs an http request.
func (dl *DartmasterLogger) LogHttpRequest(r *http.Request) {
	clientIP := r.RemoteAddr
	method := r.Method
	uri := r.RequestURI
	dl.Printf("request from IP: %s, method: %s, uri: %s", clientIP, method, uri)
}

// LogAndWriteHttpRequestError logs and writes an http-error.
func (dl *DartmasterLogger) LogAndWriteHttpRequestError(w http.ResponseWriter, status int, err error) {
	dl.Printf(err.Error())
	http.Error(w, err.Error(), status)
}

// PrintRegisteredRouterPaths prints the registered paths of a hand overed router.
func (dl *DartmasterLogger) PrintRegisteredRouterPaths(programName string, path string, router *mux.Router, port string) {
	dl.Printf("%v listening on :%v\n", programName, port)
	dl.Println("registered APIs:")
	err := router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		t, errTemp := route.GetPathTemplate()
		if errTemp != nil {
			return errTemp
		}
		dl.Println(t)
		return nil
	})
	if err != nil {
		dl.PrintfErr("error while logging registered APIs from :%v\n", port)
	}
	if path != "" {
		dl.Printf("%v reachable on: http://localhost:%v/%v", programName, port, path)
	}
	fmt.Println("")
}
