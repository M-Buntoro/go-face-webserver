package facesim

import (
	"log"

	"github.com/m-buntoro/go-face-webserver/pkg/exif"
	"github.com/m-buntoro/go-face-webserver/pkg/facenet"
	gocvmodule "github.com/m-buntoro/go-face-webserver/pkg/gocv-module"
)

func New(c Config) (res FacesimResource, err error) {

	fnet, err := facenet.NewFacenet(c.Conf.ModelPath)
	if err != nil {
		log.Println(err)
		return
	}

	err = fnet.BeginSession(&fnet)
	if err != nil {
		log.Println("fel:", err)
		return
	}

	gocv, err := gocvmodule.New(gocvmodule.Config{
		ClassifierPath: c.Conf.HaarClassifierPath,
	})
	if err != nil {
		log.Println(err)
		return
	}

	exif := exif.New()

	res = FacesimResource{
		facenet: fnet,
		exif:    exif,
		gocv:    gocv,
	}
	return res, nil
}

func (f FacesimResource) Close() (err error) {
	fnet := f.facenet.(facenet.Facenet)
	err = f.facenet.Close(&fnet)
	if err != nil {
		log.Println("error on closing facenet :", err)
		return err
	}

	f.gocv.Close()
	return nil
}
