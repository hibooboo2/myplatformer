//go:generage stringer -type=DIRECTION
package dir

type DIRECTION int32

const (
	NORTH      DIRECTION = 1
	WEST                 = 4
	SOUTH                = 8
	EAST                 = 16
	STOP                 = 256
	NORTH_WEST           = NORTH | WEST
	SOUTH_WEST           = SOUTH | WEST
	SOUTH_EAST           = SOUTH | EAST
	NORTH_EAST           = NORTH | EAST
)

const AXIS_TOLERANCE = 2000
