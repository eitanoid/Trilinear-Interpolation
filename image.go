package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"strings"
)

// Passs a 3D structure of type Vec and the color format of the Vec, return a slice of images.
func Export_Cube(cube [][][]Vec, format string) []image.Image {

	if !Supported_Formats[format] {
		panic("This color format is not supported")
	}

	res := len(cube)

	images := make([]image.Image, res)

	for num, slice := range cube {
		img := image.NewRGBA(image.Rect(0, 0, res, res))
		for row := 0; row < res; row++ {
			for col := 0; col < res; col++ {

				var color_vector = slice[row][col]
				var color color.RGBA
				// var R, G, B, A uint8

				switch format {

				case "rgba":
					color = RGBA(color_vector).Export()
				case "oklab":
					color = OKLAB(color_vector).Export()
				}
				img.SetRGBA(col, row, color)
			}
		}
		images[num] = img
	}
	return images
}

func Export_Plane(plane [][]Vec) image.Image {
	res := len(plane)
	img := image.NewRGBA(image.Rect(0, 0, res, res))
	for row := 0; row < res; row++ {
		for col := 0; col < res; col++ {
			R := uint8(plane[row][col][0])
			G := uint8(plane[row][col][1])
			B := uint8(plane[row][col][2])
			A := uint8(plane[row][col][3])

			color := color.RGBA{R, G, B, A}
			img.SetRGBA(col, row, color)
		}
	}
	return img
}

//export a plane into ansi escape sequences return [row]string
func Export_Plane_Ansi(plane [][]Vec) []string {
	res := len(plane)
	ansi := make([]string, res)

	for row := 0; row < res; row++ {
		ansi[row] = ""
		for col := 0; col < res; col++ {
			R := uint8(plane[row][col][0])
			G := uint8(plane[row][col][1])
			B := uint8(plane[row][col][2])
			ansi[row] += fmt.Sprintf("\033[48;2;%d;%d;%dm%s\033[0m", R, G, B, fmt.Sprintf("%d%d", row, col)) // ansi escape sequencefor background colored text
		}
	}
	return ansi
}

//export a cube into ansi escape sequences, return [depth][row]string where each string is the formatted column
func Export_Cube_Ansi(cube [][][]Vec, spacing int) [][]string {
	res := len(cube)
	ansi_cube := make([][]string, res)

	for dep, plane := range cube {

		ansi_cube[dep] = make([]string, res) // each plane is [row]string

		for row, line := range plane {
			ansi_cube[dep][row] = ""

			for col, pt := range line {
				ansi_cube[dep][row] += pt.To_Ansi(fmt.Sprintf("%d%d%d", dep, row, col)) + strings.Repeat(" ", spacing)
			}
			// ansi_cube[dep][row] = ansi_cube[dep][row][:res-spacing]

		}

	}
	return ansi_cube
}

// turn the first 3 ONLY entries into an ansi escape sequence for background colored text
func (v Vec) To_Ansi(text string) string {
	return fmt.Sprintf("\033[48;2;%d;%d;%dm%s\033[0m", int(v[0]), int(v[1]), int(v[2]), text)
}

func Save_PNG(img image.Image, name string) error {
	var err error

	f, err := os.Create(name)
	if err != nil {
		return err
	}
	defer f.Close()
	err = png.Encode(f, img)

	return err
}
