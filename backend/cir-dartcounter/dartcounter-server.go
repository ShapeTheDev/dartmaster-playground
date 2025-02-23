package cirdartcounter

import (
	"net/http"

	"github.com/One-Hundred-Eighty/Circle/backend/cir-dartcounter/gateway"
	dartmasterlogger "github.com/One-Hundred-Eighty/Circle/pkg/dartmaster-logger"
	"github.com/One-Hundred-Eighty/Circle/utils"
	"github.com/gorilla/mux"
)

func NewServer(logger *dartmasterlogger.DartmasterLogger, port string) *http.Server {
	router := mux.NewRouter()
	dartcounterGateway := gateway.NewDartcounterGateway(logger)

	// initiate dartcounter uris
	router.Path("/dartcounter/sse").HandlerFunc(dartcounterGateway.SSE()).Methods(http.MethodGet)

	// initiate http server
	httpServer := utils.NewHttpServer(router, port)

	// print the registered routes for debugging-purposes
	logger.PrintRegisteredRouterPaths("dartcounter", "", router, port)

	return httpServer
}
