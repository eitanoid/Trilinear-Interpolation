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
	fmt.Println(corners[1])
	fmt.Println(corners[1].ToLAB().ToRGBA())
	now := time.Now()

	cube := Trilinear_interp([2][2][2]Vec{ // parse as OKLAB
		{
			{corners[0].ToLAB().ToRaw(), corners[1].ToLAB().ToRaw()},
			{corners[2].ToLAB().ToRaw(), corners[3].ToLAB().ToRaw()}},

		{
			{corners[4].ToLAB().ToRaw(), corners[5].ToLAB().ToRaw()},
			{corners[6].ToLAB().ToRaw(), corners[7].ToLAB().ToRaw()}}},
		depth)

	fmt.Printf("Trilinear interp took: %d ms \n", time.Since(now).Milliseconds())

	// Return images
	now = time.Now()

	images := Export_Cube(cube, "oklab")
	for i, image := range images {
		Save_PNG(image, "./images/"+strconv.Itoa(i)+".png")
	}
	fmt.Printf("Drawing images took: %d ms \n", time.Since(now).Milliseconds())
}
