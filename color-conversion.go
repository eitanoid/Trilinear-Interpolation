package main

import (
	"github.com/alltom/oklab"
	"image/color"
)

type RGBA Vec // it may be nice to visualise when something is used
type OKLAB Vec

func (v RGBA) ToRaw() Vec { // RGBA is 0 1 2 3
	return Vec(v)
}

func (v OKLAB) ToRaw() Vec {
	return Vec(v)

}

func (v RGBA) Export() color.RGBA {

	return color.RGBA{uint8(v[0]), uint8(v[1]), uint8(v[2]), uint8(v[3])}
}

func (v OKLAB) Export() color.RGBA { // convert to RGBA, export as RGBA

	return v.ToRGBA().Export()
}

// RGBA to OKLAB
func (rgba RGBA) ToLAB() OKLAB { // float vector RGBA -> package conversion -> float vector OKLAB
	color_RGBA := color.RGBA{uint8(rgba[0]), uint8(rgba[1]), uint8(rgba[2]), uint8(rgba[3])}
	_lab := oklab.OklabModel.Convert(color_RGBA).(oklab.Oklab)
	floatLab := OKLAB{float64(_lab.L), float64(_lab.A), float64(_lab.B), 255}
	return floatLab
}

// OKLAB TO RGBA
func (lab_color OKLAB) ToRGBA() RGBA { // float oklab -> package oklab -> float vector RGBA
	lab := oklab.Oklab{}
	lab.L, lab.A, lab.B = lab_color[0], lab_color[1], lab_color[2]

	r, g, b, _ := lab.RGBA()
	r, g, b = r>>8, g>>8, b>>8

	return RGBA{float64(r), float64(g), float64(b), lab_color[3]}
}

func ParseFormat(v Vec, format string) Vec { // turn all formats to raw rgba vector

	switch format {

	case "rgba":
	// do nothing
	case "oklab":
		v = OKLAB(v).ToRGBA().ToRaw()
	}
	return v
}
