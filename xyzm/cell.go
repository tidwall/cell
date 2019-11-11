package xyzm

// Cell is a 128-bit integer that interleaves four coordinates.
type Cell struct {
	Hi uint64
	Lo uint64
}

func (c Cell) String() string {
	s := make([]byte, 0, 32)
	const hex = "0123456789abcdef"
	for i := 0; i < 64; i += 4 {
		s = append(s, hex[(c.Hi>>(60-i))&15])
	}
	for i := 0; i < 64; i += 4 {
		s = append(s, hex[(c.Lo>>(60-i))&15])
	}
	return string(s)
}

// The maximum coord that is less than 1.0. Equal to math.Nextafter(1, 0).
const maxCoord = 0.99999999999999988897769753748434595763683319091796875

func clip(x float64) float64 {
	if x < 0 {
		return 0
	}
	if x > maxCoord {
		return maxCoord
	}
	return x
}

// Encode returns an encoded Cell from X/Y/Z/M floating points.
// The input floating points must be within the range [0.0,1.0).
// Values outside that range are clipped.
func Encode(x, y, z, m float64) Cell {

	// Produce 32-bit integers for X/Y/Z/M -> A/B/C/D
	a := uint32(clip(x) * (1 << 32))
	b := uint32(clip(y) * (1 << 32))
	c := uint32(clip(z) * (1 << 32))
	d := uint32(clip(m) * (1 << 32))

	// Interleave A/C and B/D into 64-bit integers AC and BD
	ac := interleave(a)<<1 | interleave(c)
	bd := interleave(b)<<1 | interleave(d)

	// Interleave AC/BD into a single 128-bit ABCD (hi/lo) integer
	hi := interleave(uint32(ac>>32))<<1 | interleave(uint32(bd>>32))
	lo := interleave(uint32(ac))<<1 | interleave(uint32(bd))

	return Cell{Hi: hi, Lo: lo}
}

// Decode returns the decoded values from a cell.
func Decode(cell Cell) (x, y, z, m float64) {
	// Decoding is the inverse of the Encode logic.
	ac := (uint64(deinterleave(cell.Hi>>1)) << 32) |
		uint64(deinterleave(cell.Lo>>1))
	bd := (uint64(deinterleave(cell.Hi)) << 32) |
		uint64(deinterleave(cell.Lo))
	a := deinterleave(ac >> 1)
	b := deinterleave(bd >> 1)
	c := deinterleave(ac)
	d := deinterleave(bd)
	x = float64(a) / (1 << 32)
	y = float64(b) / (1 << 32)
	z = float64(c) / (1 << 32)
	m = float64(d) / (1 << 32)
	return x, y, z, m
}

// CellFromString returns the decoded values from a cell string.
func CellFromString(s string) Cell {
	const tbl = "" +
		"------------------------------------------------" +
		"\x00\x01\x02\x03\x04\x05\x06\x07\x08\x09" +
		"-------" +
		"\x0A\x0B\x0C\x0D\x0E\x0F" +
		"--------------------------" +
		"\x0A\x0B\x0C\x0D\x0E\x0F" +
		"---------------------------------------------------------" +
		"---------------------------------------------------------" +
		"---------------------------------------"
	var cell Cell
	for i := 0; i < len(s) && i < 16; i++ {
		cell.Hi = (cell.Hi << 4) | uint64(tbl[s[i]])
	}
	for i := 16; i < len(s) && i < 32; i++ {
		cell.Lo = (cell.Lo << 4) | uint64(tbl[s[i]])
	}
	return cell
}

// Bit interleaving thanks to the Daniel Lemire's blog entry:
// https://lemire.me/blog/2018/01/08/how-fast-can-you-bit-interleave-32-bit-integers/

func interleave(input uint32) uint64 {
	word := uint64(input)
	word = (word ^ (word << 16)) & 0x0000ffff0000ffff
	word = (word ^ (word << 8)) & 0x00ff00ff00ff00ff
	word = (word ^ (word << 4)) & 0x0f0f0f0f0f0f0f0f
	word = (word ^ (word << 2)) & 0x3333333333333333
	word = (word ^ (word << 1)) & 0x5555555555555555
	return word
}

func deinterleave(word uint64) uint32 {
	word &= 0x5555555555555555
	word = (word ^ (word >> 1)) & 0x3333333333333333
	word = (word ^ (word >> 2)) & 0x0f0f0f0f0f0f0f0f
	word = (word ^ (word >> 4)) & 0x00ff00ff00ff00ff
	word = (word ^ (word >> 8)) & 0x0000ffff0000ffff
	word = (word ^ (word >> 16)) & 0x00000000ffffffff
	return uint32(word)
}
