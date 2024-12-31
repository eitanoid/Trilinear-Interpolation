package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

// used as RGBA
type Vec []float64 // vector type
var Supported_Formats = map[string]bool{"rgba": true, "oklab": true}

func main() {
	//Handle user input
	var depth int
	args := os.Args
	_, err := fmt.Sscanf(args[1], "%d", &depth)
	if err != nil {
		fmt.Println("Running on default values")
		depth = 10
	}

	corners := []RGBA{} // generate random entries for this code as RGBA
	for i := 0; i < 8; i++ {
		corners = append(corners,
			RGBA{float64(rand.Intn(((i + 1) + 1) * 30)),
				float64(rand.Intn((i + 1) * 30)),
				float64(rand.Intn((i + 1) * 30)),
				255})
	}

	corners = []RGBA{ // constant rather than random values for debugging
		{100, 200, 250, 255},
		{50, 20, 10, 255},
		{10, 10, 10, 255},
		{250, 100, 5, 255},
		{100, 100, 200, 255},
		{150, 150, 0, 255},
		{255, 255, 0, 255},
		{0, 0, 255, 255},
	}

	now := time.Now()
	cube := Trilinear_interp([2][2][2]Vec{ // parse as OKLAB
		{
			{corners[0].ToRaw(), corners[1].ToRaw()},
			{corners[2].ToRaw(), corners[3].ToRaw()}},

		{
			{corners[4].ToRaw(), corners[5].ToRaw()},
			{corners[6].ToRaw(), corners[7].ToRaw()}}},
		depth)
	fmt.Printf("Trilinear interp took: %d ms \n", time.Since(now).Milliseconds())

	// Return images
	now = time.Now()

	images := Export_Cube(cube, "rgba")
	for i, image := range images {
		Save_PNG(image, "./images/"+strconv.Itoa(i)+".png")
	}
	fmt.Printf("Drawing images took: %d ms \n", time.Since(now).Milliseconds())
}
