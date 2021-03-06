//go:generage stringer -type=DIRECTION
package dir

type DIRECTION int32

const (
	NORTH      DIRECTION = 1
	WEST       DIRECTION = 4
	SOUTH      DIRECTION = 8
	EAST       DIRECTION = 16
	STOP       DIRECTION = 256
	NORTH_WEST DIRECTION = NORTH | WEST
	SOUTH_WEST DIRECTION = SOUTH | WEST
	SOUTH_EAST DIRECTION = SOUTH | EAST
	NORTH_EAST DIRECTION = NORTH | EAST
	ANY_DIR    DIRECTION = NORTH | SOUTH | EAST | WEST
)

const AXIS_TOLERANCE = 10000
