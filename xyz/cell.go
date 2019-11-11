package xyzm

import "github.com/tidwall/cell/xyzm"

// Cell is a 128-bit integer that interleaves three coordinates.
type Cell struct {
	Hi uint64
	Lo uint64
}

func (c Cell) String() string {
	return xyzm.Cell{Hi: c.Hi, Lo: c.Lo}.String()
}

// Encode returns an encoded Cell from X/Y/Z floating points.
// The input floating points must be within the range [0.0,1.0).
// Values outside that range are clipped.
func Encode(x, y, z float64) Cell {
	cell := xyzm.Encode(x, y, z, 0)
	return Cell{Hi: cell.Hi, Lo: cell.Lo}
}

// Decode returns the decoded values from a cell.
func Decode(cell Cell) (x, y, z float64) {
	x, y, z, _ = xyzm.Decode(xyzm.Cell{Hi: cell.Hi, Lo: cell.Lo})
	return x, y, z
}

// CellFromString returns the decoded values from a cell string.
func CellFromString(s string) Cell {
	cell := xyzm.CellFromString(s)
	return Cell{Hi: cell.Hi, Lo: cell.Lo}
}
