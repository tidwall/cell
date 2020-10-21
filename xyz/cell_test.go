package xyz

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/tidwall/lotsa"
)

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

			switch i % 1000 {
			case 543:
				x1 = -0.0000001
			case 264:
				y1 = -0.0000001
			case 812:
				z1 = -0.0000001
			case 643:
				x1 = 1.0000001
			case 129:
				y1 = 1.0000001
			case 362:
				z1 = 1.0000001
			}
			cell := Encode(x1, y1, z1)
			x2, y2, z2 := Decode(cell)
			if !prettyClose(x1, x2) || !prettyClose(y1, y2) ||
				!prettyClose(z1, z2) {
				t.Fatalf("\n"+
					"(%0.13f %0.13f %0.13f)\n"+
					"(%0.13f %0.13f %0.13f)\n",
					x1, y1, z1, x2, y2, z2)
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
	if os.Getenv("XYZ_PERF") != "1" {
		fmt.Printf("Use XYZ_PERF=1 to print performance details\n")
		return
	}

	rand.Seed(time.Now().UnixNano())
	N := 1_000_000
	xyz := make([][3]float64, N*3)
	cells := make([]Cell, N)
	for i := 0; i < N; i++ {
		for j := 0; j < 3; j++ {
			xyz[i][j] = rand.Float64()
		}
		cells[i] = Encode(xyz[i][0], xyz[i][1], xyz[i][2])
	}
	lotsa.Output = os.Stdout
	print("encode: ")
	lotsa.Ops(N, 1, func(i, _ int) {
		Encode(xyz[i][0], xyz[i][1], xyz[i][2])
	})
	print("decode: ")
	lotsa.Ops(N, 1, func(i, _ int) {
		Decode(cells[i])
	})

}
