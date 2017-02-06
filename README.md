# video-capture-and-tweet
A small program that captures a frame, converts it to picture, and tweets is.

* Calls a program that uses a blackmagic card to capture a frame from an SDI or HDMI feed. 
* Uses ffmpeg to convert that frame to a picture
* Crops the picture
* Tweets the cropped picture. 

This can be done on an interval by passing in the -i <time in seconds> to the program. 
