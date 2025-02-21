//go:build linux

// TODO: DM-97 - Check imported packages for licences
// Source: https://github.com/vladimirvivien/go4vl/tree/main/v4l2
package v4l2

/*
#cgo linux CFLAGS: -I ${SRCDIR}/../include/
#include <linux/videodev2.h>
*/
import "C"

import (
	"fmt"
	"unsafe"
)

const (
	CapVideoCapture uint32 = C.V4L2_CAP_VIDEO_CAPTURE
	CapVideoOutput  uint32 = C.V4L2_CAP_VIDEO_OUTPUT
	CapStreaming    uint32 = C.V4L2_CAP_STREAMING
)

type Capability struct {
	Driver             string
	Card               string
	BusInfo            string
	Version            uint32
	Capabilities       uint32
	DeviceCapabilities uint32
}

// IsStreamingSupported returns caps & CapStreaming
func (c Capability) IsStreamingSupported() bool {
	return c.Capabilities&CapStreaming != 0
}

// GetCapability retrieves capability info for device
func GetCapability(fd uintptr) (Capability, error) {
	var v4l2Cap C.struct_v4l2_capability
	if err := send(fd, C.VIDIOC_QUERYCAP, uintptr(unsafe.Pointer(&v4l2Cap))); err != nil {
		return Capability{}, fmt.Errorf("capability: %w", err)
	}
	return Capability{
		Driver:             C.GoString((*C.char)(unsafe.Pointer(&v4l2Cap.driver[0]))),
		Card:               C.GoString((*C.char)(unsafe.Pointer(&v4l2Cap.card[0]))),
		BusInfo:            C.GoString((*C.char)(unsafe.Pointer(&v4l2Cap.bus_info[0]))),
		Version:            uint32(v4l2Cap.version),
		Capabilities:       uint32(v4l2Cap.capabilities),
		DeviceCapabilities: uint32(v4l2Cap.device_caps),
	}, nil
}

// IsVideoCaptureSupported returns caps & CapVideoCapture
func (c Capability) IsVideoCaptureSupported() bool {
	return c.Capabilities&CapVideoCapture != 0
}

// IsVideoOutputSupported returns caps & CapVideoOutput
func (c Capability) IsVideoOutputSupported() bool {
	return c.Capabilities&CapVideoOutput != 0
}
