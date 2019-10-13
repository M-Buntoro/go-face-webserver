package facesim

import (
	"image"

	"github.com/m-buntoro/go-face-webserver/pkg/config"
)

type FacesimResource interface {
	GetImagesByteFromUrl([]string) ([]image.Image, error)
	GetFacesFromImages([]image.Image) ([]image.Image, error)
	GetFacesEmbeddings([]image.Image) ([][]float32, error)
	Close() error
}

type Facesim struct {
	resource FacesimResource
}

type Config struct {
	Conf config.Config
}

type FacesimRequest struct {
	Images []string `json:"images"`
}

type FacesimResult struct {
	Verdict          bool               `json:"verdict"`
	Similarity       float32            `json:"average_similarity"`
	DetailSimilarity map[string]float32 `json:"details"`
}
