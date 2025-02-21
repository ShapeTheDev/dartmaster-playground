//go:build linux

// TODO: DM-97 - Check imported packages for licences
// Source: https://github.com/vladimirvivien/go4vl/tree/main/v4l2
package v4l2

// #include <linux/videodev2.h>
import "C"

import (
	"fmt"
	"unsafe"
)

type StreamParamFlag = uint32

type StreamParam struct {
	Type    IOType
	Capture CaptureParam
	Output  OutputParam
}

type CaptureParam struct {
	Capability   StreamParamFlag
	CaptureMode  StreamParamFlag
	TimePerFrame Fract
	ExtendedMode uint32
	ReadBuffers  uint32
	_            [4]uint32
}

type OutputParam struct {
	Capability   StreamParamFlag
	CaptureMode  StreamParamFlag
	TimePerFrame Fract
	ExtendedMode uint32
	WriteBuffers uint32
	_            [4]uint32
}

// GetStreamParam returns streaming parameters for the driver (v4l2_streamparm).
func GetStreamParam(fd uintptr, bufType BufType) (StreamParam, error) {
	var v4l2Param C.struct_v4l2_streamparm
	v4l2Param._type = C.uint(bufType)

	if err := send(fd, C.VIDIOC_G_PARM, uintptr(unsafe.Pointer(&v4l2Param))); err != nil {
		return StreamParam{}, fmt.Errorf("stream param: %w", err)
	}

	capture := *(*CaptureParam)(unsafe.Pointer(&v4l2Param.parm[0]))
	output := *(*OutputParam)(unsafe.Pointer(uintptr(unsafe.Pointer(&v4l2Param.parm[0])) + unsafe.Sizeof(C.struct_v4l2_captureparm{})))

	return StreamParam{
		Type:    BufTypeVideoCapture,
		Capture: capture,
		Output:  output,
	}, nil
}

// GetStreamParam sets streaming parameters for the driver (v4l2_streamparm).
func SetStreamParam(fd uintptr, bufType BufType, param StreamParam) error {
	var v4l2Parm C.struct_v4l2_streamparm
	v4l2Parm._type = C.uint(bufType)
	if bufType == BufTypeVideoCapture {
		*(*C.struct_v4l2_captureparm)(unsafe.Pointer(&v4l2Parm.parm[0])) = *(*C.struct_v4l2_captureparm)(unsafe.Pointer(&param.Capture))
	}
	if bufType == BufTypeVideoOutput {
		*(*C.struct_v4l2_outputparm)(unsafe.Pointer(uintptr(unsafe.Pointer(&v4l2Parm.parm[0])) + unsafe.Sizeof(v4l2Parm.parm[0]))) =
			*(*C.struct_v4l2_outputparm)(unsafe.Pointer(&param.Output))
	}

	if err := send(fd, C.VIDIOC_S_PARM, uintptr(unsafe.Pointer(&v4l2Parm))); err != nil {
		return fmt.Errorf("stream param: %w", err)
	}

	return nil
}
