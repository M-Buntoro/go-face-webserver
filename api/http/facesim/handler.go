package facesim

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	writer "github.com/m-buntoro/go-face-webserver/pkg/response-writer"
	fservice "github.com/m-buntoro/go-face-webserver/service/facesim"
)

func New(config Config) FacesimModule {
	fcs := fservice.New(fservice.Config{
		Conf: config.Conf,
	})

	return FacesimModule{
		facesimService: fcs,
	}
}

func (f FacesimModule) Close() (err error) {
	return f.facesimService.Close()
}

func (f FacesimModule) GetRoutes() func(r chi.Router) {
	return func(r chi.Router) {
		r.Post("/request", f.HandleFacesimEndpoint)
	}
}

func (f FacesimModule) HandleFacesimEndpoint(w http.ResponseWriter, r *http.Request) {

	var (
		reqData FacesimRequest
		err     error
	)

	rqst, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		writer.ErrorWriter(w, http.StatusBadRequest, err)
		return
	}

	err = json.Unmarshal(rqst, &reqData)
	if err != nil {
		log.Println("Request Body : ", string(rqst))
		log.Println(err)
		writer.ErrorWriter(w, http.StatusBadRequest, err)
		return
	}

	if len(reqData.Images) < 2 {
		log.Println("Not enough image given")
		writer.ErrorWriter(w, http.StatusBadRequest, errors.New("Not enough image given"))
		return
	}

	res, err := f.facesimService.HandleFaceSim(reqData.FacesimRequest)

	if err != nil {
		log.Println(err)
		writer.ErrorWriter(w, http.StatusBadRequest, err)
		return
	}

	writer.SuccessWriter(w, res)

	return
}
