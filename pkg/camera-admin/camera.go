package cameraadmin

import (
	"context"
	"errors"
	"fmt"
	"time"

	camerasubscriptionhandler "github.com/One-Hundred-Eighty/Circle/pkg/camera-admin/camera-subscription-handler"
	"github.com/One-Hundred-Eighty/Circle/pkg/camera-admin/device"
	"github.com/One-Hundred-Eighty/Circle/pkg/camera-admin/v4l2"
	dartmasterlogger "github.com/One-Hundred-Eighty/Circle/pkg/dartmaster-logger"
)

type cameraAdmin struct {
	logger  *dartmasterlogger.DartmasterLogger
	cameras []*camera
}

type camera struct {
	subscriptionHandler *camerasubscriptionhandler.CameraSubscriptionHandler[[]byte]
	stopPublisherCh     chan struct{}
	outputCh            <-chan []byte
	id                  int
	devicePath          string
	device              *device.Device
}

func NewCameraAdmin() *cameraAdmin {
	cameraAdmin := &cameraAdmin{
		logger: dartmasterlogger.NewDartmasterLogger("[camera-admin] "),
		cameras: []*camera{
			{subscriptionHandler: camerasubscriptionhandler.NewCameraSubscriptionHandler[[]byte](),
				id:         1,
				devicePath: "/dev/video0"},
			{subscriptionHandler: camerasubscriptionhandler.NewCameraSubscriptionHandler[[]byte](),
				id:         2,
				devicePath: "/dev/video2"},
			{subscriptionHandler: camerasubscriptionhandler.NewCameraSubscriptionHandler[[]byte](),
				id:         3,
				devicePath: "/dev/video4"},
		},
	}
	return cameraAdmin
}

// Start starts all cameras
//
// width: resolution-width in pixels |
// height: resolution-height in pixels
func (ca *cameraAdmin) Start(width, height int) error {
	ca.logger.Println("start cameras")
	widthUint32 := uint32(width)
	heightUint32 := uint32(height)
	var err error

	for i, c := range ca.cameras {
		// open camera
		ca.cameras[i].device, err = device.Open(
			c.devicePath,
			device.WithPixFormat(v4l2.PixFormat{PixelFormat: v4l2.PixelFmtMJPEG, Width: widthUint32, Height: heightUint32}),
		)
		if err != nil {
			errMsg := fmt.Sprintf("Start() - error: opening camera (camera-id: %d, device-path: %v): %v", c.id, c.devicePath, err)
			return errors.New(errMsg)
		}

		// start camera
		err := ca.cameras[i].device.Start(context.TODO())
		if err != nil {
			errMsg := fmt.Sprintf("Start() - error: starting camera (camera-id: %d)", c.id)
			return errors.New(errMsg)
		}

		// create stop channel for the frame publisher
		stopCh := make(chan struct{})
		ca.cameras[i].stopPublisherCh = stopCh

		// grep cameraOutput channel
		ca.cameras[i].outputCh = ca.cameras[i].device.GetOutput()

		// start frame publisher
		ca.cameras[i].startFramePublisher()
	}
	return nil
}

// ShutDown shuts down all cameras
func (ca *cameraAdmin) ShutDown() error {
	ca.logger.Println("shut down cameras")
	// reset cameras
	for i, c := range ca.cameras {
		// unsubscribe all clients from the camera
		ca.cameras[i].subscriptionHandler.UnsubscribeAll()

		// close the stopPublisherCh
		// this channel is used by the frame-publisher to receive the massage that the frame-publisher should stop publishing frames
		// this is necessary before closing the cameras to avoid that the frame publisher is pulling on the frames channel of a camera before closing the camera
		if c.stopPublisherCh != nil {
			close(c.stopPublisherCh)
		}

		// based on the hardware the stopPublisherCh needs some delay time inside the frame-publisher to receive the information, that the channel is closed
		// on the raspberry pi 4 tested required minimum delay was 150ms --> 500ms should be more than enough
		time.Sleep(500 * time.Millisecond)

		// now we can close the cameras, because we can ensure, that nobody is pulling on the camera-frames anymore.
		err := ca.cameras[i].device.Close()
		if err != nil {
			errMsg := fmt.Errorf("CloseCameras() - error: closing camera (camera-id: %d, device-path: %v)", c.id, c.devicePath)
			return errMsg
		}
	}
	return nil
}

// Subscribe subscribes on a camera based on the hand-overed cameraID and returns a channel to receive the camera live view.
// The subscriberName is optional for logging purposes.
func (ca *cameraAdmin) Subscribe(cameraID int, subscriberName string) <-chan []byte {
	cameraIDX := cameraID - 1
	logCh := ca.cameras[cameraIDX].subscriptionHandler.Subscribe()

	currentSubscribers := ca.cameras[cameraIDX].subscriptionHandler.Subscriptions()
	if subscriberName != "" {
		ca.logger.Printf("client added on camera %v. client ID: %s. %d registered clients", cameraID, subscriberName, currentSubscribers)
	} else {
		ca.logger.Printf("client added on camera %v. %d registered clients", cameraID, currentSubscribers)
	}
	return logCh
}

// Unsubscribe unsubscribes from a camera based on the hand-overed cameraID and its matching log-channel.
// The subscriberName is optional for logging purposes.
func (ca *cameraAdmin) Unsubscribe(cameraID int, logChan <-chan []byte, subscriberName string) {
	cameraIDX := cameraID - 1
	ca.cameras[cameraIDX].subscriptionHandler.Unsubscribe(logChan)

	currentSubscribers := ca.cameras[cameraIDX].subscriptionHandler.Subscriptions()
	if subscriberName != "" {
		ca.logger.Printf("client removed from camera %v. client ID: %s. %d registered clients", cameraID, subscriberName, currentSubscribers)
	} else {
		ca.logger.Printf("client removed from camera %v. %d registered clients", cameraID, currentSubscribers)
	}
}

// startFramePublisher starts publishing the recorded frames with all subscribed clients.
func (c *camera) startFramePublisher() {
	go func() {
		for {
			select {
			case <-c.stopPublisherCh:
				// stop signal received, exit the publisher
				return
			case frame, ok := <-c.outputCh:
				if !ok {
					// channel was closed --> camera was shut down in the meanwhile
					return
				}
				if c.subscriptionHandler.Subscriptions() > 0 {
					c.subscriptionHandler.Publish(frame)
				} else {
					// --> no subscribed clients
				}
			}
		}
	}()
}
