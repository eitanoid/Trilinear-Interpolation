
| Example 1 | Example 2 |
|:- |:- |
|![example](https://github.com/eitanoid/Trilinear-Interpolation/blob/main/showcase/output.gif) | ![Example2](https://github.com/eitanoid/Trilinear-Interpolation/blob/main/showcase/ansi%20example.png)|

# Description

Create a color palette from 8 color verticies, using trilinear interpolation, supporting mixing in OkLab and RGBA color spaces. Implemented in Golang.
Passing in `-format=oklab` will enforce the color palette being perceptually uniform.

Running the `make_video.sh` will turn images in `./images` into a video or a gif that showcases the interpolation. Above is showcased a gif made from interpolating 100 steps in each direction. Video quality is higher.

# Usage

The following flags can be passed to the binary:

```bash

    -v     Verbose output. Shows the selected verticies, computation time.
    -f     Format. Select the color space to run the interpolation in. Current options are 'oklab' or 'rgba'. Default is rgba
    -depth Interpolate 'depth' points, for example '-depth=6' will return a 6x6x6 volume of points. Default is 6.
    -H     Hex. If printing to terminal, the palette will display hex color codes instead of the indicies.
    -i     Image. Generate `depth` image slices of the volume to './images'. If not set will print ansi formatted pallette to the terminal.
    -d     Debug. Sets the 8 verticies required for trilinear-interpolation to predefined values.

```


# TODO:
- User input
- Gif in readme
