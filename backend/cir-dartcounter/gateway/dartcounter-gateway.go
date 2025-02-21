package gateway

import (
	"encoding/json"
	"net/http"
	"time"

	dartmasterlogger "github.com/One-Hundred-Eighty/Circle/pkg/dartmaster-logger"
	"github.com/One-Hundred-Eighty/Circle/pkg/sse"
)

type dartcounterGateway struct {
	logger    *dartmasterlogger.DartmasterLogger
	sseServer *sse.SseServer
}

func NewDartcounterGateway(logger *dartmasterlogger.DartmasterLogger) *dartcounterGateway {
	dartcounterGateway := &dartcounterGateway{
		logger:    logger,
		sseServer: sse.NewSseServer("[dartcounter-sse] "),
	}
	dartcounterGateway.startSharingDataViaSse()
	return dartcounterGateway
}

// SSE serves http event-streams to all clients that call this function.
func (g *dartcounterGateway) SSE() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// log the request details
		g.logger.LogHttpRequest(r)

		// get the client's IP address
		ipAddress := r.RemoteAddr

		// serveHTTP via sse-server
		g.sseServer.ServeHTTP(w, r, ipAddress)
	}
}

// TODO
// StartSharingDataViaSse starts sharing data from the binary-admin via the sse-server.
func (g *dartcounterGateway) startSharingDataViaSse() {

	type Auto struct {
		Marke   string `json:"marke"`
		Baujahr int    `json:"baujahr"`
	}

	go func() {
		for {
			// create car
			auto := Auto{
				Marke:   "Volkswagen",
				Baujahr: 2023,
			}

			// convert struct to json
			data, err := json.Marshal(auto)
			if err != nil {
				g.logger.PrintlnErr("JSON Marshal Error:", err)
				continue
			}

			// send event via sse
			g.sseServer.SendEvent("1", "new-auto", data)
			time.Sleep(5 * time.Second)
		}
	}()
}
