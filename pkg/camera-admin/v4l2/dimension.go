// TODO: DM-97 - Check imported packages for licences
// Source: https://github.com/vladimirvivien/go4vl/tree/main/v4l2
package v4l2

type Fract struct {
	Numerator   uint32
	Denominator uint32
}

type Rect struct {
	Left   int32
	Top    int32
	Width  uint32
	Height uint32
}
