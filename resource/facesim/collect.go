package facesim

import (
	"bytes"
	"image"
	"io/ioutil"
	"log"
	"net/http"
)

func (f FacesimResource) GetImagesByteFromUrl(images []string) (res []image.Image, err error) {
	var (
		exifOrientation int
		errChan         = make(chan error)
		results         = make(chan image.Image)
		channels        int
		img             image.Image
	)

	res = []image.Image{}

	for _, item := range images {
		channels++
		go func(imgURL string) {
			defer func() {
				errChan <- err
				results <- img
			}()

			content, err := http.Get(imgURL)
			if err != nil {
				log.Println("[Facesim] Got error when getting img from url: ", err)
				return
			}

			rawbody, err := ioutil.ReadAll(content.Body)
			if err != nil {
				log.Println("[Facesim] Got error when getting img from url: ", err)
				return
			}

			// get exif
			exifs, err := f.exif.ExtractExifFromImageBytes(rawbody)
			if err != nil {
				log.Println("[Facesim]: ", err)
			}

			// decode to image
			imageReader := bytes.NewReader(rawbody)
			img, _, err = image.Decode(imageReader)
			if err != nil {
				log.Println("[Facesim]: ", err)
				return
			}

			// adjust the image orientation according to exif rotation
			if len(exifs.Orientation) > 0 {
				exifOrientation = exifs.Orientation[0]
				img = f.exif.AdjustOrientation(img, exifOrientation)
			}

		}(item)
	}

	for i := 0; i < channels; i++ {
		if procError := <-errChan; procError != nil {
			log.Println(procError)
			err = procError
		}
		if resl := <-results; resl != nil {
			res = append(res, resl)
		}
	}
	if err != nil {
		log.Println(err)
	}

	return
}
