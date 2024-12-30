package main

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

func Export_Cube(cube [][][]Vec) []image.Image {
	res := len(cube)

	images := make([]image.Image, res)

	for num, slice := range cube {
		img := image.NewRGBA(image.Rect(0, 0, res, res))
		for row := 0; row < res; row++ {
			for col := 0; col < res; col++ {
				R := uint8(slice[row][col][0])
				G := uint8(slice[row][col][1])
				B := uint8(slice[row][col][2])
				A := uint8(slice[row][col][3])

				color := color.RGBA{R, G, B, A}
				img.SetRGBA(col, row, color)
			}
		}
		images[num] = img
	}
	return images
}

func SaveImg(img image.Image, name string) error {
	var err error

	f, err := os.Create(name)
	if err != nil {
		return err
	}
	defer f.Close()
	err = png.Encode(f, img)

	return err
}
