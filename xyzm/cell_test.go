package xyzm

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/tidwall/lotsa"
)

func TestAlgorithmDetails(t *testing.T) {
	if os.Getenv("XYZM_DETAILS") != "1" {
		fmt.Printf("Use XYZM_DETAILS=1 to print algorithm details\n")
		return
	}
	// original values
	x := 0.1983718273192
	y := 0.8172398172398
	z := 0.3871239871938
	m := 0.6928301982309

	// Produce 32-bit integers for XYZM -> ABCD
	a := uint32(clip(x) * (1 << 32))
	b := uint32(clip(y) * (1 << 32))
	c := uint32(clip(z) * (1 << 32))
	d := uint32(clip(m) * (1 << 32))

	// Interleave AC and BD into 64-bit integers
	ac := interleave(a)<<1 | interleave(c)
	bd := interleave(b)<<1 | interleave(d)

	// Interleave AC and BD into a 128-bit ABCD (hi/lo) integer
	hi := interleave(uint32(ac>>32))<<1 | interleave(uint32(bd>>32))
	lo := interleave(uint32(ac))<<1 | interleave(uint32(bd))

	fmt.Printf("\nBase X/Y/Z/M values, range [0,1).\n")
	fmt.Printf("X:    %0.13f\n", x)
	fmt.Printf("Y:    %0.13f\n", y)
	fmt.Printf("Z:    %0.13f\n", z)
	fmt.Printf("M:    %0.13f\n", m)

	fmt.Printf("\nConvert X/Y/Z/M into 32-bit integers A/B/C/D.\n")
	fmt.Printf("A:    %s\n", colorize(bitstr32(a), 1))
	fmt.Printf("B:    %s\n", colorize(bitstr32(b), 1))
	fmt.Printf("C:    %s\n", colorize(bitstr32(c), 1))
	fmt.Printf("D:    %s\n", colorize(bitstr32(d), 1))

	fmt.Printf("\nInterleave A/C and B/D into 64-bit integers AC and BD.\n")
	fmt.Printf("AC:   %s\n", colorize(bitstr64(ac), 2))
	fmt.Printf("BD:   %s\n", colorize(bitstr64(bd), 2))

	fmt.Printf("\nInterleave AC/BD into a single 128-bit ABCD.\n")
	fmt.Printf("ABCD: %s\n", colorize(bitstr128(hi, lo), 4))

	fmt.Printf("\nAs string: %s\n", Cell{Hi: hi, Lo: lo}.String())
}

func prettyClose(x, y float64) bool {
	return math.Abs(x-y) < 0.00001
}

func TestEncodeDecode(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	start := time.Now()
	var count int
	for time.Since(start) < time.Millisecond*100 {
		for i := 0; i < 1000; i++ {
			x1 := rand.Float64()
			y1 := rand.Float64()
			z1 := rand.Float64()
			m1 := rand.Float64()

			switch i % 1000 {
			case 543:
				x1 = -0.0000001
			case 264:
				y1 = -0.0000001
			case 812:
				z1 = -0.0000001
			case 912:
				m1 = -0.0000001
			case 643:
				x1 = 1.0000001
			case 129:
				y1 = 1.0000001
			case 362:
				z1 = 1.0000001
			case 429:
				m1 = 1.0000001
			}
			cell := Encode(x1, y1, z1, m1)
			x2, y2, z2, m2 := Decode(cell)
			if !prettyClose(x1, x2) || !prettyClose(y1, y2) ||
				!prettyClose(z1, z2) || !prettyClose(m1, m2) {
				t.Fatalf("\n"+
					"(%0.13f %0.13f %0.13f %0.13f)\n"+
					"(%0.13f %0.13f %0.13f %0.13f)\n",
					x1, y1, z1, m1, x2, y2, z2, m2)
			}
			str := cell.String()
			cell2 := CellFromString(str)
			if cell != cell2 {
				t.Fatalf("%s != %s", cell2, cell)
			}
			count++
		}
	}
}

func TestPerf(t *testing.T) {
	if os.Getenv("XYZM_PERF") != "1" {
		fmt.Printf("Use XYZM_PERF=1 to print performance details\n")
		return
	}

	rand.Seed(time.Now().UnixNano())
	N := 1_000_000
	xyzm := make([][4]float64, N*4)
	cells := make([]Cell, N)
	for i := 0; i < N; i++ {
		for j := 0; j < 4; j++ {
			xyzm[i][j] = rand.Float64()
		}
		cells[i] = Encode(xyzm[i][0], xyzm[i][1], xyzm[i][2], xyzm[i][3])
	}
	lotsa.Output = os.Stdout
	print("encode: ")
	lotsa.Ops(N, 1, func(i, _ int) {
		Encode(xyzm[i][0], xyzm[i][1], xyzm[i][2], xyzm[i][3])
	})
	print("decode: ")
	lotsa.Ops(N, 1, func(i, _ int) {
		Decode(cells[i])
	})

}

var colors = []int{31, 32, 33, 34, 35, 36}

func bitstr32(x uint32) string {
	s := strings.Repeat("0", 32) + strconv.FormatUint(uint64(x), 2)
	return s[len(s)-32:]
}

func colorize(s string, group int) string {
	var cs string
	for i, j := 0, 0; i < len(s); i, j = i+group, j+1 {
		cs += fmt.Sprintf("\x1b[%dm%s\x1b[0m", colors[j%len(colors)], s[i:i+group])
	}
	return cs
}

func bitstr64(x uint64) string {
	s := strings.Repeat("0", 64) + strconv.FormatUint(x, 2)
	return s[len(s)-64:]
}

func bitstr128(hi, lo uint64) string {
	return bitstr64(hi) + bitstr64(lo)
}
