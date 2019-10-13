package facenet

import (
	"bytes"
	"image"
	"image/jpeg"
	_ "image/png"
	"log"
	"math"

	tf "github.com/tensorflow/tensorflow/tensorflow/go"
	"gocv.io/x/gocv"
)

func TensorFromMat(mat gocv.Mat) (output *tf.Tensor, err error) {
	GoImg, err := mat.ToImage()
	if err != nil {
		log.Println("Error on making img from mat")
		return nil, err
	}

	buffers := new(bytes.Buffer)
	err = jpeg.Encode(buffers, GoImg, nil)
	if err != nil {
		log.Println("Error on encoding image to jpeg")
		return nil, err
	}

	output, err = tf.NewTensor(string(buffers.Bytes()))
	if err != nil {
		log.Println("Error on making tensor from img")
		return nil, err
	}

	return output, nil
}

func TensorFromImg(img image.Image) (output *tf.Tensor, err error) {
	buffers := new(bytes.Buffer)
	err = jpeg.Encode(buffers, img, nil)
	if err != nil {
		log.Println("Error on encoding image to jpeg")
		return nil, err
	}

	output, err = tf.NewTensor(string(buffers.Bytes()))
	if err != nil {
		log.Println("Error on making tensor from img")
		return nil, err
	}

	return output, nil
}

func CountMean(img [][][]float32) (mean, std float32) {
	// count totals for mean
	count := len(img) * len(img[0]) * len(img[0][0])
	for _, x := range img {
		for _, y := range x {
			for _, z := range y {
				mean += z
			}
		}
	}
	mean /= float32(count)

	// count std deviation
	for _, x := range img {
		for _, y := range x {
			for _, z := range y {
				std += (z - mean) * (z - mean)
			}
		}
	}

	xstd := math.Sqrt(float64(std) / float64(count-1))
	minstd := 1.0 / math.Sqrt(float64(count))
	if xstd < minstd {
		xstd = minstd
	}

	std = float32(xstd)
	return
}
