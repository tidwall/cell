# Cell

[![GoDoc](https://img.shields.io/badge/api-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/tidwall/cell)

This project provides functions for encoding and decoding multidimensional 
z-ordered cells. A cell is an integer value that interleaves the coordinates of a point, which is useful for range and spatial based operations.

## Installing

There are three packages included in this project which handle 2, 3 and 4
dimesions. To install any or all packages:

```
go get -u github.com/tidwall/cell/xy    # 2 dimensional cells
go get -u github.com/tidwall/cell/xyz   # 3 dimensional cells
go get -u github.com/tidwall/cell/xyzm  # 4 dimensional cells
```

## Using

The input for the `Encode` function and output for the `Decode` function are floating points within the range `[0.0,1.0)`.

The 2 dimensional `xy.Encode` function results in a uint64.
The 3/4 dimensional `xyz.Encode` and `xyzm.Encode` functions result in a 128-bit integer that is represented by a struct with `Hi` and `Lo` uint64s.

## Examples

Encode/Decode a 2D point.

```go
import "github.com/tidwall/cell/xy"

cell := xy.Encode(0.331, 0.587)
fmt.Printf("%d\n", cell)

x, y := xy.Decode(cell)
fmt.Printf("%f %f\n", x, y)

// output:
// 7148508595364657900
// 0.331000 0.587000
```

Encode/Decode a Lat/Lon point.

```go
import "github.com/tidwall/cell/xy"

lat, lon := 33.1129, -112.5631
cell := xy.Encode((lon+180)/360, (lat+90)/180)
fmt.Printf("%d\n", cell)

x, y := xy.Decode(cell)
fmt.Printf("%f %f\n", x, y)

lat, lon = y*180-90, x*360-180
fmt.Printf("%f %f\n", lat, lon)

// output:
// 5548341696901379915
// 0.187325 0.683961
// 33.112900 -112.563100
```

## Performance

```
$ cd xy && XY_PERF=1 go test
encode: 1,000,000 ops in 8ms, 121,668,059/sec, 8 ns/op
decode: 1,000,000 ops in 6ms, 157,300,670/sec, 6 ns/op
```

```
$ cd xyz && XYZ_PERF=1 go test
encode: 1,000,000 ops in 19ms, 51,966,381/sec, 19 ns/op
decode: 1,000,000 ops in 19ms, 53,610,932/sec, 18 ns/op
```

```
$ cd xyzm && XYZM_PERF=1 go test
encode: 1,000,000 ops in 20ms, 49,345,055/sec, 20 ns/op
decode: 1,000,000 ops in 19ms, 52,914,482/sec, 18 ns/op
```

## Contact

Josh Baker [@tidwall](http://twitter.com/tidwall)

## License

`cell` source code is available under the MIT [License](/LICENSE).
