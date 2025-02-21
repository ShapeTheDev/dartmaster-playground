package uricaller

import (
	"bufio"
	"fmt"
	"io"
	"net/http"

	dartmasterlogger "github.com/One-Hundred-Eighty/Circle/pkg/dartmaster-logger"
)

type dartcounterUriCaller struct {
	logger   *dartmasterlogger.DartmasterLogger
	protocol string
	host     string
	port     string
}

type uriBuilder struct{}

type eventStream struct {
	Id    string
	Event string
	Data  []byte // dynamic data (e.g. binary-status-change or binary-output-change)
}

func NewDartcounterUriCaller() *dartcounterUriCaller {
	return &dartcounterUriCaller{
		logger:   dartmasterlogger.NewDartmasterLogger("[uri-caller] "),
		protocol: "http",
		host:     "localhost",
		port:     "8888",
	}
}

// buildURL constructs the complete URL for a given endpoint.
func (duc *dartcounterUriCaller) buildURL(path string) string {
	return fmt.Sprintf("%s://%s:%s%s", duc.protocol, duc.host, duc.port, path)
}

// sse calls the sse uri and returns a channel to get eventStreams.
// path: "/sse"
func (duc *dartcounterUriCaller) SSE() (<-chan eventStream, error) {
	// create a channel to stream sse events
	eventStreamChan := make(chan eventStream, 10)

	// start the sse connection in a new goroutine
	go func() {
		defer close(eventStreamChan) // close the channel when the goroutine ends

		// create request
		req, err := http.NewRequest("GET", duc.buildURL("/sse"), nil)
		if err != nil {
			errMsg := fmt.Errorf("SSE() - error: creating request: %v", err)
			duc.logger.PrintlnErr(errMsg)
			return
		}
		req.Header.Set("Accept", "text/event-stream")

		// call uri
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			errMsg := fmt.Errorf("SSE() - error sending request: %v", err)
			duc.logger.PrintlnErr(errMsg)
			return
		}
		defer resp.Body.Close()

		// ensure the response is for SSE
		if resp.StatusCode != http.StatusOK || resp.Header.Get("Content-Type") != "text/event-stream" {
			errMsg := fmt.Errorf("SSE() - invalid response from /sse endpoint: %v", resp.Status)
			duc.logger.PrintlnErr(errMsg)
			return
		}

		// read the sse stream
		reader := bufio.NewReader(resp.Body)
		var eventStream eventStream
		var eventTracker int

		for {
			// read each line from the stream
			line, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					duc.logger.Println("SSE() - connection closed by server")
					break
				}
				duc.logger.PrintfErr("SSE() - error: reading SSE stream: %v", err)
				break
			}

			if line == "\n" || eventTracker == 3 { // termination criterion of an event stream
				eventStreamChan <- eventStream // send event
				eventTracker = 0
			} else {
				// get id, event, data
				if line[:3] == "id:" {
					eventStream.Id = line[4:(len(line) - 1)] // (len(line)-1) --> cuts the last "\n" from the id
					eventTracker += 1
				}
				if line[:6] == "event:" {
					eventStream.Event = line[7:(len(line) - 1)] // (len(line)-1) --> cuts the last "\n" from the event
					eventTracker += 1
				}
				if line[:5] == "data:" {
					if line != "data: \n" {
						eventStream.Data = []byte(line[6:])
						eventTracker += 1
					}
				}
			}
		}
	}()

	return eventStreamChan, nil
}
