//go:build linux

// TODO: DM-97 - Check imported packages for licences
// Source: https://github.com/vladimirvivien/go4vl/tree/main/v4l2
package v4l2

// #include <linux/videodev2.h>
import "C"

import (
	"fmt"
	"unsafe"

	sys "golang.org/x/sys/unix"
)

type BufType = uint32

const (
	BufTypeVideoCapture BufType = C.V4L2_BUF_TYPE_VIDEO_CAPTURE
	BufTypeVideoOutput  BufType = C.V4L2_BUF_TYPE_VIDEO_OUTPUT
)

type IOType = uint32

const (
	IOTypeMMAP   IOType = C.V4L2_MEMORY_MMAP
	IOTypeDMABuf IOType = C.V4L2_MEMORY_DMABUF
)

type BufFlag = uint32

const (
	BufFlagMapped BufFlag = C.V4L2_BUF_FLAG_MAPPED
	BufFlagError  BufFlag = C.V4L2_BUF_FLAG_ERROR
)

type RequestBuffers struct {
	Count        uint32
	StreamType   uint32
	Memory       uint32
	Capabilities uint32
	_            [1]uint32
}

type Buffer struct {
	Index     uint32
	Type      uint32
	BytesUsed uint32
	Flags     uint32
	Field     uint32
	Timestamp sys.Timeval
	Timecode  Timecode
	Sequence  uint32
	Memory    uint32
	Info      BufferInfo
	Length    uint32
	Reserved2 uint32
	RequestFD int32
}

type BufferInfo struct {
	Offset  uint32
	UserPtr uintptr
	Planes  *Plane
	FD      int32
}

type Plane struct {
	BytesUsed  uint32
	Length     uint32
	Info       PlaneInfo
	DataOffset uint32
}

type PlaneInfo struct {
	MemOffset uint32
	UserPtr   uintptr
	FD        int32
}

// makeBuffer makes a Buffer value from C.struct_v4l2_buffer
func makeBuffer(v4l2Buf C.struct_v4l2_buffer) Buffer {
	return Buffer{
		Index:     uint32(v4l2Buf.index),
		Type:      uint32(v4l2Buf._type),
		BytesUsed: uint32(v4l2Buf.bytesused),
		Flags:     uint32(v4l2Buf.flags),
		Field:     uint32(v4l2Buf.field),
		Timestamp: *(*sys.Timeval)(unsafe.Pointer(&v4l2Buf.timestamp)),
		Timecode:  *(*Timecode)(unsafe.Pointer(&v4l2Buf.timecode)),
		Sequence:  uint32(v4l2Buf.sequence),
		Memory:    uint32(v4l2Buf.memory),
		Info:      *(*BufferInfo)(unsafe.Pointer(&v4l2Buf.m[0])),
		Length:    uint32(v4l2Buf.length),
		Reserved2: uint32(v4l2Buf.reserved2),
		RequestFD: *(*int32)(unsafe.Pointer(&v4l2Buf.anon0[0])),
	}
}

// StreamOn requests streaming to be turned on for
// capture (or output) that uses memory map, user ptr, or DMA buffers.
func StreamOn(dev StreamingDevice) error {
	bufType := dev.BufferType()
	if err := send(dev.Fd(), C.VIDIOC_STREAMON, uintptr(unsafe.Pointer(&bufType))); err != nil {
		return fmt.Errorf("stream on: %w", err)
	}
	return nil
}

// StreamOff requests streaming to be turned off for
// capture (or output) that uses memory map, user ptr, or DMA buffers.
func StreamOff(dev StreamingDevice) error {
	bufType := dev.BufferType()
	if err := send(dev.Fd(), C.VIDIOC_STREAMOFF, uintptr(unsafe.Pointer(&bufType))); err != nil {
		return fmt.Errorf("stream off: %w", err)
	}
	return nil
}

// InitBuffers sends buffer allocation request (VIDIOC_REQBUFS) to initialize buffer IO
// for video capture or video output when using either mem map, user pointer, or DMA buffers.
func InitBuffers(dev StreamingDevice) (RequestBuffers, error) {
	if dev.MemIOType() != IOTypeMMAP && dev.MemIOType() != IOTypeDMABuf {
		return RequestBuffers{}, fmt.Errorf("request buffers: %w", ErrorUnsupported)
	}
	var req C.struct_v4l2_requestbuffers
	req.count = C.uint(dev.BufferCount())
	req._type = C.uint(dev.BufferType())
	req.memory = C.uint(dev.MemIOType())

	if err := send(dev.Fd(), C.VIDIOC_REQBUFS, uintptr(unsafe.Pointer(&req))); err != nil {
		return RequestBuffers{}, fmt.Errorf("request buffers: %w: type not supported", err)
	}

	return *(*RequestBuffers)(unsafe.Pointer(&req)), nil
}

// GetBuffer retrieves buffer info for allocated buffers at provided index.
// This call should take place after buffers are allocated with RequestBuffers (for mmap for instance).
func GetBuffer(dev StreamingDevice, index uint32) (Buffer, error) {
	var v4l2Buf C.struct_v4l2_buffer
	v4l2Buf._type = C.uint(dev.BufferType())
	v4l2Buf.memory = C.uint(dev.MemIOType())
	v4l2Buf.index = C.uint(index)

	if err := send(dev.Fd(), C.VIDIOC_QUERYBUF, uintptr(unsafe.Pointer(&v4l2Buf))); err != nil {
		return Buffer{}, fmt.Errorf("query buffer: type not supported: %w", err)
	}

	return makeBuffer(v4l2Buf), nil
}

// mapMemoryBuffer creates a local buffer mapped to the address space of the device specified by fd.
func mapMemoryBuffer(fd uintptr, offset int64, len int) ([]byte, error) {
	data, err := sys.Mmap(int(fd), offset, len, sys.PROT_READ|sys.PROT_WRITE, sys.MAP_SHARED)
	if err != nil {
		return nil, fmt.Errorf("map memory buffer: %w", err)
	}
	return data, nil
}

// MapMemoryBuffers creates mapped memory buffers for specified buffer count of device.
func MapMemoryBuffers(dev StreamingDevice) ([][]byte, error) {
	bufCount := int(dev.BufferCount())
	buffers := make([][]byte, bufCount)
	for i := 0; i < bufCount; i++ {
		buffer, err := GetBuffer(dev, uint32(i))
		if err != nil {
			return nil, fmt.Errorf("mapped buffers: %w", err)
		}

		offset := buffer.Info.Offset
		length := buffer.Length
		mappedBuf, err := mapMemoryBuffer(dev.Fd(), int64(offset), int(length))
		if err != nil {
			return nil, fmt.Errorf("mapped buffers: %w", err)
		}
		buffers[i] = mappedBuf
	}
	return buffers, nil
}

// unmapMemoryBuffer removes the buffer that was previously mapped.
func unmapMemoryBuffer(buf []byte) error {
	if err := sys.Munmap(buf); err != nil {
		return fmt.Errorf("unmap memory buffer: %w", err)
	}
	return nil
}

// UnmapMemoryBuffers unmaps all mapped memory buffer for device
func UnmapMemoryBuffers(dev StreamingDevice) error {
	if dev.Buffers() == nil {
		return fmt.Errorf("unmap buffers: uninitialized buffers")
	}
	for i := 0; i < len(dev.Buffers()); i++ {
		if err := unmapMemoryBuffer(dev.Buffers()[i]); err != nil {
			return fmt.Errorf("unmap buffers: %w", err)
		}
	}
	return nil
}

// QueueBuffer enqueues a buffer in the device driver (as empty for capturing, or filled for video output)
// when using either memory map, user pointer, or DMA buffers. Buffer is returned with
// additional information about the queued buffer.
func QueueBuffer(fd uintptr, ioType IOType, bufType BufType, index uint32) (Buffer, error) {
	var v4l2Buf C.struct_v4l2_buffer
	v4l2Buf._type = C.uint(bufType)
	v4l2Buf.memory = C.uint(ioType)
	v4l2Buf.index = C.uint(index)

	if err := send(fd, C.VIDIOC_QBUF, uintptr(unsafe.Pointer(&v4l2Buf))); err != nil {
		return Buffer{}, fmt.Errorf("buffer queue: %w", err)
	}

	return makeBuffer(v4l2Buf), nil
}

// DequeueBuffer dequeues a buffer in the device driver, marking it as consumed by the application,
// when using either memory map, user pointer, or DMA buffers. Buffer is returned with
// additional information about the dequeued buffer.
func DequeueBuffer(fd uintptr, ioType IOType, bufType BufType) (Buffer, error) {
	var v4l2Buf C.struct_v4l2_buffer
	v4l2Buf._type = C.uint(bufType)
	v4l2Buf.memory = C.uint(ioType)

	err := send(fd, C.VIDIOC_DQBUF, uintptr(unsafe.Pointer(&v4l2Buf)))
	if err != nil {
		return Buffer{}, fmt.Errorf("buffer dequeue: %w", err)
	}

	return makeBuffer(v4l2Buf), nil
}
