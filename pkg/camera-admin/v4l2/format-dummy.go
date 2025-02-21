//go:build !linux

package v4l2

// FourCCType represents the four character encoding value
type FourCCType = uint32

// Some Predefined pixel format definitions
var (
	PixelFmtMJPEG FourCCType = 0
)

type PixFormat struct {
	Width       uint32
	Height      uint32
	PixelFormat FourCCType
}

// GetPixFormat retrieves pixel information for the specified driver (via v4l2_format and v4l2_pix_format)
func GetPixFormat(fd uintptr) (PixFormat, error) {
	pixFormat := PixFormat{}
	return pixFormat, nil
}

// SetPixFormat sets the pixel format information for the specified driver
func SetPixFormat(fd uintptr, pixFmt PixFormat) error {
	return nil
}
