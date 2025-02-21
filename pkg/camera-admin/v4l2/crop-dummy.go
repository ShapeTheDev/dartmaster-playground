//go:build !linux

package v4l2

type CropCapability struct {
	DefaultRect Rect
}

// GetCropCapability  retrieves cropping info for specified device
func GetCropCapability(fd uintptr, bufType BufType) (CropCapability, error) {
	capability := CropCapability{}
	return capability, nil
}

// SetCropRect sets the cropping dimension for specified device
func SetCropRect(fd uintptr, r Rect) error {
	return nil
}
