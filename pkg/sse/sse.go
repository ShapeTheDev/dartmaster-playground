package sse

import (
	"bytes"
	"fmt"
	"net/http"

	dartmasterlogger "github.com/One-Hundred-Eighty/Circle/pkg/dartmaster-logger"
	subscriptionhandler "github.com/One-Hundred-Eighty/Circle/pkg/subscription-handler"
)

type WriteListener func(data []byte)

type SseServer struct {
	logger              *dartmasterlogger.DartmasterLogger
	subscriptionHandler *subscriptionhandler.SubscriptionHandler[eventStream]
	writeListener       WriteListener
}

type eventStream struct {
	id        []byte
	eventType []byte
	data      []byte
}

// NewSseServer creates a new SSE-server object and start the sseServer.
func NewSseServer(loggerPrefix string) *SseServer {
	sseServer := &SseServer{
		logger:              dartmasterlogger.NewDartmasterLogger(loggerPrefix),
		subscriptionHandler: subscriptionhandler.NewSubscriptionHandler[eventStream](),
	}
	return sseServer
}

// SetWriteListener sets a listener for writing messages.
func (sse *SseServer) SetWriteListener(listener WriteListener) {
	sse.writeListener = listener
}

// SendEvent sends the given id, eventTyoe and data to all connected SSE clients.
func (sse *SseServer) SendEvent(id string, eventType string, data []byte) {
	eventStream := eventStream{
		id:        []byte(id),
		eventType: []byte(eventType),
		data:      data,
	}
	sse.subscriptionHandler.Publish(eventStream)
}

// ServeHTTP serves the SSE page to forward events to a client. Each connected client starts this method in its own request thread.
// This method is running for each client as long as the connection is not closed or the SSE server is terminated.
func (sse *SseServer) ServeHTTP(rw http.ResponseWriter, req *http.Request, subscriber string) {
	fail := func(msg error) {
		sse.logger.LogAndWriteHttpRequestError(rw, http.StatusInternalServerError, msg)
	}
	// make sure that the writer supports flushing
	flusher, ok := rw.(http.Flusher)
	if !ok {
		errMsg := fmt.Errorf("sse streaming not supported")
		fail(errMsg)
		return
	}

	// set response-writer-header
	rw.Header().Set("Content-Type", "text/event-stream")
	rw.Header().Set("Cache-Control", "no-cache")
	rw.Header().Set("Connection", "keep-alive")
	rw.Header().Set("Access-Control-Allow-Origin", "*")

	// register new sse client
	eventChan := sse.subscribe(subscriber)
	defer sse.unsubscribe(eventChan, subscriber)

	// listen to connection close and un-register eventStreamChan
	ctxDone := req.Context().Done()
	for {
		select {
		case <-ctxDone:
			// connection was closed, aborting
			sse.logger.Println("sse client: connection closed")
			return
		case event := <-eventChan:
			dataLines := bytes.Split(event.data, []byte("\n"))
			// received event from server publish - forward it to the connected client
			for i, dataLine := range dataLines {
				var message string
				eof := ""
				if i >= len(dataLines)-1 {
					eof = "\n" // adds a newline for the last dataLine --> termination criterion to end the event
				}
				message += fmt.Sprintf("id: %s\n", event.id)
				message += fmt.Sprintf("event: %s\n", event.eventType)
				message += fmt.Sprintf("data: %s\n", dataLine)

				totalMessage := message + eof

				n, err := fmt.Fprintf(rw, "%s", totalMessage)
				if err != nil {
					errMsg := fmt.Errorf("sse client: event could not be sent: %v", err)
					fail(errMsg)
					return
				}

				if sse.writeListener != nil {
					sse.writeListener([]byte(totalMessage[:n]))
				}
			}
			// make sure it is actually sent to the client in time and not buffered
			flusher.Flush()
		}
	}
}

// subscribe subscribes on the sseServer.
func (sse *SseServer) subscribe(subscriber string) <-chan eventStream {
	clientChan := sse.subscriptionHandler.Subscribe()
	sse.logger.Printf("client added. client IP: %s. %d registered clients", subscriber, sse.subscriptionHandler.Subscriptions())
	return clientChan
}

// unsubscribe unsubscribes from the sseServer.
func (sse *SseServer) unsubscribe(logChan <-chan eventStream, subscriber string) {
	sse.subscriptionHandler.Unsubscribe(logChan)
	sse.logger.Printf("client removed. client IP: %s. %d registered clients", subscriber, sse.subscriptionHandler.Subscriptions())
}
