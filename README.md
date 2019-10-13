# go-face-webserver
Implementation of face detection/similarity from images using FaceNet, OpenCV and Tensorflow in Golang

# requirements
To run this app you will need to install several binaries which are listed below :
* OpenCV version 3.4.2
* Tensorflow version 1.7

# models
this implementation uses model from https://github.com/davidsandberg/facenet 
which uses dataset from CASIA-WebFace and VGGFace2 as training data
and Haar filter provided by opencv for face detection created by Rainer Lienhart
