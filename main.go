package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

// used as RGBA
type Vec []float64 // Vector

// we will indenfity a color vector by []float64{R, G, B ,A} (Vec{R,G,B,A})
// we will identify a cube by [height][width][length][color]
// verticies are idenfitied in [][]..verts indexed as 1 or 0 meaning top or bottom of the dimension for eg [1][0][0] is top height, bottom rows, bottom cols

// start and end are vectors of the same dimension, steps is the length of the resulting interpolation [][]float64
func Linear_interp(verts [2]Vec, steps int) []Vec {

	var dimension int = len(verts[0]) // the vector dimension, for RGBA colors this is 4
	if len(verts[1]) != dimension {
		panic("linear: cannot interpolate between vectors of different dimensions")
	}

	start, end := verts[0], verts[1]

	diff := make([]float64, len(start))
	for i := 0; i < dimension; i++ {
		diff[i] = (end[i] - start[i]) / float64(steps+1) //x .. x + i* diff ... x + steps*diff = y, we want the end length to be steps + 2, so we need diff to be the 1/(steps+1) fraction
	}

	//initialize the return variable
	interp := make([]Vec, steps+2) // the length after interpolating steps times
	interp[0] = start
	interp[steps+1] = end

	for i := 1; i < steps+1; i++ {
		interp[i] = make(Vec, dimension)
		for pos, val := range diff {

			interp[i][pos] = interp[0][pos] + val*float64(i)
		}

	}
	return interp
}

// varialbes are the labeled corners as b = 0 , t = 1 indexed as rowcol
// bb bt
// tb tt

// interp bb -> bt and tb -> tt then bb -> tb and bt -> tt
func Bilinear_interp(verts [2][2]Vec, steps int) [][]Vec {

	var dimension int = len(verts[0][0]) // the vector dimension, for RGBA colors this is 4

	for i := 0; i < 4; i++ {
		if len(verts[i/2][i%2]) != dimension { // 00 01 10 11
			panic("bilinear: cannot interpolate between vectors of different dimensions")
		}
	}

	// initialize
	plane := make([][]Vec, steps+2) // [row][col][color] where bb is [0][0] and tt is [steps+2][steps+2]
	for i := range plane {
		plane[i] = make([]Vec, steps+2)
	}

	first_col := Linear_interp([2]Vec{verts[0][0], verts[1][0]}, steps)
	last_col := Linear_interp([2]Vec{verts[0][1], verts[1][1]}, steps)

	for i := 0; i < steps+2; i++ { // no need to interp the 0th and steps+1th rows
		plane[i] = Linear_interp([2]Vec{first_col[i], last_col[i]}, steps) // each row is the interp of the start and end of the row
	}

	return plane
}

// looping order:
// @ is calculated # is empty
//
// @###@
// #####
// #####
// @###@
//
// this isn't written into the array yet, but the top and bottom row interp is calculated.
// @@@@@
// #####
// #####
// @@@@@
//
// use the saved slices of the top and bottom interp to do another linear interp each col, repeat this step for each col
// @@@@@
// @####
// @####
// @@@@@

// left to right the index of verts is left to right the string btb ([0][1][0])
func Trilinear_interp(verts [2][2][2]Vec, steps int) [][][]Vec { // height row col color
	// bbb, bbt, btb, btt, tbb, tbt, ttb, ttt Vec
	var dimension = len(verts[0][0][0])
	for i := 0; i < 8; i++ {
		if len(verts[i/4][(i/2)%2][i%2]) != dimension {
			panic("triliniear: cannot interpolate between vectors of different dimensions")
		}
	}

	cube := make([][][]Vec, steps+2)

	corners := make([][]Vec, steps+2) // [0,1,2,3 are bb bt tb tt], height, color
	corners[0] = Linear_interp([2]Vec{verts[0][0][0], verts[1][0][0]}, steps)
	corners[1] = Linear_interp([2]Vec{verts[0][0][1], verts[1][0][1]}, steps)
	corners[2] = Linear_interp([2]Vec{verts[0][1][0], verts[1][1][0]}, steps)
	corners[3] = Linear_interp([2]Vec{verts[0][1][1], verts[1][1][1]}, steps)

	// trilinear_interp
	// iterate over the height of the cube and interpolate each slice using the current heights bb,bt,tb,tt

	for i := 0; i < steps+2; i++ {
		bb, bt, tb, tt := corners[0][i], corners[1][i], corners[2][i], corners[3][i]
		cube[i] = Bilinear_interp([2][2]Vec{{bb, bt}, {tb, tt}}, steps)
	}

	return cube
}

// looping order:
// eg on a 4x4 cube
// input:
// 0     1     2     3
// @##@  ####  ####  @##@
// ####  ####  ####  ####
// ####  ####  ####  ####
// @##@  ####  ####  @##@
//
// calculate the corners
//
//  0     1     2     3
//  @##@  @##@  @##@  @##@
//  ####  ####  ####  ####
//  ####  ####  ####  ####
//  @##@  @##@  @##@  @##@
//
//  fill in the layers
//
//  0     1     2     3
//  @@@@  @##@  @##@  @##@
//  @@@@  ####  ####  ####
//  @@@@  ####  ####  ####
//  @@@@  @##@  @##@  @##@
//  goes this way ->

func main() {
	// test code
	// bb, bt, tb, tt := []float64{0, 0}, []float64{0, 10}, []float64{10, 0}, []float64{10, 20}
	//
	// plane := bilinear_interp(bb, bt, tb, tt, 9)
	//
	// for _, row := range plane {
	// 	fmt.Println(row)
	// }

	// bbb, bbt, btb, btt := []float64{200, 0, 0, 255}, []float64{0, 0, 10, 255}, []float64{0, 10, 0, 255}, []float64{0, 100, 10, 255}
	// tbb, tbt, ttb, ttt := []float64{30, 100, 40, 255}, []float64{0, 100, 0, 255}, []float64{5, 50, 200, 255}, []float64{255, 100, 50, 255}
	// cube := Trilinear_interp(bbb, bbt, btb, btt, tbb, tbt, ttb, ttt, 200)

	corners := []Vec{} // generate random entries for this code
	for i := 0; i < 8; i++ {
		corners = append(corners,
			Vec{float64(rand.Intn(((i + 1) + 1) * 30)),
				float64(rand.Intn((i + 1) * 30)),
				float64(rand.Intn((i + 1) * 30)),
				255})
	}

	now := time.Now()

	cube := Trilinear_interp([2][2][2]Vec{
		{{corners[0], corners[1]},
			{corners[2], corners[3]}},

		{{corners[4], corners[5]},
			{corners[6], corners[7]}}},
		200)

	fmt.Printf("Trilinear interp took: %d ms \n", time.Since(now).Milliseconds())
	now = time.Now()

	images := Export_Cube(cube)

	for i, image := range images {
		SaveImg(image, "./images/"+strconv.Itoa(i)+".png")
	}
	fmt.Printf("Drawing images took: %d ms \n", time.Since(now).Milliseconds())
}
