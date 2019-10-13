package facesim

import (
	"image"

	"github.com/m-buntoro/go-face-webserver/pkg/config"
	"github.com/m-buntoro/go-face-webserver/pkg/exif"
	"github.com/m-buntoro/go-face-webserver/pkg/facenet"
	tf "github.com/tensorflow/tensorflow/tensorflow/go"
	"gocv.io/x/gocv"
)

type GoCVInterface interface {
	Detect(img gocv.Mat) (largest image.Image, faces []image.Image, err error)
	Close() error
}

type ExifInterface interface {
	ExtractExifFromImageBytes([]byte) (exif.ImageExif, error)
	AdjustOrientation(img image.Image, exifOrientation int) image.Image
}

type FacenetInterface interface {
	Close(*facenet.Facenet) error
	BeginSession(*facenet.Facenet) (err error)
	ExtractFeatures(tensor *tf.Tensor) ([]float32, error)
	PrepareMat(img gocv.Mat) (output *tf.Tensor, err error)
	Preprocess(img image.Image) (output *tf.Tensor, err error)
	Normalize(tensor *tf.Tensor) (tensorout *tf.Tensor, err error)
}

type FacesimResource struct {
	facenet FacenetInterface
	exif    ExifInterface
	gocv    GoCVInterface
}

type Config struct {
	Conf config.Config
}
