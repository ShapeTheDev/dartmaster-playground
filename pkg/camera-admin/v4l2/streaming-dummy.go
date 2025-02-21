//go:build !linux

package v4l2

type BufType = uint32

const (
	BufTypeVideoCapture BufType = 0
	BufTypeVideoOutput  BufType = 0
)

type IOType = uint32

const (
	IOTypeMMAP IOType = 0
)

type BufFlag = uint32

const (
	BufFlagMapped BufFlag = 0
	BufFlagError  BufFlag = 0
)

type RequestBuffers struct {
	Count uint32
}

type Buffer struct {
	Index     uint32
	BytesUsed uint32
	Flags     uint32
}

// StreamOn requests streaming to be turned on for
// capture (or output) that uses memory map, user ptr, or DMA buffers.
func StreamOn(dev StreamingDevice) error {
	return nil
}

// StreamOff requests streaming to be turned off for
// capture (or output) that uses memory map, user ptr, or DMA buffers.
func StreamOff(dev StreamingDevice) error {
	return nil
}

// InitBuffers sends buffer allocation request (VIDIOC_REQBUFS) to initialize buffer IO
// for video capture or video output when using either mem map, user pointer, or DMA buffers.
func InitBuffers(dev StreamingDevice) (RequestBuffers, error) {
	requestBuffers := RequestBuffers{}
	return requestBuffers, nil
}

// MapMemoryBuffers creates mapped memory buffers for specified buffer count of device.
func MapMemoryBuffers(dev StreamingDevice) ([][]byte, error) {
	buffers := make([][]byte, 0)
	return buffers, nil
}

// UnmapMemoryBuffers unmaps all mapped memory buffer for device
func UnmapMemoryBuffers(dev StreamingDevice) error {
	return nil
}

// QueueBuffer enqueues a buffer in the device driver (as empty for capturing, or filled for video output)
// when using either memory map, user pointer, or DMA buffers. Buffer is returned with
// additional information about the queued buffer.
func QueueBuffer(fd uintptr, ioType IOType, bufType BufType, index uint32) (Buffer, error) {
	buffer := Buffer{}
	return buffer, nil
}

// DequeueBuffer dequeues a buffer in the device driver, marking it as consumed by the application,
// when using either memory map, user pointer, or DMA buffers. Buffer is returned with
// additional information about the dequeued buffer.
func DequeueBuffer(fd uintptr, ioType IOType, bufType BufType) (Buffer, error) {
	buffer := Buffer{}
	return buffer, nil
}
