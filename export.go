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

				pt := ParseFormat(slice[row][col], format)
				R := uint8(pt[0])
				G := uint8(pt[1])
				B := uint8(pt[2])
				A := uint8(pt[3])

				img.SetRGBA(col, row, color.RGBA{R, G, B, A})
			}
		}
		images[num] = img
	}
	return images
}

// Pass a 2D structure of type Vec and the color format of the Vec, return an image.Image
func Export_Plane(plane [][]Vec, format string) image.Image {
	res := len(plane)
	img := image.NewRGBA(image.Rect(0, 0, res, res))
	for row := 0; row < res; row++ {
		for col := 0; col < res; col++ {
			pt := ParseFormat(plane[row][col], format)
			R := uint8(pt[0])
			G := uint8(pt[1])
			B := uint8(pt[2])
			A := uint8(pt[3])
			color := color.RGBA{R, G, B, A}
			img.SetRGBA(col, row, color)
		}
	}
	return img
}

//export a plane into ansi escape sequences return [row]string
func Export_Plane_Ansi(plane [][]Vec, format string) []string {
	res := len(plane)
	ansi := make([]string, res)

	for row := 0; row < res; row++ {
		ansi[row] = ""
		for col := 0; col < res; col++ {
			pt := ParseFormat(plane[row][col], format)

			R := uint8(pt[0])
			G := uint8(pt[1])
			B := uint8(pt[2])
			ansi[row] += fmt.Sprintf("\033[48;2;%d;%d;%dm%s\033[0m", R, G, B, fmt.Sprintf("%d%d", row, col)) // ansi escape sequencefor background colored text
		}
	}
	return ansi
}

// if show_codes = true, show the hex value of each color, otherwise print the index
//
//export a cube into ansi escape sequences, return [depth][row]string where each string is the formatted column
func Export_Cube_Ansi(cube [][][]Vec, format string, spacing int, show_codes int) [][]string {
	res := len(cube)
	ansi_cube := make([][]string, res)
	index_channel_maxlen := len(fmt.Sprintf("%x", res-1))

	for dep, plane := range cube {

		ansi_cube[dep] = make([]string, res) // each plane is [row]string

		for row, line := range plane {
			// ansi_cube[dep][row]
			row_builder := strings.Builder{}

			for col := range line {
				pt := ParseFormat(line[col], format)

				switch show_codes {
				case 1: // 1 for hex
					row_builder.WriteString(pt.To_Ansi(pt.To_HexCode()))
					row_builder.WriteString(strings.Repeat(" ", spacing))
				case 2: // 2 for none
					row_builder.WriteString(pt.To_Ansi("  "))
				default: //indecies
					var index string
					if res > 16 { //overflows so I need to adjust spacing
						f_string := fmt.Sprintf("%%%dx%%%dx%%%dx", index_channel_maxlen, index_channel_maxlen, index_channel_maxlen)
						index = fmt.Sprintf(f_string, dep, row, col)
					} else { //normally behaved string
						index = fmt.Sprintf("%x%x%x", dep, row, col)
					}
					row_builder.WriteString(pt.To_Ansi(index))
					row_builder.WriteString(strings.Repeat(" ", spacing))
				}
				ansi_cube[dep][row] = row_builder.String()

			}

		}

	}
	return ansi_cube
}

// parse the first 3 entries in v into an ansi escape sequence with a contrasting foreground text color
func (v Vec) To_Ansi(text string) string {

	// determine the foreground color based on the RGB value's luma
	// luma := 0.299*(v[0])/255 + 0.587*(v[1])/255 + 0.144*(v[2])/255
	luma := 299*int(v[0]) + 587*int(v[1]) + 144*int(v[2])
	fg := ""

	if luma > 600*255 { // multiplying here by 255 is equivalent to dividing each component by 255
		fg = "30;"
	}

	return fmt.Sprintf("\x1b[%s48;2;%d;%d;%dm%s\x1b[m", fg, int(v[0]), int(v[1]), int(v[2]), text)
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

// parse the first 3 entries of v into an #RRGGBB hex code
func (v Vec) To_HexCode() string {
	return fmt.Sprintf("#%02X%02X%02X", int(v[0]), int(v[1]), int(v[2]))
}
