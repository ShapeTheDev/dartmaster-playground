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

// FourCCType represents the four character encoding value
type FourCCType = uint32

// Some Predefined pixel format definitions
var (
	PixelFmtRGB24 FourCCType = C.V4L2_PIX_FMT_RGB24
	PixelFmtGrey  FourCCType = C.V4L2_PIX_FMT_GREY
	PixelFmtYUYV  FourCCType = C.V4L2_PIX_FMT_YUYV
	PixelFmtYYUV  FourCCType = C.V4L2_PIX_FMT_YYUV
	PixelFmtYVYU  FourCCType = C.V4L2_PIX_FMT_YVYU
	PixelFmtUYVY  FourCCType = C.V4L2_PIX_FMT_UYVY
	PixelFmtVYUY  FourCCType = C.V4L2_PIX_FMT_VYUY
	PixelFmtMJPEG FourCCType = C.V4L2_PIX_FMT_MJPEG
	PixelFmtJPEG  FourCCType = C.V4L2_PIX_FMT_JPEG
	PixelFmtMPEG  FourCCType = C.V4L2_PIX_FMT_MPEG
	PixelFmtH264  FourCCType = C.V4L2_PIX_FMT_H264
	PixelFmtMPEG4 FourCCType = C.V4L2_PIX_FMT_MPEG4
)

type FieldType = uint32

type ColorspaceType = uint32

type YCbCrEncodingType = uint32

type HSVEncodingType = YCbCrEncodingType

type QuantizationType = uint32

type XferFunctionType = uint32

type PixFormat struct {
	Width        uint32
	Height       uint32
	PixelFormat  FourCCType
	Field        FieldType
	BytesPerLine uint32
	SizeImage    uint32
	Colorspace   ColorspaceType
	Priv         uint32
	Flags        uint32
	YcbcrEnc     YCbCrEncodingType
	HSVEnc       HSVEncodingType
	Quantization QuantizationType
	XferFunc     XferFunctionType
}

// GetPixFormat retrieves pixel information for the specified driver (via v4l2_format and v4l2_pix_format)
func GetPixFormat(fd uintptr) (PixFormat, error) {
	var v4l2Format C.struct_v4l2_format
	v4l2Format._type = C.uint(BufTypeVideoCapture)

	if err := send(fd, C.VIDIOC_G_FMT, uintptr(unsafe.Pointer(&v4l2Format))); err != nil {
		return PixFormat{}, fmt.Errorf("pix format failed: %w", err)
	}

	v4l2PixFmt := *(*C.struct_v4l2_pix_format)(unsafe.Pointer(&v4l2Format.fmt[0]))
	return PixFormat{
		Width:        uint32(v4l2PixFmt.width),
		Height:       uint32(v4l2PixFmt.height),
		PixelFormat:  uint32(v4l2PixFmt.pixelformat),
		Field:        uint32(v4l2PixFmt.field),
		BytesPerLine: uint32(v4l2PixFmt.bytesperline),
		SizeImage:    uint32(v4l2PixFmt.sizeimage),
		Colorspace:   uint32(v4l2PixFmt.colorspace),
		Priv:         uint32(v4l2PixFmt.priv),
		Flags:        uint32(v4l2PixFmt.flags),
		YcbcrEnc:     *(*uint32)(unsafe.Pointer(&v4l2PixFmt.anon0[0])),
		HSVEnc:       *(*uint32)(unsafe.Pointer(uintptr(unsafe.Pointer(&v4l2PixFmt.anon0[0])) + unsafe.Sizeof(C.uint(0)))),
		Quantization: uint32(v4l2PixFmt.quantization),
		XferFunc:     uint32(v4l2PixFmt.xfer_func),
	}, nil
}

// SetPixFormat sets the pixel format information for the specified driver
func SetPixFormat(fd uintptr, pixFmt PixFormat) error {
	var v4l2Format C.struct_v4l2_format
	v4l2Format._type = C.uint(BufTypeVideoCapture)
	*(*C.struct_v4l2_pix_format)(unsafe.Pointer(&v4l2Format.fmt[0])) = *(*C.struct_v4l2_pix_format)(unsafe.Pointer(&pixFmt))

	if err := send(fd, C.VIDIOC_S_FMT, uintptr(unsafe.Pointer(&v4l2Format))); err != nil {
		return fmt.Errorf("pix format failed: %w", err)
	}
	return nil
}
