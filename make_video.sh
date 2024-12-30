#!/bin/env bash

ffmpeg -framerate 15 -i images/%d.png -vf "pad=ceil(iw/2)*2:ceil(ih/2)*2" -r 30 -pix_fmt yuv420p slideshow.mp4
