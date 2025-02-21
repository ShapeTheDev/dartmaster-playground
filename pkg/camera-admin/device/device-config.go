// TODO: DM-97 - Check imported packages for licences
// Source: https://github.com/vladimirvivien/go4vl/tree/main/device
package device

import "github.com/One-Hundred-Eighty/Circle/pkg/camera-admin/v4l2"

type config struct {
	ioType    v4l2.IOType
	pixFormat v4l2.PixFormat
	bufSize   uint32
	fps       uint32
	bufType   uint32
}

type Option func(*config)

func WithPixFormat(pixFmt v4l2.PixFormat) Option {
	return func(o *config) {
		o.pixFormat = pixFmt
	}
}
