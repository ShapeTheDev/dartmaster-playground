// TODO: DM-97 - Check imported packages for licences
// Source: https://github.com/vladimirvivien/go4vl/tree/main/v4l2
package v4l2

import "context"

// Device is the base interface for a v4l2 device
type Device interface {
	Name() string
	Fd() uintptr
	Capability() Capability
	MemIOType() IOType
	GetOutput() <-chan []byte
	SetInput(<-chan []byte)
	Close() error
}

// StreamingDevice represents device that supports streaming IO via mapped buffer sharing.
type StreamingDevice interface {
	Device
	Buffers() [][]byte
	BufferType() BufType
	BufferCount() uint32
	Start(context.Context) error
	Stop() error
}
