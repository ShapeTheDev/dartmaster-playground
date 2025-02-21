//go:build !linux

package v4l2

type StreamParam struct {
	Capture CaptureParam
	Output  OutputParam
}

type CaptureParam struct {
	TimePerFrame Fract
}

type OutputParam struct {
	TimePerFrame Fract
}

// GetStreamParam returns streaming parameters for the driver (v4l2_streamparm).
func GetStreamParam(fd uintptr, bufType BufType) (StreamParam, error) {
	streamParam := StreamParam{}
	return streamParam, nil
}

// GetStreamParam sets streaming parameters for the driver (v4l2_streamparm).
func SetStreamParam(fd uintptr, bufType BufType, param StreamParam) error {
	return nil
}
