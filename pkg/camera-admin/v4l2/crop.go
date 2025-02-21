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

type CropCapability struct {
	StreamType  uint32
	Bounds      Rect
	DefaultRect Rect
	PixelAspect Fract
	_           [4]uint32
}

// GetCropCapability  retrieves cropping info for specified device
func GetCropCapability(fd uintptr, bufType BufType) (CropCapability, error) {
	var cap C.struct_v4l2_cropcap
	cap._type = C.uint(bufType)

	if err := send(fd, C.VIDIOC_CROPCAP, uintptr(unsafe.Pointer(&cap))); err != nil {
		return CropCapability{}, fmt.Errorf("crop capability: %w", err)
	}

	return *(*CropCapability)(unsafe.Pointer(&cap)), nil
}

// SetCropRect sets the cropping dimension for specified device
func SetCropRect(fd uintptr, r Rect) error {
	var crop C.struct_v4l2_crop
	crop._type = C.uint(BufTypeVideoCapture)
	crop.c = *(*C.struct_v4l2_rect)(unsafe.Pointer(&r))

	if err := send(fd, C.VIDIOC_S_CROP, uintptr(unsafe.Pointer(&crop))); err != nil {
		return fmt.Errorf("set crop: %w", err)
	}
	return nil
}
