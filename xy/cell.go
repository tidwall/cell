package xy

// Cell is a 64-bit integer that interleaves two coordinates.
type Cell uint64

func (c Cell) String() string {
	s := make([]byte, 0, 32)
	const hex = "0123456789abcdef"
	for i := 0; i < 64; i += 4 {
		s = append(s, hex[(c>>(60-i))&15])
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

// Encode returns an encoded Cell from X/Y floating points.
// The input floating points must be within the range [0.0,1.0).
// Values outside that range are clipped.
func Encode(x, y float64) Cell {

	// Produce 32-bit integers for X/Y-> A/B
	a := uint32(clip(x) * (1 << 32))
	b := uint32(clip(y) * (1 << 32))

	// Interleave A/B into 64-bit integers AB
	ab := interleave(a)<<1 | interleave(b)

	return Cell(ab)
}

// Decode returns the decoded values from a cell.
func Decode(cell Cell) (x, y float64) {
	// Decoding is the inverse of the Encode logic.
	ab := uint64(cell)
	a := deinterleave(ab >> 1)
	b := deinterleave(ab)
	x = float64(a) / (1 << 32)
	y = float64(b) / (1 << 32)
	return x, y
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
	var ab uint64
	for i := 0; i < len(s) && i < 16; i++ {
		ab = (ab << 4) | uint64(tbl[s[i]])
	}
	return Cell(ab)
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
