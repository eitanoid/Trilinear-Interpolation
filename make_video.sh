#!/bin/env bash

input_folder="images"
output_file="slideshow.mp4"
framerate=30
bitrate="6000k"
crf="1" #constant rate factor
preset="veryslow"


read -p "y for mp4 video, else generates a gif (y/n)?" CONT
if [ "$CONT" = "y" ]; then
	echo "generating video";
	ffmpeg -framerate $framerate -i $input_folder/%d.png -vf "pad=ceil(iw/2)*2:ceil(ih/2)*2,boxblur=10:1,format=yuv420p" -r $framerate -b:v $bitrate  -crf $crf -preset $preset $output_file

	else
	echo "generating gif";
	
# ffmpeg -framerate $framerate -i $input_folder/%d.png -vf "fps=$framerate,scale=iw:ih:flags=lanczos" output.gif
ffmpeg -framerate $framerate -i $input_folder/%d.png -filter_complex "[0:v]fps=$framerate,scale=iw:ih:flags=lanczos,palettegen=[palette];[0:v][palette]paletteuse=dither=floyd_steinberg" output.gif

fi



