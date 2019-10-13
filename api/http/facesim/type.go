package facesim

import (
	"github.com/m-buntoro/go-face-webserver/pkg/config"
	fservice "github.com/m-buntoro/go-face-webserver/service/facesim"
)

type FacesimRequest struct {
	fservice.FacesimRequest
}
type FacesimResult struct {
	fservice.FacesimResult
}

type FacesimService interface {
	HandleFaceSim(fservice.FacesimRequest) (fservice.FacesimResult, error)
	Close() error
}

type FacesimModule struct {
	facesimService FacesimService
}

type Config struct {
	Conf config.Config
}
