package utils

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// NewHttpServer returns a new http server.
func NewHttpServer(router *mux.Router, port string) *http.Server {
	httpServer := &http.Server{
		ReadTimeout: 5 * time.Second,
		Handler:     router,
		Addr:        fmt.Sprintf("0.0.0.0:%v", port),
	}
	return httpServer
}
