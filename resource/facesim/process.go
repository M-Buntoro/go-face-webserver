package facesim

import (
	"errors"
	"image"
	"log"

	"gocv.io/x/gocv"
)

func (f FacesimResource) GetFacesFromImages(images []image.Image) (res []image.Image, err error) {
	for _, item := range images {
		var imres image.Image

		imgmat, err := gocv.ImageToMatRGBA(item)
		defer func() {
			err := imgmat.Close()
			if err != nil {
				log.Println("[Facesim] Error in closing matrix: ", err)
			}
		}()
		if err != nil {
			log.Println("[Facesim]: ", err)
		}

		imres, faces, err := f.gocv.Detect(imgmat)
		if len(faces) == 0 {
			err = errors.New("No face found on image")
			continue
		}

		res = append(res, imres)

	}

	if err != nil {
		log.Println(err)
	}

	return
}

func (f FacesimResource) GetFacesEmbeddings(faces []image.Image) (embeddings [][]float32, err error) {
	var (
		errChan  = make(chan error)
		results  = make(chan []float32)
		channels int
	)
	for _, item := range faces {
		channels++
		go func(img image.Image) {
			embed := []float32{}
			defer func() {
				errChan <- err
				results <- embed
			}()

			imgtensor, err := f.facenet.Preprocess(img)
			if err != nil {
				log.Println("[Facesim] Got error when get prepare image: ", err)
			}

			embed, err = f.facenet.ExtractFeatures(imgtensor)
			if err != nil {
				log.Println("[Facesim] Got error when generate embedding: ", err)
			}

		}(item)
	}

	for i := 0; i < channels; i++ {
		if procError := <-errChan; procError != nil {
			log.Println(procError)
			err = procError
		}
		if resl := <-results; resl != nil {
			log.Println(len(resl))
			embeddings = append(embeddings, resl)
		}
	}
	if err != nil {
		log.Println(err)
	}
	return
}
