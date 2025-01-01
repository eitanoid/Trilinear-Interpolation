package main

import (
	"flag"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// used as RGBA
type Vec []float64 // vector type
var Supported_Formats = map[string]bool{"rgba": true, "oklab": true}

func Print_Input(verts [2][2][2]Vec, depth int, format string) {
	_depth := strconv.Itoa(depth - 1)

	front_face := map[string]string{
		"top":    fmt.Sprintf("%s%s", verts[0][0][0].To_Ansi("000"), verts[0][0][1].To_Ansi("00"+_depth)),
		"bottom": fmt.Sprintf("%s%s", verts[0][1][0].To_Ansi("010"), verts[0][1][1].To_Ansi("0"+_depth+_depth)),
	}
	back_face := map[string]string{
		"top":    fmt.Sprintf("%s%s", verts[1][0][0].To_Ansi(_depth+"00"), verts[1][0][1].To_Ansi(_depth+"0"+_depth)),
		"bottom": fmt.Sprintf("%s%s", verts[1][1][0].To_Ansi(_depth+_depth+"0"), verts[1][1][1].To_Ansi(_depth+_depth+_depth)),
	}

	fmt.Printf("Interpolating %d times, in %s format, between: \n%s   %s\n%s   %s\n", depth, format, front_face["top"], back_face["top"], front_face["bottom"], back_face["bottom"])

}

func main() {
	//Handle user input
	_format := flag.String("format", "rgba", "In which format to interpolate.")
	_depth := flag.Int("depth", 6, "Interpolate to 'depth' points.")
	_verbose := flag.Bool("v", false, "Verbose")
	_debug := flag.Bool("d", false, "Set vertecies to debug mode")
	_generate_images := flag.Bool("i", false, "Generate 'depth' images of the interpolation")

	flag.Parse()

	format := *_format
	depth := *_depth
	verbose := *_verbose
	generate_images := *_generate_images
	debug := *_debug

	// fmt.Printf("%v,%v,%v,%v\n", format, depth, verbose, generate_images)

	corners := []RGBA{} // generate random entries for this code as RGBA
	for i := 0; i < 8; i++ {
		corners = append(corners,
			RGBA{float64(rand.Intn(((i + 1) + 1) * 30)),
				float64(rand.Intn((i + 1) * 30)),
				float64(rand.Intn((i + 1) * 30)),
				255})
	}

	if debug {
		corners = []RGBA{ // constant rather than random values for debugging
			{255, 255, 255, 255}, //#ffffff
			{0, 0, 0, 255},       //#000000
			{255, 0, 0, 255},     //#ff0000
			{0, 255, 0, 255},     //#00ff00
			{0, 0, 255, 255},     //#0000ff
			{255, 0, 255, 255},   //#ff00ff
			{255, 255, 0, 255},   //#ffff00
			{0, 255, 255, 255},   //#00ffff
		}
	}

	// verts is indexed as [forward 0 / backward 1][top 0 / bottom 1][left 0/ right 1]
	verts := [2][2][2]Vec{ // corners[n].ToLAB().ToRaw() to lerpas as OKLAB.
		{
			{corners[0].ToRaw(), corners[1].ToRaw()},
			{corners[2].ToRaw(), corners[3].ToRaw()}},

		{
			{corners[4].ToRaw(), corners[5].ToRaw()},
			{corners[6].ToRaw(), corners[7].ToRaw()}},
	}

	if verbose {
		Print_Input(verts, depth, format)
	}

	// Run Trilerp
	now := time.Now()
	cube := Trilinear_interp(verts, depth)
	fmt.Printf("Trilinear interp took: %d ms \n", time.Since(now).Milliseconds())

	// Return images
	now = time.Now()
	if generate_images {

		images := Export_Cube(cube, "rgba")
		for i, image := range images {
			Save_PNG(image, "./images/"+strconv.Itoa(i)+".png")
		}
	} else { // print colors to terminal in groups of 3 planes per row

		_cspace := 1 // how much space between each color in the planes
		_hspace := 3 // how much space between planes horizontally
		_vspace := 3 // how much space between planes vertically

		hspace := strings.Repeat(" ", _hspace)
		vspace := strings.Repeat("\n", _vspace)

		ansi_cube := Export_Cube_Ansi(cube, _cspace)
		// fmt.Println(ansi_cube)
		// fmt.Printf("%v %d \n", ansi_cube[0], len(ansi_cube[0]))
		// fmt.Printf("%v %d \n", ansi_cube[0][0], len(ansi_cube[0][0]))
		for dep, res := 0, len(cube); dep < res; dep += 3 {

			for row := 0; row < res; row++ {
				switch {
				case dep+2 < res: // print 3 points
					fmt.Printf("%s%s%s%s%s\n", ansi_cube[dep][row], hspace, ansi_cube[dep+1][row], hspace, ansi_cube[dep+2][row]) // col+space+col+space
				case dep+2 == res: // print 2 planes
					fmt.Printf("%s%s%s\n", ansi_cube[dep][row], hspace, ansi_cube[dep+1][row]) // col+space+col+space
				case dep+1 == res: // print 1 plane
					fmt.Printf("%s\n", ansi_cube[dep][row]) // col
				}
			}
			fmt.Print(vspace) // print vspace at the end of each series of planes
		}
	}

	if verbose {
		fmt.Printf("Output took: %d ms \n", time.Since(now).Milliseconds())
	}
}
