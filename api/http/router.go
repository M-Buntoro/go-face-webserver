package handler

import (
	"net/http"

	"github.com/go-chi/chi"

	"github.com/m-buntoro/go-face-webserver/api/http/facesim"
	"github.com/m-buntoro/go-face-webserver/pkg/config"
)

type FacesimHandlerInterface interface {
	HandleFacesimEndpoint(w http.ResponseWriter, r *http.Request)
	Close() error
	GetRoutes() func(r chi.Router)
}

type RouterModule struct {
	facesim FacesimHandlerInterface
}

type Config struct {
	Conf config.Config
}

func New(config Config) RouterModule {
	return RouterModule{
		facesim: facesim.New(facesim.Config{
			Conf: config.Conf,
		}),
	}
}

func (r RouterModule) Handler() (res *chi.Mux) {
	res = chi.NewRouter()
	res.Route("/facesim/v1", r.facesim.GetRoutes())

	return
}

func (r RouterModule) Close() (err error) {
	return r.facesim.Close()
}
