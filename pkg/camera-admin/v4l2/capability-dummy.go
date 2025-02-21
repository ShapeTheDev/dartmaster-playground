//go:build !linux

package v4l2

type Capability struct {
}

// IsStreamingSupported returns caps & CapStreaming
func (c Capability) IsStreamingSupported() bool {
	return false
}

// GetCapability retrieves capability info for device
func GetCapability(fd uintptr) (Capability, error) {
	capability := Capability{}
	return capability, nil
}

// IsVideoCaptureSupported returns caps & CapVideoCapture
func (c Capability) IsVideoCaptureSupported() bool {
	return false
}

// IsVideoOutputSupported returns caps & CapVideoOutput
func (c Capability) IsVideoOutputSupported() bool {
	return false
}
