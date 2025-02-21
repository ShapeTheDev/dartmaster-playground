// TODO: DM-97 - Check imported packages for licences
// Source: https://github.com/vladimirvivien/go4vl/tree/main/v4l2
package v4l2

type TimecodeType = uint32

type TimecodeFlag = uint32

type Timecode struct {
	Type    TimecodeType
	Flags   TimecodeFlag
	Frames  uint8
	Seconds uint8
	Minutes uint8
	Hours   uint8
	_       [4]uint8 // userbits
}
