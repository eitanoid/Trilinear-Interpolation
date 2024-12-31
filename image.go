package main

import (
	"image"
	"image/color"
	"image/png"
	"os"
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
