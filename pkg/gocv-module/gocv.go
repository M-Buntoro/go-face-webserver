package gocv

import (
	"errors"
	"image"
	"log"

	"gocv.io/x/gocv"
)

type GoCVModule struct {
	facedetect gocv.CascadeClassifier
}

type Config struct {
	ClassifierPath string
}

func New(c Config) (res GoCVModule, err error) {
	// load classifier to detect faces with opencv
	classifier := gocv.NewCascadeClassifier()

	if !classifier.Load(c.ClassifierPath) {
		log.Printf("Error reading cascade file: %v\n", c.ClassifierPath)
		return res, errors.New("Failed to load cascade file")
	}

	res = GoCVModule{
		facedetect: classifier,
	}

	return
}

func (g GoCVModule) Close() (err error) {
	err = g.facedetect.Close()
	if err != nil {
		log.Println("error on closing cascades :", err)
	}
	return err
}

func (g GoCVModule) Detect(img gocv.Mat) (largest image.Image, faces []image.Image, err error) {
	bboxs := g.facedetect.DetectMultiScale(img)
	largestArea := []int{0, 0}

	if len(bboxs) <= 0 {
		log.Println("[Facesim] No face found on image")
		return
	}

	for _, bbox := range bboxs {
		crop := img.Region(bbox)
		cropped := crop.Clone()

		defer func() {
			err = crop.Close()
			if err != nil {
				log.Println("[Facesim][Google] Error in closing matrix: ", err)
			}

			err = cropped.Close()
			if err != nil {
				log.Println("[Facesim][Google] Error in closing matrix: ", err)
			}
		}()

		size := cropped.Size()
		imgcropped, err := cropped.ToImage()
		if err != nil {
			log.Println("[Facesim][OPENCV]: ", err)
			break
		}

		if len(size) < 2 {
			log.Println("[Facesim][OPENCV] Invalid size for face detection")
			continue
		}

		if largestArea[0]*largestArea[1] < size[0]*size[1] {
			largestArea[0] = size[0]
			largestArea[1] = size[1]
			largest = imgcropped
		}
		faces = append(faces, imgcropped)
	}

	return
}
